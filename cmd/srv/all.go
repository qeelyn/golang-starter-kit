package srv

import "github.com/qeelyn/go-common/grpcx/registry"

func RunAll(configPath *string, register registry.Registry) error {
	go RunGreeter(configPath, register)
	return RunGateway(configPath, register)
}
