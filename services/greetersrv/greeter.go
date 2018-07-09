package greetersrv

import (
	"context"
	"github.com/qeelyn/golang-starter-kit/schemas/greeter"
)

type greeterService struct {
}

func NewGreeterService() *greeterService {
	return &greeterService{}
}

func (*greeterService) Hello(ctx context.Context, req *greeter.Request) (*greeter.Response, error) {
	return &greeter.Response{Msg: "hello world"}, nil
}
