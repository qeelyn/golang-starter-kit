package srv

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/qeelyn/gin-contrib/errorhandle"
	ginTracing "github.com/qeelyn/gin-contrib/tracing"
	"github.com/qeelyn/go-common/config"
	"github.com/qeelyn/go-common/config/options"
	"github.com/qeelyn/go-common/grpcx/dialer"
	"github.com/qeelyn/go-common/grpcx/registry"
	"github.com/qeelyn/go-common/logger"
	"github.com/qeelyn/golang-starter-kit/gateway/app"
	"github.com/qeelyn/golang-starter-kit/gateway/router"
	"github.com/qeelyn/golang-starter-kit/schemas/greeter"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func RunGateway(cnfOpts options.Options, register registry.Registry) error {
	var (
		err    error
		tracer opentracing.Tracer
	)
	cnfOpts.FileName = "gateway.yaml"
	// load application configurations
	if app.Config, err = config.LoadConfig(&cnfOpts); err != nil {
		return err
	}

	appName := app.Config.GetString("appname")
	listen := app.Config.GetString("listen")
	app.IsDebug = app.Config.GetBool("debug")
	// create the logger
	app.Logger = newLogger(app.Config)
	defer app.Logger.GetZap().Sync()

	if app.Caches, err = newBatchCache(app.Config); err != nil {
		panic(err)
	}
	//tracing
	tracer = newTracing(app.Config, appName)
	tracerCnf := map[string]interface{}{"useOpentracing": app.Config.IsSet("opentracing")}
	app.TracerFunc = ginTracing.TracingHandleFunc(tracerCnf)

	//rpc client
	cc := newDialer(app.Config.GetString("rpc.greeter"), tracer)
	app.GreeterClient = greeter.NewGreeterClient(cc)

	router := routers.DefaultRouter()
	initRouter(router)

	server := &http.Server{
		Addr:    listen,
		Handler: router,
	}

	if err = graceful(server); err != nil {
		return fmt.Errorf("Server run error:", err)
	}
	return nil
}

func newDialer(serviceName string, tracer opentracing.Tracer) *grpc.ClientConn {
	cc, err := dialer.Dial(serviceName,
		dialer.WithDialOption(
			grpc.WithInsecure(),
			grpc.WithBalancerName("round_robin"),
		),
		dialer.WithUnaryClientInterceptor(
			grpc_prometheus.UnaryClientInterceptor,
			dialer.WithAuth(),
		),
		dialer.WithTracer(tracer),
	)
	if err != nil {
		log.Panicf("dialer error: %v", err)
	}
	return cc
}

func initRouter(g *gin.Engine) {
	g.Use(app.TracerFunc)
	if app.Config.IsSet("log.access") {
		c := logger.NewFileLogger(app.Config.GetStringMap("log.access"))
		accessLogger := logger.NewLogger(c)
		g.Use(app.AccessLogHandleFunc(accessLogger.GetZap(), time.RFC3339, false))
	}
	// load error messages
	ef := app.Config.GetString("error-template")
	if ef != "" {
		g.Use(errorhandle.ErrorHandle(map[string]interface{}{
			"error-template": ef,
		}, app.Logger))
	}

	if app.Config.GetBool("jwt.enable") {
		pubKeyKey := "jwt.public-key"
		if app.Config.IsSet(pubKeyKey) {
			if err := config.ResetFromSource(app.Config, pubKeyKey); err != nil {
				panic(err)
			}
		}
		authConfig := app.Config.GetStringMap("jwt")
		//init middleware
		app.AuthHanlerFunc = app.NewAuthMiddleware(authConfig).Handle()
	}
	// check access
	if app.Config.IsSet("auth") {
		app.CheckAccessMiddleware = app.NewCheckAccessMiddleware(app.Config.GetStringMap("auth"))
	}

	routers.SetupRouterGroup(g)
	routers.SetGraphQlRouterGroup(g)
}

func graceful(srv *http.Server) error {
	go func() {
		// service connections
		log.Println("http Server is ready for listening at:", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
