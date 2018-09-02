package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qeelyn/go-common/config"
	"github.com/qeelyn/go-common/config/etcdv3"
	"github.com/qeelyn/go-common/config/options"
	"github.com/qeelyn/go-common/grpcx/registry"
	_ "github.com/qeelyn/go-common/grpcx/registry/etcdv3"
	"github.com/qeelyn/golang-starter-kit/cmd/srv"
	"google.golang.org/grpc/resolver"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	configPath    = flag.String("c", "config", "app config path")
	namingAddr    = flag.String("n", "", "service register and discovery server address")
	monitorListen = flag.String("m", "", "http listen for prometheus monitor")
)

func usage() {
	fmt.Fprintf(os.Stderr, "golang start kit services\n")
	fmt.Fprintf(os.Stderr, "USAGE\n")
	fmt.Fprintf(os.Stderr, "  serve command [flags]\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "The commands are\n")
	fmt.Fprintf(os.Stderr, "  all          Boots all services\n")
	fmt.Fprintf(os.Stderr, "  gateway      Api gateway\n")
	fmt.Fprintf(os.Stderr, "  gteeter      Greeter service\n")
	fmt.Fprintf(os.Stderr, "Flags\n")
	fmt.Fprintf(os.Stderr, "  -c			  Config file path\n")
	fmt.Fprintf(os.Stderr, "  -n			  Service discovery address\n")
	fmt.Fprintf(os.Stderr, "  -m			  Http listen for prometheus monitor\n")
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {

	flag.CommandLine.Parse(os.Args[2:])
	configOptions := []options.Option{config.Path(*configPath)}

	var run func(options.Options, registry.Registry) error

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch cmd := strings.ToLower(os.Args[1]); cmd {
	case "all":
		run = srv.RunAll
	case "gateway":
		run = srv.RunGateway
	case "greeter":
		run = srv.RunGreeter
	default:
		usage()
		os.Exit(1)
	}

	//discover service
	//registry.DefaultRegistry
	//rrBalancer := balancer.Get("round_robin")
	var (
		register registry.Registry
		err      error
	)

	if *namingAddr != "" {
		register, err = registry.DefaultRegistry(
			registry.Dsn(*namingAddr),
		)

		if err != nil {
			panic(err)
		}
		resolver.Register(register.(resolver.Builder))
		configOptions = append(configOptions, config.Registry(register))
	}
	if *monitorListen != "" {
		go func() {
			log.Printf("starting prometheus http server at:%s", *monitorListen)
			http.Handle("/metrics", promhttp.Handler())
			httpServer := &http.Server{
				Addr: *monitorListen,
			}
			httpServer.ListenAndServe()
		}()
	}

	cnfOpts := config.ParseOptions(configOptions...)
	if register != nil {
		// TODO 只支持了etcd3
		etcdv3.Build(cnfOpts)
	}

	if err := run(*cnfOpts, register); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
