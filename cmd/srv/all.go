package srv

import (
	"errors"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/qeelyn/go-common/cache"
	qconfig "github.com/qeelyn/go-common/config"
	"github.com/qeelyn/go-common/config/options"
	"github.com/qeelyn/go-common/grpcx"
	"github.com/qeelyn/go-common/grpcx/authfg"
	"github.com/qeelyn/go-common/grpcx/dialer"
	"github.com/qeelyn/go-common/grpcx/registry"
	"github.com/qeelyn/go-common/grpcx/tracing"
	"github.com/qeelyn/go-common/logger"
	"github.com/qeelyn/go-common/tracer"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"strings"
	"time"
)

func RunAll(cnfOpts options.Options, register registry.Registry) error {
	go RunGreeter(cnfOpts, register)
	return RunGateway(cnfOpts, register)
}

// viper is the section for rpc
func newDialer(isGateway bool, viper *viper.Viper, tracer opentracing.Tracer, gopts ...grpc.DialOption) *grpc.ClientConn {
	dialOptions := append(gopts,
		grpc.WithWaitForHandshake(),
		grpc.WithInsecure(),
		grpc.WithBalancerName("round_robin"),
	)
	if viper.GetInt("keepalive") != 0 {
		cp := keepalive.ClientParameters{Time: time.Duration(viper.GetInt("keepalive")) * time.Second}
		dialOptions = append(dialOptions, grpc.WithKeepaliveParams(cp))
	}
	cc, err := dialer.Dial(viper.GetString("name"),
		dialer.WithDialOption(dialOptions...),
		dialer.WithUnaryClientInterceptor(
			grpc_prometheus.UnaryClientInterceptor,
			authfg.WithAuthClient(isGateway),
		),
		dialer.WithTracer(tracer),
		dialer.WithTraceIdFunc(tracing.DefaultClientTraceIdFunc(isGateway)),
	)
	if err != nil {
		log.Panicf("dialer error: %v", err)
	}
	return cc
}

func newTracing(viper *viper.Viper, serviceName string) opentracing.Tracer {
	var ter opentracing.Tracer
	if viper.IsSet("opentracing") {
		cfg := &config.Configuration{}
		viper.Sub("opentracing").Unmarshal(cfg)
		ter = tracer.NewTracer(cfg, serviceName)
		if _, ok := opentracing.GlobalTracer().(opentracing.NoopTracer); ok {
			opentracing.InitGlobalTracer(ter)
		} else if strings.Contains(cfg.ServiceName, "gateway") {
			//set gateway's ter to global ter
			opentracing.InitGlobalTracer(ter)
		}
	}
	return ter
}

func newLogger(viper *viper.Viper) *logger.Logger {
	fl := logger.NewFileLogger(viper.GetStringMap("log.file"))
	if viper.GetBool("debug") {
		return logger.NewLogger(fl, logger.NewStdLogger())
	} else {
		return logger.NewLogger(fl)
	}
}

func tryAppendAuthInterceptor(viper *viper.Viper, opts []grpcx.Option) []grpcx.Option {
	if viper.GetBool("jwt.enable") {
		if viper.IsSet("jwt.public-key") {
			if err := qconfig.ResetFromSource(viper, "jwt.public-key"); err != nil {
				panic(err)
			}
		}
		opts = append(opts, grpcx.WithAuthFunc(authfg.ServerJwtAuthFunc(viper.GetStringMap("jwt"))))
	}
	return opts
}

func tryAppendKeepAlive(viper *viper.Viper, opts []grpcx.Option) []grpcx.Option {
	if viper.IsSet("keepalive") {
		ksp := keepalive.ServerParameters{
			Time: viper.GetDuration("keepalive") * time.Second,
		}
		return append(opts, grpcx.WithGrpcOption(grpc.KeepaliveParams(ksp)))
	}
	return opts
}

// mgo logger adapter
type mgoLogger struct {
	*logger.Logger
}

func NewMgoLogger(log *logger.Logger) *mgoLogger {
	return &mgoLogger{
		Logger: log,
	}
}

func (t mgoLogger) Output(calldepth int, s string) error {
	t.Logger.Sugared().Info(s)
	return nil
}

func newBatchCache(viper *viper.Viper) (caches map[string]cache.Cache, err error) {
	if viper.IsSet("cache") {
		return nil, nil
	}
	batchCnf := viper.GetStringMap("cache")
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
