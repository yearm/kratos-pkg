package status

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/yearm/kratos-pkg/ecodes"
	"github.com/yearm/kratos-pkg/env"
	"github.com/yearm/kratos-pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type Option func(*options)
type options struct {
	message string
	level   log.Level
}

// WithMessage used to set the error message.
func WithMessage(message string) Option {
	return func(o *options) {
		o.message = message
	}
}

// WithLevel used to set the error log level.
func WithLevel(level log.Level) Option {
	return func(o *options) {
		o.level = level
	}
}

// Error create a gRPC status error carrying the error details.
func Error(ctx context.Context, code ecodes.Code, err error, opts ...Option) error {
	opt := options{
		message: code.Message(ctx),
		level:   log.LevelError,
	}
	for _, o := range opts {
		o(&opt)
	}

	err = errors.WrapDepth(err, 3)
	detail := (&ErrorDetail{
		Code:    code,
		Message: opt.message,
		Level:   opt.level,
		Callers: errors.Callers(err),
	}).ToStructPB()
	st, _ := status.New(codes.Code(101), fmt.Sprintf("[%s] %v", env.GetServiceName(), err)).WithDetails(detail)
	return st.Err()
}

// WrapError wrap the existing errors and append the new call location information.
func WrapError(ctx context.Context, err error, msg ...string) error {
	if err == nil {
		return nil
	}
	st, detail, ok := FromError(err)
	if !ok || detail == nil {
		return Error(ctx, ecodes.UnknownError, err)
	}

	errorDetail := ToErrorDetail(detail).WrapCallers(3)
	message := st.Message()
	if len(msg) > 0 && msg[0] != "" {
		message = fmt.Sprintf("%s: %s", msg[0], st.Message())
	}
	afterSt, _ := status.New(st.Code(), message).WithDetails(errorDetail.ToStructPB())
	return afterSt.Err()
}

// FromError extract gRPC status and error detail from error.
func FromError(err error) (*status.Status, *structpb.Struct, bool) {
	if err == nil {
		return nil, nil, false
	}
	st, ok := status.FromError(err)
	if !ok || st == nil {
		return nil, nil, false
	}
	for _, detail := range st.Details() {
		if v, ok := detail.(*structpb.Struct); ok {
			return st, v, true
		}
	}
	return st, nil, true
}

// FromErrorDetail parse the ErrorDetail structure from the error.
func FromErrorDetail(ctx context.Context, err error) *ErrorDetail {
	if err == nil {
		return nil
	}
	_, detail, ok := FromError(err)
	if !ok || detail == nil {
		err = errors.WrapDepth(err, 3)
		return &ErrorDetail{
			Code:    ecodes.UnknownError,
			Message: ecodes.UnknownError.Message(ctx),
			Level:   log.LevelError,
			Callers: errors.Callers(err),
		}
	}
	return ToErrorDetail(detail)
}
