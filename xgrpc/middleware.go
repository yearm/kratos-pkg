package xgrpc

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/go-kratos/aegis/circuitbreaker"
	"github.com/go-kratos/aegis/circuitbreaker/sre"
	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-playground/validator/v10"
	"github.com/yearm/kratos-pkg/ecodes"
	"github.com/yearm/kratos-pkg/logger"
	"github.com/yearm/kratos-pkg/utils/bytesconv"
	"github.com/yearm/kratos-pkg/utils/gjson"
	"github.com/yearm/kratos-pkg/utils/group"
	"github.com/yearm/kratos-pkg/xgrpc/status"
	"google.golang.org/grpc/peer"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	// DefaultServerMiddlewares default middleware chain for servers.
	DefaultServerMiddlewares = []middleware.Middleware{
		tracing.Server(),
		metadata.Server(),
		Log(),
		Recovery(),
		RateLimit(),
		Validator(),
	}

	// DefaultClientMiddlewares default middleware chain for clients.
	DefaultClientMiddlewares = []middleware.Middleware{
		tracing.Client(),
		metadata.Client(),
		Recovery(),
		ClientBreaker(),
	}
)

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
					err = status.Error(ctx, ecodes.InternalServerError, fmt.Errorf(errString))
				}
			}()
			return handler(ctx, req)
		}
	}
}

// RateLimit is a server rate limiter middleware
func RateLimit(opts ...bbr.Option) middleware.Middleware {
	limiter := bbr.NewLimiter(opts...)
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			done, e := limiter.Allow()
			if e != nil {
				return nil, status.Error(ctx, ecodes.TooManyRequests, e)
			}
			reply, err = handler(ctx, req)
			done(ratelimit.DoneInfo{Err: err})
			return
		}
	}
}

// Log is a server logging middleware.
func Log() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			startTime := time.Now()

			var requestLog logger.RequestLog
			if tr, ok := transport.FromServerContext(ctx); ok {
				requestLog.Kind = tr.Kind().String()
				requestLog.Endpoint = tr.Endpoint()
				requestLog.Method = tr.Operation()
				requestLog.Header = gjson.MustMarshalToString(tr.RequestHeader())
				requestLog.Req = protoToString(req)
				if p, ok := peer.FromContext(ctx); ok {
					requestLog.ClientIP = p.Addr.String()
				}
			}

			reply, err = handler(ctx, req)

			level := log.LevelInfo
			responseLog := logger.ResponseLog{
				Code:  uint32(ecodes.OK),
				Reply: protoToString(reply),
			}
			st, detail, ok := status.FromError(err)
			if ok {
				responseLog.Code = uint32(st.Code())
				responseLog.Error = st.Message()
				if detail != nil {
					errorDetail := status.ToErrorDetail(detail)
					level = errorDetail.Level
					responseLog.ErrorDetail = &logger.ErrorDetailLog{
						Code:    uint32(errorDetail.Code),
						Message: errorDetail.Message,
						Level:   strings.ToLower(errorDetail.Level.String()),
						Callers: errorDetail.Callers,
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

// protoToString convert proto to string.
func protoToString(m any) string {
	if v, ok := m.(proto.Message); ok {
		b, _ := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(v)
		return bytesconv.BytesToString(b)
	}
	return fmt.Sprintf("%+v", m)
}

// Validator is a validator middleware.
func Validator() middleware.Middleware {
	validate := validator.New()
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (reply any, err error) {
			if e := validate.StructCtx(ctx, req); e != nil {
				return nil, status.Error(ctx, ecodes.InvalidArgument, e, status.WithLevel(log.LevelWarn))
			}
			return handler(ctx, req)
		}
	}
}

// ClientBreaker circuit breaker middleware will return ServiceUnavailable when the circuit
// breaker is triggered and the request is rejected directly.
func ClientBreaker() middleware.Middleware {
	gp := group.NewGroup(func() any {
		return sre.NewBreaker()
	})
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			info, _ := transport.FromClientContext(ctx)
			breaker := gp.Get(info.Operation()).(circuitbreaker.CircuitBreaker)
			if err := breaker.Allow(); err != nil {
				breaker.MarkFailed()
				return nil, status.Error(ctx, ecodes.ServiceUnavailable, err)
			}
			reply, err := handler(ctx, req)
			if err != nil && (kerrors.IsInternalServer(err) || kerrors.IsServiceUnavailable(err) || kerrors.IsGatewayTimeout(err)) {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
			return reply, err
		}
	}
}
