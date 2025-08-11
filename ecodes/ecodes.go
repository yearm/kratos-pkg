package ecodes

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
)

// Code define the types of business error codes.
type Code uint32

const (
	OK Code = iota

	Canceled

	UnknownError

	NotImplemented

	ServiceUnavailable

	InternalServerError

	TooManyRequests

	RequestTimeout

	BadRequest

	Conflict

	InvalidArgument

	NotFound

	AccessDenied

	Unauthorized
)

// Is check whether the target code is equal to the code.
func (c Code) Is(code Code) bool {
	return c == code
}

// CodeDetail define the detailed information structure of the error code.
type CodeDetail struct {
	Message string
	// LocalizeConfig: reserved field for international multilingual support.
	// LocalizeConfig *i18n.LocalizeConfig
}

// codeMap global error code registry.
var codeMap = map[Code]CodeDetail{
	OK:                  {Message: "成功"},
	Canceled:            {Message: "操作已取消"},
	UnknownError:        {Message: "未知错误"},
	NotImplemented:      {Message: "此接口未实现"},
	ServiceUnavailable:  {Message: "服务暂时不可用"},
	InternalServerError: {Message: "服务器内部错误"},
	TooManyRequests:     {Message: "请求过于频繁"},
	RequestTimeout:      {Message: "请求超时"},
	BadRequest:          {Message: "错误请求"},
	Conflict:            {Message: "资源冲突"},
	InvalidArgument:     {Message: "无效的参数"},
	NotFound:            {Message: "资源不存在"},
	AccessDenied:        {Message: "拒绝访问"},
	Unauthorized:        {Message: "未经授权"},
}

var mu sync.Mutex

// RegisterErrorCodes register error codes, start from iota+1001、iota+2001...
func RegisterErrorCodes(cm map[Code]CodeDetail) {
	mu.Lock()
	defer mu.Unlock()
	for code, detail := range cm {
		if _, ok := codeMap[code]; ok {
			panic("duplicate register error code")
		}
		codeMap[code] = detail
	}
}

// Message error code description information.
// ctx is a reserved parameter for subsequent support of dynamic localization.
func (c Code) Message(ctx context.Context) string {
	cd, ok := codeMap[c]
	if !ok {
		log.Context(ctx).Errorf("unregistered error code[%v] accessed", c)
		return ""
	}
	return cd.Message
}

// FromGRPCCode converts a gRPC error code into the corresponding ecodes.Code.
func FromGRPCCode(code codes.Code) Code {
	switch code {
	case codes.OK:
		return OK
	case codes.Canceled:
		return Canceled
	case codes.Unknown:
		return UnknownError
	case codes.InvalidArgument:
		return InvalidArgument
	case codes.DeadlineExceeded:
		return RequestTimeout
	case codes.NotFound:
		return NotFound
	case codes.AlreadyExists:
		return Conflict
	case codes.PermissionDenied:
		return AccessDenied
	case codes.ResourceExhausted:
		return TooManyRequests
	case codes.FailedPrecondition:
		return BadRequest
	case codes.Aborted:
		return Conflict
	case codes.OutOfRange:
		return BadRequest
	case codes.Unimplemented:
		return NotImplemented
	case codes.Internal:
		return InternalServerError
	case codes.Unavailable:
		return ServiceUnavailable
	case codes.DataLoss:
		return InternalServerError
	case codes.Unauthenticated:
		return Unauthorized
	}
	return InternalServerError
}
