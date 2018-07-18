package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qeelyn/go-common/grpcx/registry"
	_ "github.com/qeelyn/go-common/grpcx/registry/etcdv3"
	"github.com/qeelyn/golang-starter-kit/cmd/srv"
	"google.golang.org/grpc/resolver"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	configPath    = flag.String("c", "config", "app config path")
	namingAddr    = flag.String("n", "", "service register server address")
	monitorListen = flag.String("m", "", "http listen for prometheus monitor")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Qeelyn fof services\n")
	fmt.Fprintf(os.Stderr, "USAGE\n")
	fmt.Fprintf(os.Stderr, "  serve <option> <mode> [flags]\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "Services\n")
	fmt.Fprintf(os.Stderr, "  all          Boots all services\n")
	fmt.Fprintf(os.Stderr, "  api          api gateway\n")
	fmt.Fprintf(os.Stderr, "  gteeter      greeter service\n")
	fmt.Fprintf(os.Stderr, "Options\n")
	fmt.Fprintf(os.Stderr, "  -c			  config file path\n")
	fmt.Fprintf(os.Stderr, "  -n			  service discovery address\n")
	fmt.Fprintf(os.Stderr, "  -m			  http listen for prometheus monitor\n")
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	flag.Parse()
	var run func(cp *string, register registry.Registry) error

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch cmd := strings.ToLower(os.Args[len(os.Args)-1]); cmd {
	case "all":
		run = srv.RunAll
	case "api":
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
			registry.Addrs(*namingAddr),
			registry.Timeout(30*time.Second),
		)

		if err != nil {
			panic(err)
		}

		resolver.Register(register.(resolver.Builder))
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
	if err := run(configPath, register); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
