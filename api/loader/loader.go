package loader

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/dataloader"
	"github.com/qeelyn/go-common/cache"
	"github.com/qeelyn/golang-starter-kit/api/app"
	"time"
)

var (
	dataloaderCache dataloader.Cache
)

const loaderInContextKey string = "loader_cls"

type Loader struct {
}

type Collection struct {
	lookup map[loaderKey]*dataloader.Loader
}

func NewLoaderCollection() Collection {
	var cacheOpt dataloader.Option
	if c, ok := app.Caches["dataloader"]; ok {
		//use cache's default duration
		dataloaderCache = NewLoaderCache(c, time.Duration(app.Config.GetInt("cache.dataloader.duration"))*time.Second)
		cacheOpt = dataloader.WithCache(dataloaderCache)
	} else {
		cacheOpt = dataloader.WithCache(&dataloader.NoCache{})
	}

	return Collection{
		lookup: map[loaderKey]*dataloader.Loader{
			UserLoaderKey: NewUserLoader(cacheOpt),
		},
	}
}

func (c Collection) Attach(ctx *gin.Context) context.Context {
	ctx.Set(loaderInContextKey, &c)
	//for k, batchFn := range c.lookup {
	//	ctx.Set(string(k), dataloader.NewBatchedLoader(batchFn, dataloader.WithCache(dataloaderCache)))
	//	//ctx = context.WithValue(ctx, k, dataloader.NewBatchedLoader(batchFn))
	//}
	return ctx
}

func (c Collection) GetLoader(k loaderKey) *dataloader.Loader {
	ldr, ok := c.lookup[k]
	if ok {
		return ldr
	}
	switch k {
	}
	return nil
}

func extract(ctx context.Context, k loaderKey) (*dataloader.Loader, error) {
	// k need same type as attach ctx.Set type
	coll, ok := ctx.Value(loaderInContextKey).(*Collection)
	if !ok {
		return nil, fmt.Errorf("unable to find %s loader on the request context", k)
	}
	// find the loader
	ldr := coll.GetLoader(k)
	if ldr == nil {
		return nil, fmt.Errorf("unable to find %s loader on loader collection", k)
	}

	return ldr, nil
}

func Load(tk loaderKey, ctx context.Context, key dataloader.Key) (interface{}, error) {
	ldr, err := extract(ctx, tk)
	if err != nil {
		return nil, err
	}
	thunk := ldr.Load(ctx, key)
	data, err := thunk()
	if err != nil {
		ldr.Clear(ctx, key)
		return nil, err
	}
	return data, nil
}

type loaderCache struct {
	cache    cache.Cache
	duration time.Duration
}

func NewLoaderCache(c cache.Cache, duration time.Duration) *loaderCache {
	return &loaderCache{cache: c, duration: duration}
}

func (l *loaderCache) Get(ctx context.Context, key dataloader.Key) (dataloader.Thunk, bool) {
	var thunk dataloader.Thunk
	err := l.cache.Get(key.String(), &thunk)
	if err != nil {
		return nil, false
	}

	return thunk, true
}

func (l *loaderCache) Set(ctx context.Context, key dataloader.Key, thunk dataloader.Thunk) {
	l.cache.Set(key.String(), thunk, l.duration)
}

func (l *loaderCache) Delete(ctx context.Context, key dataloader.Key) bool {
	return l.cache.Delete(key.String()) == nil
}

func (l *loaderCache) Clear() {
	l.cache.FlushAll()
}
