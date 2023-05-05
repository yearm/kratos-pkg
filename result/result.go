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
	Errorf string       `json:"-"`
	Level  log.Level    `json:"-"`
}

// SetMessage ...
func (r *Result) SetMessage(msg string) {
	if r == nil {
		return
	}
	r.Msg = msg
}

// SetLevel ...
func (r *Result) SetLevel(level log.Level) {
	if r == nil {
		return
	}
	r.Level = level
}

// New ...
func New(status ecode.Status, data interface{}, levels ...log.Level) *Result {
	return &Result{
		Status: status,
		Msg:    status.Message(),
		Data:   data,
		Errorf: debug.Caller(2, 3),
		Level:  level(levels...),
	}
}

// NewWithMsg use custom msg
func NewWithMsg(status ecode.Status, data interface{}, msg string, levels ...log.Level) *Result {
	return &Result{
		Status: status,
		Msg:    msg,
		Data:   data,
		Errorf: debug.Caller(2, 3),
		Level:  level(levels...),
	}
}

// NewWithErrorf ...
func NewWithErrorf(status ecode.Status, data interface{}, errorf func() string, levels ...log.Level) *Result {
	return &Result{
		Status: status,
		Msg:    status.Message(),
		Data:   data,
		Errorf: errorf(),
		Level:  level(levels...),
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
func ErrorIs(r *Result, status ecode.Status) bool {
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
	errorf := func() string { return debug.Caller(5, 3) }
	switch st.Code() {
	case codes.Canceled:
		return NewWithErrorf(ecode.StatusCancelled, err, errorf, log.LevelWarn)
	case codes.Unknown:
		return NewWithErrorf(ecode.StatusUnknownError, err, errorf)
	case codes.DeadlineExceeded:
		return NewWithErrorf(ecode.StatusRequestTimeout, err, errorf)
	case codes.Internal:
		return NewWithErrorf(ecode.StatusInternalServerError, err, errorf)
	case codes.Unavailable:
		return NewWithErrorf(ecode.StatusTemporarilyUnavailable, err, errorf)
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
			result := NewWithErrorf(ecode.Status(gconv.String(structMap["status"])), err, errorf, log.ParseLevel(gconv.String(structMap["level"])))
			if useErrMsg {
				result.SetMessage(gconv.String(structMap["msg"]))
			}
			return result
		}
		return NewWithErrorf(ecode.StatusInternalServerError, err, errorf)
	default:
		return NewWithErrorf(ecode.StatusInternalServerError, err, errorf)
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

// NilErrorf ...
var NilErrorf = func() string { return "" }
