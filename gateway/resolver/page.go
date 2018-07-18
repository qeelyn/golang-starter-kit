package resolver

import (
	"github.com/graph-gophers/graphql-go"
	"github.com/qeelyn/go-common/protobuf/paginate"
)

type pageInfoResolver struct {
	startCursor     graphql.ID
	endCursor       graphql.ID
	hasNextPage     bool
	hasPreviousPage bool
}

func newPageInfoResolver(pi *paginate.PageInfo) *pageInfoResolver {
	return &pageInfoResolver{
		startCursor:     graphql.ID(pi.StartCursor),
		endCursor:       graphql.ID(pi.EndCursor),
		hasNextPage:     pi.HasNextPage,
		hasPreviousPage: pi.HasPreviousPage,
	}
}

func (r *pageInfoResolver) StartCursor() *graphql.ID {
	return &r.startCursor
}

func (r *pageInfoResolver) EndCursor() *graphql.ID {
	return &r.endCursor
}

func (r *pageInfoResolver) HasNextPage() bool {
	return r.hasNextPage
}

func (r *pageInfoResolver) HasPreviousPage() bool {
	return r.hasPreviousPage
}

type Pagination struct {
	first  int32
	after  string
	last   int32
	before string
}

func (Pagination) ImplementsGraphQLType(name string) bool {
	return name == "Pagination"
}

func (t *Pagination) UnmarshalGraphQL(input interface{}) error {
	t, _ = input.(*Pagination)
	return nil
}
