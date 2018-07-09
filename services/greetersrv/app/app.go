package app

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/qeelyn/gin-contrib/ginzap"
	"github.com/qeelyn/go-common/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	IsDebug bool
	Db      *gorm.DB
	Logger  *ginzap.Logger
)

// get User Id from grpc metadata
func GetUserContext(ctx context.Context) (*auth.Identity, error) {
	x := ctx.Value("user")
	if ua, ok := ctx.Value("user").(*auth.Identity); ok && x != nil {
		return ua, nil
	}
	return nil, status.Error(codes.Unauthenticated, "UNAUTHENTICATED")
}
