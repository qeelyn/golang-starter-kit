package srv

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/qeelyn/gin-contrib/errorhandle"
	"github.com/qeelyn/gin-contrib/ginzap"
	"github.com/qeelyn/go-common/grpcx/dialer"
	"github.com/qeelyn/go-common/grpcx/registry"
	"github.com/qeelyn/go-common/tracing"
	"github.com/qeelyn/golang-starter-kit/api/app"
	"github.com/qeelyn/golang-starter-kit/api/router"
	"github.com/qeelyn/golang-starter-kit/helper/apph"
	"github.com/qeelyn/golang-starter-kit/schemas/greeter"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"
)

func RunApi(configPath *string, register registry.Registry) error {
	var (
		err error
	)
	// load application configurations
	if app.Config, err = apph.LoadConfig(path.Join(*configPath, "api.yaml")); err != nil {
		return err
	}

	appName := app.Config.GetString("appname")
	listen := app.Config.GetString("listen")

	app.IsDebug = app.Config.GetBool("debug")
	// create the logger
	fl := ginzap.NewFileLogger(app.Config.GetStringMap("log.file"))
	if app.IsDebug {
		app.Logger = ginzap.NewLogger(fl, ginzap.NewStdLogger())
	} else {
		app.Logger = ginzap.NewLogger(fl)
	}
	defer app.Logger.GetZap().Sync()
	// cache
	cacheConfig := app.Config.GetStringMap("cache")
	if cacheConfig != nil {
		if err := app.NewCache(cacheConfig); err != nil {
			return err
		}
	}

	cfg := &jaegerconfig.Configuration{}
	app.Config.Sub("opentracing").Unmarshal(cfg)
	tracer := tracing.NewTracer(cfg, appName)
	if tracer != nil {
		opentracing.InitGlobalTracer(tracer)
	}

	//rpc client
	cc := newDialer(app.Config.GetString("rpc.fof"), tracer)

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
	if app.Config.IsSet("log.access") {
		c := ginzap.NewFileLogger(app.Config.GetStringMap("log.access"))
		accessLogger := ginzap.NewLogger(c)
		g.Use(app.AccessLogHandleFunc(accessLogger.GetZap(), time.RFC3339, false))
	}
	// load error messages
	ef := app.Config.GetString("error-template")
	if ef != "" {
		g.Use(errorhandle.ErrorHandle(map[string]interface{}{
			"error-template": ef,
		}, app.Logger))
	}

	authConfig := app.Config.GetStringMap("auth")
	//init middleware
	app.AuthHandleFunc(authConfig)
	app.CheckAccessHandleFunc(authConfig)

	routers.SetupRouterGroup(g)
	routers.SetGraphQlRouterGroup(g)
}

func graceful(srv *http.Server) error {
	go func() {
		// service connections
		log.Println("Server is ready for listening at:", srv.Addr)
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
