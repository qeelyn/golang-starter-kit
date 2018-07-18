package errors

import (
	"errors"
	graphql "github.com/graph-gophers/graphql-go/errors"
	"github.com/qeelyn/gin-contrib/errorhandle"
)

var (
	//check data loader key,and service returns type
	ErrLoaderWrongType  = errors.New("LOADER_WRONG_TYPE")
	ErrPermissionDenied = errors.New("PERMISSION_DENIED")
	ErrUnauthorized     = errors.New("UNAUTHORIZED")
)

type slicer interface {
	Slice() []error
}

type indexedCauser interface {
	Index() int
	Cause() error
}

func Expand(errs []*graphql.QueryError) []*graphql.QueryError {
	expanded := make([]*graphql.QueryError, 0, len(errs))

	for _, err := range errs {
		switch t := err.ResolverError.(type) {
		case slicer:
			for _, e := range t.Slice() {
				qe := &graphql.QueryError{
					Message:   err.Message,
					Locations: err.Locations,
					Path:      err.Path,
				}

				if ic, ok := e.(indexedCauser); ok {
					qe.Path = append(qe.Path, ic.Index())
					if errorhandle.ErrMessage != nil {
						qe.Message = errorhandle.ErrMessage.GetErrorDescription(ic.Cause()).Message
					} else {
						qe.Message = ic.Cause().Error()
					}
				}

				expanded = append(expanded, qe)
			}
		default:
			if errorhandle.ErrMessage != nil && err.ResolverError != nil {
				err.Message = errorhandle.ErrMessage.GetErrorDescription(err.ResolverError).Message
			}
			expanded = append(expanded, err)
		}
	}

	return expanded
}
