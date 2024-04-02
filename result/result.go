package result

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	"github.com/yearm/kratos-pkg/ecode"
	"github.com/yearm/kratos-pkg/errs"
	kstatus "github.com/yearm/kratos-pkg/status"
	"github.com/yearm/kratos-pkg/util/debug"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type RenderTyp string

const (
	JSON RenderTyp = "JSON"
	HTML RenderTyp = "HTML"
	TEXT RenderTyp = "TEXT"
)

type (
	Option func(r *Result)
	Result struct {
		Status    ecode.Status `json:"status"`
		Msg       string       `json:"message"`
		Data      interface{}  `json:"data"`
		renderTyp RenderTyp
		httpCode  int
		caller    string
		level     log.Level
	}
)

// Message msg
func Message(msg string) Option {
	return func(r *Result) {
		r.Msg = msg
	}
}

// HttpCode http code
func HttpCode(code int) Option {
	return func(r *Result) {
		r.httpCode = code
	}
}

// Caller caller
func Caller(caller string) Option {
	return func(r *Result) {
		r.caller = caller
	}
}

// Level log level
func Level(level log.Level) Option {
	return func(r *Result) {
		r.level = level
	}
}

// UseStatusMsg use status msg
func UseStatusMsg() Option {
	return func(r *Result) {
		r.Msg = r.Status.Message()
	}
}

// RenderType render type
func RenderType(renderTyp RenderTyp) Option {
	return func(r *Result) {
		r.renderTyp = renderTyp
	}
}

// HttpCode ...
func (r *Result) HttpCode() int {
	return r.httpCode
}

// RenderType ...
func (r *Result) RenderType() RenderTyp {
	return r.renderTyp
}

// Caller ...
func (r *Result) Caller() string {
	return r.caller
}

// Level ...
func (r *Result) Level() log.Level {
	return r.level
}

// StatusIs ...
func (r *Result) StatusIs(status ecode.Status) bool {
	if r == nil {
		return false
	}
	return r.Status == status
}

// StatusIsNotFound ...
func (r *Result) StatusIsNotFound(status ecode.Status) bool {
	if r == nil {
		return false
	}
	if lo.Contains([]ecode.Status{ecode.StatusNotFound, ecode.StatusRecordNotFound}, status) {
		return true
	}
	return false
}

// New ...
func New(status ecode.Status, data interface{}, opts ...Option) *Result {
	result := &Result{
		Status:    status,
		Msg:       status.Message(),
		Data:      data,
		renderTyp: JSON,
		level:     status.Level(),
	}
	if e, ok := data.(*errs.ValidateError); ok {
		result.Msg = e.Error()
	}
	for _, opt := range opts {
		opt(result)
	}
	if result.caller == "" {
		result.caller = debug.Caller(2, 3)
	}
	return result
}

// FromRPCError ...
func FromRPCError(err error, opts ...Option) *Result {
	status, ok := gstatus.FromError(err)
	if !ok {
		return nil
	}

	var (
		code   ecode.Status
		result *Result
	)
	defer func() {
		for _, opt := range opts {
			opt(result)
		}
	}()

	switch status.Code() {
	case ecode.RPCBusinessError:
		for _, detail := range status.Details() {
			if st, ok := detail.(*structpb.Struct); ok {
				structMap := st.AsMap()
				result = &Result{
					Status:    ecode.Status(gconv.String(structMap[kstatus.DetailStatusKey])),
					Msg:       gconv.String(structMap[kstatus.DetailMessageKey]),
					Data:      err,
					renderTyp: JSON,
					caller:    debug.Caller(2, 3),
					level:     log.ParseLevel(gconv.String(structMap[kstatus.DetailLevelKey])),
				}
				return result
			}
		}
		code = ecode.StatusInternalServerError
	case codes.Canceled:
		code = ecode.StatusCancelled
	case codes.Unknown:
		code = ecode.StatusUnknownError
	case codes.DeadlineExceeded:
		code = ecode.StatusRequestTimeout
	case codes.Internal:
		code = ecode.StatusInternalServerError
	case codes.Unavailable:
		code = ecode.StatusTemporarilyUnavailable
	default:
		code = ecode.StatusInternalServerError
	}
	result = &Result{
		Status:    code,
		Msg:       code.Message(),
		Data:      err,
		renderTyp: JSON,
		caller:    debug.Caller(2, 3),
		level:     code.Level(),
	}
	return result
}

// StatusIs ...
func StatusIs(result *Result, status ecode.Status) bool {
	if result == nil {
		return false
	}
	return result.Status == status
}
