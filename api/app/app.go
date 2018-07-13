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

func NewCache(cnf map[string]interface{}) error {
	for key, value := range cnf {
		config := value.(map[string]interface{})
		if ins, err := cache.NewCache(config["type"].(string), config); err != nil {
			return err
		} else {
			Caches[key] = ins
			if Cache == nil && key == "default" {
				Cache = ins
			}
		}
	}
	if Cache == nil {
		panic("initial cache failure,miss default cache")
	}
	if len(Caches) == 0 {
		panic("initial cache failure,please check the config")
	}
	return nil
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
