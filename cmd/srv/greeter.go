package srv

import (
	"fmt"
	"github.com/qeelyn/go-common/config/options"
	"github.com/qeelyn/go-common/logger"
	"github.com/spf13/viper"
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
	"github.com/qeelyn/golang-starter-kit/schemas/greeter"
	"github.com/qeelyn/golang-starter-kit/services/greetersrv"
)

const greeterSrvName = "srv-greeter"

func RunGreeter(cnfOpts options.Options, register registry.Registry) error {
	var (
		err    error
		cnf    *viper.Viper
		tracer opentracing.Tracer
		db     *gorm.DB
		dlog   *logger.Logger
	)
	cnfOpts.FileName = "greeter.yaml"

	if cnf, err = config.LoadConfig(&cnfOpts); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}

	appName, listen, isDebug := cnf.GetString("appname"), cnf.GetString("listen"), cnf.GetBool("debug")
	// create the logger
	dlog = newLogger(cnf)
	defer dlog.Strict().Sync()

	dlog.ToZapField = func(values []interface{}) []zapcore.Field {
		return gormx.CreateGormLog(values).ToZapFields()
	}
	//db
	if cnf.IsSet("db") {
		db, err = gormx.NewDb(cnf.GetStringMap("db.default"))
		if err != nil {
			panic(err)
		}
		db.LogMode(isDebug)
		defer db.Close()

		if !isDebug {
			db.SetLogger(dlog)
		}
	}
	//opentracing
	tracer = newTracing(cnf, appName)

	// debug enable ?
	serverPayloadLoggingDecider := func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
		return isDebug
	}

	service := greetersrv.NewGreeterService()
	service.Db = db

	var opts = []grpcx.Option{
		grpcx.WithTracer(tracer),
		grpcx.WithLogger(dlog.Strict()),
		grpcx.WithUnaryServerInterceptor(grpc_zap.PayloadUnaryServerInterceptor(dlog.Strict(), serverPayloadLoggingDecider)),
		grpcx.WithPrometheus(cnf.GetString("metrics.listen")),
		grpcx.WithRegistry(register, greeterSrvName, cnf.GetString("registryListen")),
	}

	opts = tryAppendKeepAlive(cnf, opts)
	opts = tryAppendAuthInterceptor(cnf, opts)
	server, err := grpcx.Micro(appName, opts...)

	if err != nil {
		panic(fmt.Errorf("%s server start error:%s", greeterSrvName, err))
	}

	rpc := server.BuildGrpcServer()
	greeter.RegisterGreeterServer(rpc, service)
	if err = server.Run(rpc, listen); err != nil {
		return fmt.Errorf("%s server run error:%s", greeterSrvName, err)
	}
	return nil
}
