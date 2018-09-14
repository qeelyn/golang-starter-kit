package app

import (
	_ "github.com/patrickmn/go-cache"
	"github.com/qeelyn/go-common/cache"
	_ "github.com/qeelyn/go-common/cache/local"
	_ "github.com/qeelyn/go-common/cache/memcache"
	"github.com/spf13/viper"

	"context"
	"github.com/jinzhu/gorm"
	"github.com/qeelyn/go-common/logger"
	"github.com/qeelyn/golang-starter-kit/schemas/greeter"
)

var (
	Config        *viper.Viper
	IsDebug       bool
	Cache         cache.Cache
	Caches        map[string]cache.Cache
	Logger        *logger.Logger
	Db            *gorm.DB
	GreeterClient greeter.GreeterClient
)

func init() {
	Caches = make(map[string]cache.Cache)
}

func GetUserId(ctx context.Context) string {
	v := ctx.Value("userid")
	if v == nil {
		return ""
	}
	return v.(string)
}

func GetOrgId(ctx context.Context) string {
	v := ctx.Value("orgid")
	if v == nil {
		return "0"
	}
	return v.(string)
}
