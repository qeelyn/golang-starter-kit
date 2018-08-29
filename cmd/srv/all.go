package srv

import (
	"github.com/qeelyn/go-common/config/options"
	"github.com/qeelyn/go-common/grpcx/registry"
)

func RunAll(cnfOpts options.Options, register registry.Registry) error {
	go RunGreeter(cnfOpts, register)
	return RunGateway(cnfOpts, register)
}
