package handle

import (
	"context"
	navErr "errors"
	"github.com/gin-gonic/gin"
	"github.com/graph-gophers/graphql-go"
	"github.com/qeelyn/gin-contrib/auth"
	"github.com/qeelyn/go-common/logger"
	"github.com/qeelyn/golang-starter-kit/gateway/app"
	"github.com/qeelyn/golang-starter-kit/gateway/errors"
	"github.com/qeelyn/golang-starter-kit/gateway/loader"
	"github.com/qeelyn/golang-starter-kit/gateway/resolver"
	"github.com/qeelyn/golang-starter-kit/gateway/schema"
	"net/http"
	"sync"
)

type query struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type request struct {
	queries []query
	isBatch bool
}

type GraphQL struct {
	Schema      *graphql.Schema
	Loaders     loader.Collection
	CheckAccess *auth.CheckAccess
}

func ServeGraphqlResource(r *gin.RouterGroup) {
	//tracer := graphql.Tracer(opentracing.GlobalTracer())
	graphqlSchema := graphql.MustParseSchema(schema.GetRootSchema(), &resolver.Resolver{})
	graphql.Logger(NewGraphqlLogger(app.Logger))
	checkAccess := app.NewCheckAccess(app.Config.GetStringMap("auth"))
	h := &GraphQL{
		Schema:      graphqlSchema,
		Loaders:     loader.NewLoaderCollection(),
		CheckAccess: checkAccess,
	}
	r.POST("query", h.Query)
}

type graphqlLoggerAdapter struct {
	logger *logger.Logger
}

func NewGraphqlLogger(l *logger.Logger) *graphqlLoggerAdapter {
	return &graphqlLoggerAdapter{
		logger: l,
	}
}

func (t *graphqlLoggerAdapter) LogPanic(ctx context.Context, value interface{}) {
	app.Logger.Error(value)
}

func (t *GraphQL) Query(c *gin.Context) {
	req, err := parse(c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	n := len(req.queries)
	if n == 0 {
		c.AbortWithError(http.StatusBadRequest, navErr.New("err-request"))
		return
	}
	var (
		ctx       = t.Loaders.Attach(c)
		responses = make([]*graphql.Response, n)
		wg        sync.WaitGroup
	)
	wg.Add(n)
	for i, q := range req.queries {
		go func(i int, q query) {
			res := t.Schema.Exec(ctx, q.Query, q.OperationName, q.Variables)
			res.Errors = errors.Expand(res.Errors)

			responses[i] = res
			wg.Done()
		}(i, q)
	}

	wg.Wait()
	if req.isBatch {
		c.JSON(200, responses)
	} else if len(responses) > 0 {
		c.JSON(200, responses[0])
	}
}
