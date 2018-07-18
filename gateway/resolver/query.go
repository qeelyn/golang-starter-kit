package resolver

import (
	"context"
	"github.com/qeelyn/golang-starter-kit/gateway/app"
	"github.com/qeelyn/golang-starter-kit/schemas/greeter"
)

func (t *Resolver) Hello(ctx context.Context, args struct{ Name string }) (*greeterResolver, error) {
	req := &greeter.Request{
		Name: args.Name,
	}
	res, err := app.GreeterClient.Hello(ctx, req)
	if err != nil {
		return nil, err
	}
	return &greeterResolver{
		Response: res,
	}, nil
}
