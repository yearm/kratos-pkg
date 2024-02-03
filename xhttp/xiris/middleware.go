package xiris

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	thttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/iris-contrib/middleware/cors"
	jwtMiddleware "github.com/iris-contrib/middleware/jwt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	iriscontext "github.com/kataras/iris/v12/context"
	"github.com/yearm/kratos-pkg/config/env"
	"github.com/yearm/kratos-pkg/ecode"
	"github.com/yearm/kratos-pkg/metrics"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// DefaultMiddlewares ...
var DefaultMiddlewares = []middleware.Middleware{StartAt(), tracing.Server(), Metrics(), Recovery(), Cors(), RateLimit()}

// Middlewares return middlewares wrapper
func Middlewares(m ...middleware.Middleware) iriscontext.Handler {
	chain := middleware.Chain(m...)
	return func(c iris.Context) {
		next := func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			c.ResetRequest(c.Request().WithContext(ctx))
			c.Next()
			return
		}
		next = chain(next)
		ctx := NewIrisContext(c.Request().Context(), c)
		if irisCtx, ok := FromIrisContext(ctx); ok {
			thttp.SetOperation(ctx, irisCtx.GetCurrentRoute().Path())
		}
		_, _ = next(ctx, c.Request())
	}
}

// Recovery ...
func Recovery() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			defer func() {
				if _err := recover(); _err != nil {
					buf := make([]byte, 64<<10)
					n := runtime.Stack(buf, false)
					errInfo, _ := jsoniter.MarshalToString(map[string]interface{}{
						"error": fmt.Sprintf("%v", _err),
						"req":   fmt.Sprintf("%+v", req),
						"stack": fmt.Sprintf("%s", buf[:n]),
					})
					log.Context(ctx).Error(errInfo)
					if irisCtx, ok := FromIrisContext(ctx); ok {
						irisCtx.StatusCode(http.StatusInternalServerError)
						irisCtx.StopExecution()
					}
				}
			}()
			return handler(ctx, req)
		}
	}
}

// Cors ...
func Cors() middleware.Middleware {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		MaxAge:           3600,
		AllowCredentials: false,
		Debug:            false,
	})
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if irisCtx, ok := FromIrisContext(ctx); ok {
				if irisCtx.Method() == http.MethodOptions {
					corsHandler(irisCtx)
					return
				} else {
					// NOTE: 因 iris cors Serve() 会将非 Options 请求添加 Header 后 ctx.Next()，所以此处必须单独处理非 Options 请求
					irisCtx.Header("Access-Control-Allow-Origin", "*")
					irisCtx.Header("Access-Control-Allow-Credentials", "false")
				}
			}
			return handler(ctx, req)
		}
	}
}

// StartAt ...
func StartAt() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if irisCtx, ok := FromIrisContext(ctx); ok {
				irisCtx.Values().Set("startAt", time.Now().UnixMilli())
			}
			return handler(ctx, req)
		}
	}
}

// RateLimit ...
func RateLimit() middleware.Middleware {
	limiter := bbr.NewLimiter()
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			done, err := limiter.Allow()
			if err != nil {
				if irisCtx, ok := FromIrisContext(ctx); ok {
					irisCtx.StatusCode(http.StatusTooManyRequests)
					irisCtx.StopExecution()
					return handler(ctx, req)
				}
			}
			defer done(ratelimit.DoneInfo{Err: err})
			return handler(ctx, req)
		}
	}
}

var (
	prom     *metrics.Prom
	promOnce sync.Once
)

// Metrics ...
func Metrics() middleware.Middleware {
	promOnce.Do(func() {
		labelNames := []string{"app", "method", "path", "code"}
		prom = metrics.NewProm("").
			RegisterCounter("http_request_handle_total", labelNames).
			RegisterHistogram("http_request_handle_seconds", labelNames)
	})

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			defer func() {
				if irisCtx, ok := FromIrisContext(ctx); ok {
					labels := []string{env.GetServiceName(), irisCtx.Method(), irisCtx.GetCurrentRoute().Path(), strconv.Itoa(irisCtx.GetStatusCode())}
					prom.CounterIncr(labels...)
					processTime := NewContext(irisCtx).ProcessTime()
					prom.HistogramObserve(float64(processTime)/1e3, labels...)
				}
			}()
			resp, err := handler(ctx, req)
			return resp, err
		}
	}
}

const (
	// JwtTokenKey ...
	JwtTokenKey = "jwtToken"
)

// JwtAuth iris auth middleware
func JwtAuth(secret ...string) iris.Handler {
	var _secret = "secret"
	if len(secret) > 0 {
		_secret = secret[0]
	}
	_jwt := jwtMiddleware.New(jwtMiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(_secret), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return func(ctx iriscontext.Context) {
		err := _jwt.CheckJWT(ctx)
		if err != nil {
			NewContext(ctx).Render(ecode.StatusUnauthorized, err, http.StatusUnauthorized)
			return
		}
		ctx.Values().Set(JwtTokenKey, _jwt.Get(ctx))
		ctx.Next()
	}
}
