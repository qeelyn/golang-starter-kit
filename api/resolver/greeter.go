package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/qeelyn/golang-starter-kit/helper/relay"
	"github.com/qeelyn/golang-starter-kit/schemas/greeter"
)

type greeterResolver struct {
	*greeter.Response
}

func NewGreeterResolverList(data *greeter.Response) (*greeterResolver, error) {
	ret := &greeterResolver{Response: data}
	return ret, nil
}

func (t *greeterResolver) ID() graphql.ID {
	// kind 一般统一定义在protobuf中
	kind := "greeter"
	return relay.MarshalID(kind, t.GetId())
}
