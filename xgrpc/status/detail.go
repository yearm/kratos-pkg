package status

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/yearm/kratos-pkg/ecodes"
	"github.com/yearm/kratos-pkg/utils/debug"
	"google.golang.org/protobuf/types/known/structpb"
	"strings"
)

// ErrorDetail error detail structure.
type ErrorDetail struct {
	Code    ecodes.Code `json:"code"`
	Message string      `json:"message"`
	Level   log.Level   `json:"level"`
	Callers string      `json:"callers"`
}

// ToMap convert ErrorDetail to the map[string]any
func (e *ErrorDetail) ToMap() map[string]any {
	return map[string]any{
		"code":    uint32(e.Code),
		"message": e.Message,
		"level":   strings.ToLower(e.Level.String()),
		"callers": e.Callers,
	}
}

// ToStructPB convert ErrorDetail to the structpb.Struct.
func (e *ErrorDetail) ToStructPB() *structpb.Struct {
	st, _ := structpb.NewStruct(e.ToMap())
	return st
}

// WrapCallers appends current caller information to the stack trace.
func (e *ErrorDetail) WrapCallers(depth ...int) *ErrorDetail {
	d := 2
	if len(depth) > 0 {
		d = depth[0]
	}
	if e.Callers == "[]" {
		e.Callers = strings.Replace(e.Callers, "[", fmt.Sprintf(`["%s"`, debug.Caller(d)), 1)
	} else {
		e.Callers = strings.Replace(e.Callers, "[", fmt.Sprintf(`["%s", `, debug.Caller(d)), 1)
	}
	return e
}

// ToErrorDetail convert structpb.Struct to the ErrorDetail.
func ToErrorDetail(st *structpb.Struct) *ErrorDetail {
	m := st.AsMap()
	code, _ := m["code"].(float64)
	message, _ := m["message"].(string)
	level, _ := m["level"].(string)
	callers, _ := m["callers"].(string)
	return &ErrorDetail{
		Code:    ecodes.Code(code),
		Message: message,
		Level:   log.ParseLevel(level),
		Callers: callers,
	}
}
