package errors

import (
	"errors"
	graphql "github.com/graph-gophers/graphql-go/errors"
	"github.com/qeelyn/gin-contrib/errorhandle"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	//check data loader key,and service returns type
	ErrLoaderWrongType  = errors.New("LOADER_WRONG_TYPE")
	ErrPermissionDenied = errors.New("PERMISSION_DENIED")
	ErrUnauthorized     = errors.New("UNAUTHORIZED")
	ErrStockNoTFound    = errors.New("StockNoTFound")
	ErrGRPCUnavailable  = errors.New("GRPCUnavailable")
	ErrDeadlineExceeded = errors.New("DeadlineExceeded")
)

func Expand(errs []*graphql.QueryError) []*graphql.QueryError {
	expanded := make([]*graphql.QueryError, 0, len(errs))

	for _, err := range errs {
		switch t := err.ResolverError.(type) {
		case interface{ GRPCStatus() *status.Status }: //for grpc
			switch t.GRPCStatus().Code() {
			case codes.DeadlineExceeded: //timeout
				err.Message = errorhandle.ErrMessage.GetErrorDescription(ErrDeadlineExceeded).Message
			case codes.Unavailable:
				err.Message = errorhandle.ErrMessage.GetErrorDescription(ErrGRPCUnavailable).Message
			default:
			}
			expanded = append(expanded, err)
		default:
			if errorhandle.ErrMessage != nil && err.ResolverError != nil {
				err.Message = errorhandle.ErrMessage.GetErrorDescription(err.ResolverError).Message
			}
			expanded = append(expanded, err)
		}
	}

	return expanded
}

//func handError(err *graphql.QueryError) *graphql.QueryError {
//	switch t := err.ResolverError.(type) {
//	case indexedCauser:
//		qe := &graphql.QueryError{
//			Message:   err.Message,
//			Locations: err.Locations,
//			Path:      err.Path,
//		}
//	default:
//		return err
//	}
//}
