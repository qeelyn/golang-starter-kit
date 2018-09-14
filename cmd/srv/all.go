package srv

import (
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/qeelyn/go-common/cache"
	qconfig "github.com/qeelyn/go-common/config"
	"github.com/qeelyn/go-common/config/options"
	"github.com/qeelyn/go-common/grpcx"
	"github.com/qeelyn/go-common/grpcx/registry"
	"github.com/qeelyn/go-common/logger"
	"github.com/qeelyn/go-common/tracing"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"
)

func RunAll(cnfOpts options.Options, register registry.Registry) error {
	go RunGreeter(cnfOpts, register)
	return RunGateway(cnfOpts, register)
}

func newTracing(viper *viper.Viper, serviceName string) opentracing.Tracer {
	var tracer opentracing.Tracer
	if viper.IsSet("opentracing") {
		cfg := &config.Configuration{}
		viper.Sub("opentracing").Unmarshal(cfg)
		tracer = tracing.NewTracer(cfg, serviceName)
		if tracer != nil && opentracing.GlobalTracer() == nil {
			opentracing.InitGlobalTracer(tracer)
		}
	}
	return tracer
}

func newLogger(viper *viper.Viper) *logger.Logger {
	fl := logger.NewFileLogger(viper.GetStringMap("log.file"))
	if viper.GetBool("debug") {
		return logger.NewLogger(fl, logger.NewStdLogger())
	} else {
		return logger.NewLogger(fl)
	}
}

func appendAuthInterceptor(viper *viper.Viper, opts []grpcx.Option) []grpcx.Option {
	if viper.GetBool("jwt.enable") {
		if viper.IsSet("jwt.public-key") {
			if err := qconfig.ResetFromSource(viper, "jwt.public-key"); err != nil {
				panic(err)
			}
		}
		opts = append(opts, grpcx.WithAuthFunc(grpcx.JwtAuthFunc(viper.GetStringMap("jwt"))))
	}
	return opts
}

func newBatchCache(batchCnf map[string]interface{}) (caches map[string]cache.Cache, err error) {
	caches = make(map[string]cache.Cache)
	for key, value := range batchCnf {
		cnf := value.(map[string]interface{})
		if ins, err := cache.NewCache(cnf["type"].(string), cnf); err != nil {
			return nil, err
		} else {
			caches[key] = ins
		}
	}
	if len(caches) == 0 {
		return nil, errors.New("initial cache failure,please check the config")
	}
	return
}
