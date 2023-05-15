package middleware

import (
	"context"
	"errors"
	"fmt"
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
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/yearm/kratos-pkg/ecode"
	"github.com/yearm/kratos-pkg/status"
	"github.com/yearm/kratos-pkg/util/group"
	"google.golang.org/grpc/codes"
	"reflect"
	"runtime"
	"strings"
	"time"
)

var (
	// DefaultServerMiddlewares ...
	DefaultServerMiddlewares = []middleware.Middleware{StartAt(), tracing.Server(), Log(), Recovery(), metadata.Server(), RateLimit(), Validator()}
	// DefaultClientMiddlewares ...
	DefaultClientMiddlewares = []middleware.Middleware{Recovery(), tracing.Client(), metadata.Client(), ClientBreaker()}
)

// Validator ...
func Validator() middleware.Middleware {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		return field.Tag.Get("errMsg")
	})
	zhT := zh.New()
	uni := ut.New(zhT, zhT)
	trans, _ := uni.GetTranslator("zh")
	if err := zhtrans.RegisterDefaultTranslations(validate, trans); err != nil {
		logrus.Panicln("validator middleware register translations error:", err)
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if _err := validate.Struct(req); _err != nil {
				if fieldErrors, ok := (_err).(validator.ValidationErrors); ok {
					for _, fieldError := range fieldErrors {
						var errMsg string
						translateValue := fieldError.Translate(trans)
						// NOTE: Field() 和 StructField() 不相等说明取到了 errMsg tag 值
						if fieldError.Field() != fieldError.StructField() {
							// NOTE: 翻译时取的值是 Field()，由于前面 RegisterTagNameFunc 取的是 errMsg tag 对应的值，所以这里翻译后要替换成 StructField()
							translateValue = strings.Replace(translateValue, fieldError.Field(), fieldError.StructField(), 1)
							errMsg = fieldError.Field()
						}
						_fieldError := errors.New(fmt.Sprintf("%s:%s", fieldError.StructNamespace(), translateValue))
						if errMsg != "" {
							return nil, status.ErrorWithMsg(_fieldError, ecode.StatusInvalidRequest, errMsg)
						}
						return nil, status.Error(_fieldError, ecode.StatusInvalidRequest)
					}
				}
				return nil, status.Error(_err, ecode.StatusInvalidRequest)
			}
			return handler(ctx, req)
		}
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
					errStr := fmt.Sprintf("%v", _err)
					errInfo, _ := jsoniter.MarshalToString(map[string]interface{}{
						"error": errStr,
						"req":   fmt.Sprintf("%+v", req),
						"stack": fmt.Sprintf("%s", buf[:n]),
					})
					log.Context(ctx).Error(errInfo)
					err = status.Error(errors.New(errStr), ecode.StatusInternalServerError)
				}
			}()
			return handler(ctx, req)
		}
	}
}

// RateLimit ...
func RateLimit() middleware.Middleware {
	limiter := bbr.NewLimiter()
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			done, e := limiter.Allow()
			if e != nil {
				return nil, status.Error(e, ecode.StatusTooManyRequests)
			}
			reply, err = handler(ctx, req)
			done(ratelimit.DoneInfo{Err: err})
			return
		}
	}
}

// StartAt  ...
func StartAt() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			ctx = context.WithValue(ctx, "startAt", time.Now().UnixMilli())
			return handler(ctx, req)
		}
	}
}

// Log  ...
func Log() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				params = make(map[string]interface{})
				result = make(map[string]interface{})
			)
			if info, ok := transport.FromServerContext(ctx); ok {
				params["method"] = info.Operation()
				headerStr, _ := jsoniter.MarshalToString(info.RequestHeader())
				params["header"] = headerStr
				params["req"] = fmt.Sprint(req)
				params["endpoint"] = info.Endpoint()
			}

			defer func() {
				var level = log.LevelWarn
				st, _struct := status.FromError(err)
				result["grpcCode"] = codes.OK
				if st != nil {
					result["grpcCode"] = st.Code()
					result["error"] = st.Err().Error()
					if _struct != nil {
						for k, v := range _struct.AsMap() {
							result[k] = v
						}
					}
				}

				var processTime int64
				if startAt, _ := ctx.Value("startAt").(int64); startAt > 0 {
					processTime = time.Now().UnixMilli() - startAt
				}
				log.Context(ctx).Log(level,
					"@field", map[string]interface{}{
						"params":      params,
						"processTime": processTime,
						"result":      result,
					},
				)
			}()
			reply, err = handler(ctx, req)
			return
		}
	}
}

// ClientBreaker ...
func ClientBreaker() middleware.Middleware {
	gp := group.NewGroup(func() interface{} {
		return sre.NewBreaker()
	})
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			info, _ := transport.FromClientContext(ctx)
			breaker := gp.Get(info.Operation()).(circuitbreaker.CircuitBreaker)
			if err := breaker.Allow(); err != nil {
				breaker.MarkFailed()
				return nil, status.Error(err, ecode.StatusTemporarilyUnavailable)
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
