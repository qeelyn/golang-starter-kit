package srv

import (
	"fmt"
	"github.com/qeelyn/go-common/logger"
	"github.com/spf13/viper"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap/zapcore"

	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/opentracing/opentracing-go"
	"github.com/qeelyn/go-common/config"
	"github.com/qeelyn/go-common/gormx"
	"github.com/qeelyn/go-common/grpcx"
	"github.com/qeelyn/go-common/grpcx/registry"
	"github.com/qeelyn/go-common/tracing"
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
		db     *gorm.DB
		dlog   *logger.Logger
	)
	if cnf, err = config.LoadConfig(path.Join(*configPath, "greeter.yaml")); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}

	appName := cnf.GetString("appname")
	listen := cnf.GetString("listen")

	isDebug := cnf.GetBool("debug")
	db, err = gormx.NewDb(cnf.GetStringMap("db.default"))
	if err != nil {
		panic(err)
	}
	db.LogMode(isDebug)
	defer db.Close()

	// create the logger
	fl := logger.NewFileLogger(cnf.GetStringMap("log.file"))
	if isDebug {
		dlog = logger.NewLogger(fl, logger.NewStdLogger())
	} else {
		dlog = logger.NewLogger(fl)
	}
	defer dlog.GetZap().Sync()

	dlog.ToZapField = func(values []interface{}) []zapcore.Field {
		return gormx.CreateGormLog(values).ToZapFields()
	}
	if !isDebug {
		db.SetLogger(dlog)
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
		return isDebug
	}

	service := greetersrv.NewGreeterService()
	service.Db = db
	service.Logger = dlog

	var opts = []grpcx.Option{
		grpcx.WithTracer(tracer),
		grpcx.WithLogger(service.Logger.GetZap()),
		grpcx.WithUnaryServerInterceptor(grpc_zap.PayloadUnaryServerInterceptor(service.Logger.GetZap(), serverPayloadLoggingDecider)),
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
