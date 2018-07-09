package resolver

import (
	"github.com/graph-gophers/graphql-go"
)

type nodeInterface interface {
	ID() graphql.ID
}

type nodeResolver struct {
	nodeInterface
}

func (t *nodeResolver) ToGreeter() (*greeterResolver, bool) {
	c, ok := t.nodeInterface.(*greeterResolver)
	return c, ok
}
