package app

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	"github.com/qeelyn/gin-contrib/auth"
	"github.com/qeelyn/gin-contrib/tracing"
	errors2 "github.com/qeelyn/golang-starter-kit/gateway/errors"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"
)

var (
	AuthMiddleware        *auth.GinJWTMiddleware
	CheckAccessMiddleware *auth.CheckAccess
	checkAccessUrl        string
	checkAccessTimeout    int
	JeagerTracer          *tracing.JeagerTracer
)

// Ginzap returns a gin.HandlerFunc (middleware) that logs requests using uber-go/zap.
//
// Requests with errors are logged using zap.Error().
// Requests without errors are logged using zap.Info().
//
func AccessLogHandleFunc(logger *zap.Logger, timeFormat string, utc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// some evil middlewares modify this values
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		bodyCopy := &bytes.Buffer{}
		if c.Request.Method == "POST" {
			switch c.ContentType() {
			case binding.MIMEJSON, binding.MIMEPOSTForm, binding.MIMEXML:
				io.Copy(bodyCopy, c.Request.Body)
				c.Request.Body = ioutil.NopCloser(bytes.NewReader(bodyCopy.Bytes()))
			}
		}
		if orgId := c.GetHeader("Qeelyn-Org-Id"); orgId != "" {
			c.Set("orgid", orgId)
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		if utc {
			end = end.UTC()
		}

		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.ByteString("body", bodyCopy.Bytes()),
			zap.String("ip", c.ClientIP()),
			zap.String("auth", c.GetString("userId")),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("time", end.Format(timeFormat)),
			zap.Duration("latency", latency),
			tracing.TracingField(c, logger),
		)
	}
}

// auth will check the jwt token basically
func NewAuthMiddleware(config map[string]interface{}) *auth.GinJWTMiddleware {
	// the jwt middleware
	pKey, _ := config["public-key"].([]byte)
	eKey, _ := config["encryption-key"].(string)
	algo, _ := config["algorithm"].(string)
	if strings.HasPrefix(algo, "RS") && (pKey == nil) {
		panic("miss pubKeyFile or priKeyFile setting when in RS signing algorithm")
	}
	if strings.HasPrefix(algo, "HS") && eKey == "" {
		panic("miss encryption-key setting when in HS signing algorithm")
	}
	middle := &auth.GinJWTMiddleware{
		Realm:            "auth server",
		PubKeyFile:       pKey,
		Key:              []byte(eKey),
		SigningAlgorithm: config["algorithm"].(string), //RS256
		UnauthorizedHandle: func(c *gin.Context, code int, message string) bool {
			if IsDebug && c.GetHeader("Authorization") == "" {
				if tid, ok := config["testuserid"]; ok {
					c.Set("userId", tid.(string))
				}
				c.Next()
				return false
			}
			c.JSON(code, gin.H{
				"errors": []map[string]interface{}{
					{
						"code":    code,
						"message": message,
					},
				},
			})
			return true
		},
		TokenValidator: func(token *jwt.Token, c *gin.Context) bool {
			c.Set("authorization", c.GetHeader("Authorization"))
			return true
		},
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
	}
	return middle

}

// userId will be exist after bearer auth middleware execute
func NewCheckAccessMiddleware(config map[string]interface{}) *auth.CheckAccess {
	checkAccessUrl = path.Join(config["auth-server"].(string), config["check-access"].(string))
	routerPrefix := config["router-prefix"].(string)
	checkAccessTimeout = config["check-access-timeout"].(int)
	instance := &auth.CheckAccess{
		GetPermissionFunc: func(context *gin.Context) string {
			reqPath := context.Request.URL.Path
			if strings.HasPrefix(reqPath, routerPrefix) {
				return reqPath[len(routerPrefix):]
			} else {
				return reqPath
			}
		},
		CheckFunc: func(context *http.Request, userId string, permission string, params map[string]interface{}) int {
			if IsDebug && context.Header.Get("Authorization") == "" {
				return http.StatusOK
			}
			body, err := json.Marshal(map[string]interface{}{
				"permission": permission,
				"params":     params,
			})
			if err != nil {
				Logger.Errorf("error on CheckFunc : %s", err)
				return http.StatusBadRequest
			}
			client := http.Client{
				Timeout: time.Duration(checkAccessTimeout) * time.Millisecond,
			}
			req, _ := http.NewRequest("POST", checkAccessUrl, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", context.Header.Get("Authorization"))

			if authRes, err := client.Do(req); err == nil {
				return authRes.StatusCode
			} else {
				Logger.Errorf("error on auth client request : %s", err)
				return http.StatusInternalServerError
			}
		},
	}
	return instance
}

func CheckAccess(ctx context.Context, permission string, params map[string]interface{}) (bool, error) {
	var userId, orgId string
	var ok bool
	var err error
	if userId, ok = ctx.Value("userid").(string); !ok {
		err = errors2.ErrUnauthorized
	}
	if orgId, ok = ctx.Value("orgid").(string); ok {
		if params == nil {
			params = map[string]interface{}{}
		}
		params["org_id"] = orgId
	}
	req := ctx.Value(0).(*http.Request)

	if code := CheckAccessMiddleware.CheckFunc(req, userId, permission, params); code != http.StatusOK {
		if code == http.StatusForbidden {
			err = errors2.ErrPermissionDenied
			Logger.Warnf("userId %s has no permission at %s", userId, permission)
		}
		err = errors.New(http.StatusText(code))
	}
	return err == nil, err
}

func NewJeagerTracer() gin.HandlerFunc {
	JeagerTracer = &tracing.JeagerTracer{}
	return JeagerTracer.Handle()
}
