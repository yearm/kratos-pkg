package result

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/yearm/kratos-pkg/ecode"
	"github.com/yearm/kratos-pkg/util/debug"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// Result ...
type Result struct {
	Status ecode.Status `json:"status"`
	Msg    string       `json:"message"`
	Data   interface{}  `json:"data"`
	caller string
	level  log.Level
}

// SetMessage ...
func (r *Result) SetMessage(msg string) {
	if r == nil {
		return
	}
	r.Msg = msg
}

// Caller ...
func (r *Result) Caller() string {
	return r.caller
}

// Level ...
func (r *Result) Level() log.Level {
	return r.level
}

// SetLevel ...
func (r *Result) SetLevel(level log.Level) {
	if r == nil {
		return
	}
	r.level = level
}

// New ...
func New(status ecode.Status, data interface{}, levels ...log.Level) *Result {
	return &Result{
		Status: status,
		Msg:    status.Message(),
		Data:   data,
		caller: debug.Caller(2, 3),
		level:  level(levels...),
	}
}

// NewWithMsg use custom msg
func NewWithMsg(status ecode.Status, data interface{}, msg string, levels ...log.Level) *Result {
	return &Result{
		Status: status,
		Msg:    msg,
		Data:   data,
		caller: debug.Caller(2, 3),
		level:  level(levels...),
	}
}

// NewWithCaller ...
func NewWithCaller(status ecode.Status, data interface{}, caller func() string, levels ...log.Level) *Result {
	return &Result{
		Status: status,
		Msg:    status.Message(),
		Data:   data,
		caller: caller(),
		level:  level(levels...),
	}
}

// NewFromRPCError ...
func NewFromRPCError(err error) *Result {
	return fromError(err, false)
}

// NewWithMsgFromRPCError ...
func NewWithMsgFromRPCError(err error) *Result {
	return fromError(err, true)
}

// ErrorIs ...
// Deprecated: Use StatusIs
func ErrorIs(r *Result, status ecode.Status) bool {
	if r == nil {
		return false
	}
	return r.Status == status
}

// StatusIs ...
func StatusIs(r *Result, status ecode.Status) bool {
	if r == nil {
		return false
	}
	return r.Status == status
}

// fromError ...
func fromError(err error, useErrMsg bool) *Result {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return nil
	}
	caller := func() string { return debug.Caller(5, 3) }
	switch st.Code() {
	case codes.Canceled:
		return NewWithCaller(ecode.StatusCancelled, err, caller, log.LevelWarn)
	case codes.Unknown:
		return NewWithCaller(ecode.StatusUnknownError, err, caller)
	case codes.DeadlineExceeded:
		return NewWithCaller(ecode.StatusRequestTimeout, err, caller)
	case codes.Internal:
		return NewWithCaller(ecode.StatusInternalServerError, err, caller)
	case codes.Unavailable:
		return NewWithCaller(ecode.StatusTemporarilyUnavailable, err, caller)
	case ecode.RPCBusinessError:
		var (
			_struct *structpb.Struct
			_ok     bool
		)
		for _, detail := range st.Details() {
			if _struct, _ok = detail.(*structpb.Struct); _ok {
				break
			}
		}
		if _struct != nil {
			structMap := _struct.AsMap()
			result := NewWithCaller(ecode.Status(gconv.String(structMap["status"])), err, caller, log.ParseLevel(gconv.String(structMap["level"])))
			if useErrMsg {
				result.SetMessage(gconv.String(structMap["msg"]))
			}
			return result
		}
		return NewWithCaller(ecode.StatusInternalServerError, err, caller)
	default:
		return NewWithCaller(ecode.StatusInternalServerError, err, caller)
	}
}

// level ...
func level(levels ...log.Level) log.Level {
	level := log.LevelWarn
	if len(levels) > 0 {
		level = levels[0]
	}
	return level
}

// NilCaller ...
var NilCaller = func() string { return "" }
