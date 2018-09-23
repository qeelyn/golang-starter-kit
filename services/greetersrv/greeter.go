package greetersrv

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/qeelyn/golang-starter-kit/schemas/greeter"
)

type greeterService struct {
	Db *gorm.DB
}

func NewGreeterService() *greeterService {
	return &greeterService{}
}

func (*greeterService) Hello(ctx context.Context, req *greeter.Request) (*greeter.Response, error) {
	return &greeter.Response{Msg: "hello world"}, nil
}
