package srv

import (
	"fmt"
	"github.com/qeelyn/go-common/logger"
	"github.com/spf13/viper"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap/zapcore"

	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/opentracing/opentracing-go"
	"github.com/qeelyn/go-common/config"
	"github.com/qeelyn/go-common/gormx"
	"github.com/qeelyn/go-common/grpcx"
	"github.com/qeelyn/go-common/grpcx/registry"
	"github.com/qeelyn/go-common/tracing"
	"github.com/qeelyn/golang-starter-kit/api/app"
	"github.com/qeelyn/golang-starter-kit/schemas/greeter"
	"github.com/qeelyn/golang-starter-kit/services/greetersrv"
	"path"
)

const greeterSrvName = "srv-greeter"

func RunGreeter(configPath *string, register registry.Registry) error {

	var (
		err    error
		cnf    *viper.Viper
		tracer opentracing.Tracer
	)
	if cnf, err = config.LoadConfig(path.Join(*configPath, "greeter.yaml")); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}

	appName := cnf.GetString("appname")
	listen := cnf.GetString("listen")

	app.IsDebug = cnf.GetBool("debug")
	app.Db, err = gormx.NewDb(cnf.GetStringMap("db.default"))
	if err != nil {
		panic(err)
	}
	app.Db.LogMode(app.IsDebug)
	defer app.Db.Close()

	// create the logger
	fl := logger.NewFileLogger(cnf.GetStringMap("log.file"))
	if app.IsDebug {
		app.Logger = logger.NewLogger(fl, logger.NewStdLogger())
	} else {
		app.Logger = logger.NewLogger(fl)
	}
	defer app.Logger.GetZap().Sync()

	app.Logger.ToZapField = func(values []interface{}) []zapcore.Field {
		return gormx.CreateGormLog(values).ToZapFields()
	}
	if !app.IsDebug {
		app.Db.SetLogger(app.Logger)
	}

	if cnf.IsSet("opentracing") {
		cfg := &jaegerconfig.Configuration{}
		cnf.Sub("opentracing").Unmarshal(cfg)
		tracer = tracing.NewTracer(cfg, appName)

	}

	// debug enable ?
	serverPayloadLoggingDecider := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		if fullMethodName == "healthcheck" {
			return false
		}
		return true
	}

	var opts = []grpcx.Option{
		grpcx.WithTracer(tracer),
		grpcx.WithLogger(app.Logger.GetZap()),
		grpcx.WithUnaryServerInterceptor(grpc_zap.PayloadUnaryServerInterceptor(app.Logger.GetZap(), serverPayloadLoggingDecider)),
		grpcx.WithAuthFunc(grpcx.AuthFunc(cnf.GetString("auth.public-key"))),
		grpcx.WithPrometheus(cnf.GetString("metrics.listen")),
		grpcx.WithRegistry(register, greeterSrvName, cnf.GetString("registryListen")),
	}

	server, err := grpcx.Micro(appName, opts...)

	if err != nil {
		panic(fmt.Errorf("fof server start error:%s", err))
	}

	rpc := server.BuildGrpcServer()
	greeter.RegisterGreeterServer(rpc, greetersrv.NewGreeterService())
	if err = server.Run(rpc, listen); err != nil {
		return fmt.Errorf("Server run error:", err)
	}
	return nil
}
