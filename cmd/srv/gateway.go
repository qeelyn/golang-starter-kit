package srv

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/opentracing/opentracing-go"
	"github.com/qeelyn/gin-contrib/errorhandle"
	ginTracing "github.com/qeelyn/gin-contrib/tracing"
	"github.com/qeelyn/go-common/config"
	"github.com/qeelyn/go-common/config/options"
	"github.com/qeelyn/go-common/grpcx/registry"
	"github.com/qeelyn/go-common/logger"
	"github.com/qeelyn/golang-starter-kit/gateway/app"
	"github.com/qeelyn/golang-starter-kit/gateway/router"
	"github.com/qeelyn/golang-starter-kit/schemas/greeter"
	"net/http"
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

	appName, listen := app.Config.GetString("appname"), app.Config.GetString("listen")
	app.IsDebug = app.Config.GetBool("debug")
	// create the logger
	app.Logger = newLogger(app.Config)
	defer app.Logger.Strict().Sync()
	//use grpc log for rpc client
	grpc_zap.ReplaceGrpcLogger(app.Logger.Strict())

	if app.Caches, err = newBatchCache(app.Config); err != nil {
		panic(err)
	}
	//tracing
	tracer = newTracing(app.Config, appName)

	app.TracerFunc = ginTracing.HandleFunc(map[string]interface{}{"useOpentracing": tracer != nil})

	//rpc client
	cc := newDialer(true, app.Config.Sub("rpc.greeter"), tracer)
	app.GreeterClient = greeter.NewGreeterClient(cc)
	defer cc.Close()

	router := routers.DefaultRouter()
	initRouter(router)

	server := &http.Server{
		Addr:    listen,
		Handler: router,
	}

	return server.ListenAndServe()
}

func initRouter(g *gin.Engine) {
	g.Use(app.TracerFunc)
	if glevel := app.Config.GetInt("gzip"); glevel != 0 {
		g.Use(gzip.Gzip(glevel))
	}
	if app.Config.IsSet("log.access") {
		c := logger.NewFileLogger(app.Config.GetStringMap("log.access"))
		accessLogger := logger.NewLogger(c)
		g.Use(app.AccessLogHandleFunc(accessLogger.Strict(), time.RFC3339, false))
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
