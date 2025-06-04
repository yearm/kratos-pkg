package xgin

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	thttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/yearm/kratos-pkg/ecodes"
	"github.com/yearm/kratos-pkg/logger"
	"github.com/yearm/kratos-pkg/utils/bytesconv"
	"github.com/yearm/kratos-pkg/utils/gjson"
	"github.com/yearm/kratos-pkg/xhttp/api"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// MiddlewareOption middleware option.
type MiddlewareOption func(*middlewareOptions)

// middlewareOptions middleware options.
type middlewareOptions struct {
	tracingOptions []tracing.Option
	logValuerMap   map[string]LogValuer
	bbrOptions     []bbr.Option
	// ...
}

// WithMiddlewareTracingOptions used to set the tracing.Server options.
func WithMiddlewareTracingOptions(tracingOpts []tracing.Option) MiddlewareOption {
	return func(options *middlewareOptions) {
		options.tracingOptions = tracingOpts
	}
}

// WithMiddlewareLogValuerMap used to set the Log options.
func WithMiddlewareLogValuerMap(m map[string]LogValuer) MiddlewareOption {
	return func(options *middlewareOptions) {
		options.logValuerMap = m
	}
}

// WithMiddlewareBBROptions used to set the RateLimit options.
func WithMiddlewareBBROptions(bbrOpts []bbr.Option) MiddlewareOption {
	return func(options *middlewareOptions) {
		options.bbrOptions = bbrOpts
	}
}

// DefaultServerMiddlewares default middleware chain for servers.
func DefaultServerMiddlewares(opts ...MiddlewareOption) []middleware.Middleware {
	opt := middlewareOptions{}
	for _, o := range opts {
		o(&opt)
	}
	ms := []middleware.Middleware{
		tracing.Server(opt.tracingOptions...),
		Log(opt.logValuerMap),
		Recovery(),
		RateLimit(opt.bbrOptions...),
	}
	return ms
}

// Middlewares return middlewares wrapper
func Middlewares(m []middleware.Middleware) gin.HandlerFunc {
	chain := middleware.Chain(m...)
	return func(c *gin.Context) {
		next := func(ctx context.Context, req any) (any, error) {
			c.Next()
			return c.Writer, nil
		}
		next = chain(next)
		ctx := NewGinContext(c.Request.Context(), c)
		c.Request = c.Request.WithContext(ctx)
		thttp.SetOperation(ctx, c.FullPath())
		if ginCtx, ok := FromGinContext(ctx); ok {
			thttp.SetOperation(ctx, ginCtx.FullPath())
		}
		_, _ = next(ctx, c.Request)
	}
}

// Recovery is a recovery middleware.
func Recovery() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			defer func() {
				if e := recover(); e != nil {
					buf := make([]byte, 64<<10)
					n := runtime.Stack(buf, false)
					errString := fmt.Sprintf("%v", e)
					errInfo := gjson.MustMarshalToString(map[string]any{
						"error": errString,
						"req":   fmt.Sprintf("%+v", req),
						"stack": fmt.Sprintf("%s", buf[:n]),
					})
					log.Context(ctx).Error(errInfo)
					if ginCtx, ok := FromGinContext(ctx); ok {
						r := api.NewError(ctx, ecodes.InternalServerError, fmt.Errorf(errString)).WithHTTPCode(http.StatusInternalServerError)
						RenderWithAbort(ginCtx, r)
						return
					}
				}
			}()
			return handler(ctx, req)
		}
	}
}

// LogValuer is returns a log value.
type LogValuer func(ctx *gin.Context) any

// Log is a server logging middleware.
func Log(logValuerMap ...map[string]LogValuer) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			startTime := time.Now()

			var requestLog logger.RequestLog
			if tr, ok := transport.FromServerContext(ctx); ok {
				requestLog.Kind = tr.Kind().String()
				requestLog.Endpoint = tr.Endpoint()
				requestLog.Header = gjson.MustMarshalToString(tr.RequestHeader())
			}
			if ginCtx, ok := FromGinContext(ctx); ok {
				requestLog.Method = ginCtx.HandlerName()
				requestLog.Req = gjson.MustMarshalToString(map[string]any{
					"method": ginCtx.Request.Method,
					"uri":    ginCtx.Request.RequestURI,
					"proto":  ginCtx.Request.Proto,
					"body": func() string {
						form, _ := ginCtx.MultipartForm()
						if form != nil && len(form.File) > 0 {
							files := make([]string, 0, len(form.File))
							for _, headers := range form.File {
								for _, header := range headers {
									files = append(files, fmt.Sprintf("%v[%v]", header.Filename, bytesconv.FormatFileSize(header.Size)))
								}
							}
							return strings.Join(files, ", ")
						}
						body, err := ginCtx.GetRawData()
						if err != nil {
							return fmt.Sprintf("ginCtx.GetRawData failed: %v", err)
						}
						ginCtx.Request.Body = io.NopCloser(bytes.NewReader(body))
						return bytesconv.BytesToString(body)
					}(),
				})
				requestLog.ClientIP = ginCtx.ClientIP()
				if len(logValuerMap) > 0 {
					extra := make(map[string]any)
					for key, valuer := range logValuerMap[0] {
						extra[key] = valuer(ginCtx)
					}
					requestLog.Extra = gjson.MustMarshalToString(extra)
				}
			}

			reply, err = handler(ctx, req)

			level := log.LevelInfo
			var responseLog logger.ResponseLog
			if ginCtx, ok := FromGinContext(ctx); ok {
				apiResponse, _ := ginCtx.Get(responseKey)
				if v, ok := apiResponse.(*api.Response); ok {
					level = v.GetLevel()
					responseLog.Code = uint32(v.GetHTTPCode())
					responseLog.Reply = v.RenderString()
					if v.Code != ecodes.OK {
						responseLog.Error = func() string {
							var errString string
							if v.GetError() != nil {
								errString = v.GetError().Error()
							}
							return errString
						}()
						responseLog.ErrorDetail = &logger.ErrorDetailLog{
							Code:    uint32(v.Code),
							Message: v.Message,
							Level:   strings.ToLower(v.GetLevel().String()),
							Callers: v.GetCallers(),
						}
					}
				}
			}
			responseLog.Latency = time.Since(startTime).Milliseconds()

			log.Context(ctx).Log(level,
				log.DefaultMessageKey, "",
				logger.RequestKey, requestLog,
				logger.ResponseKey, responseLog,
			)
			return
		}
	}
}

// RateLimit is a server rate limiter middleware
func RateLimit(opts ...bbr.Option) middleware.Middleware {
	limiter := bbr.NewLimiter(opts...)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			done, err := limiter.Allow()
			if err != nil {
				if ginCtx, ok := FromGinContext(ctx); ok {
					r := api.NewError(ctx, ecodes.TooManyRequests, err).WithHTTPCode(http.StatusTooManyRequests)
					RenderWithAbort(ginCtx, r)
					return
				}
			}
			reply, err = handler(ctx, req)
			done(ratelimit.DoneInfo{Err: err})
			return
		}
	}
}
