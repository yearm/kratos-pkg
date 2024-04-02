package status

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/samber/lo"
	"github.com/yearm/kratos-pkg/config/env"
	"github.com/yearm/kratos-pkg/ecode"
	"github.com/yearm/kratos-pkg/util/debug"
	"github.com/yearm/kratos-pkg/xerrors"
	gstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	DetailStatusKey  = "status"
	DetailMessageKey = "msg"
	DetailLevelKey   = "level"
	DetailCallersKey = "callers"
)

// Error ...
func Error(err error, status ecode.Status, levels ...log.Level) error {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	errCallers := xerrors.Callers(err)
	callers := make([]interface{}, 0, len(errCallers)+1)
	callers = append(callers, debug.Caller(2, 3))
	callers = append(callers, gconv.Interfaces(errCallers)...)
	level := lo.If(len(levels) == 0, status.Level()).ElseF(func() log.Level { return levels[0] })
	detail, _ := structpb.NewStruct(map[string]interface{}{
		DetailStatusKey:  status.String(),
		DetailMessageKey: status.Message(),
		DetailLevelKey:   level.String(),
		DetailCallersKey: callers,
	})
	st, _ := gstatus.New(ecode.RPCBusinessError, fmt.Sprintf("[%s] %s", env.GetServiceName(), errMsg)).WithDetails(detail)
	return st.Err()
}

// ErrorWithMsg ...
func ErrorWithMsg(err error, status ecode.Status, msg string, levels ...log.Level) error {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	errCallers := xerrors.Callers(err)
	callers := make([]interface{}, 0, len(errCallers)+1)
	callers = append(callers, debug.Caller(2, 3))
	callers = append(callers, gconv.Interfaces(errCallers)...)
	level := lo.If(len(levels) == 0, status.Level()).ElseF(func() log.Level { return levels[0] })
	detail, _ := structpb.NewStruct(map[string]interface{}{
		DetailStatusKey:  status.String(),
		DetailMessageKey: msg,
		DetailLevelKey:   level.String(),
		DetailCallersKey: callers,
	})
	st, _ := gstatus.New(ecode.RPCBusinessError, fmt.Sprintf("[%s] %s", env.GetServiceName(), errMsg)).WithDetails(detail)
	return st.Err()
}

// FromError ...
func FromError(err error) (*gstatus.Status, *structpb.Struct) {
	if err == nil {
		return nil, nil
	}
	if st, ok := gstatus.FromError(err); ok {
		for _, detail := range st.Details() {
			if detailSt, ok := detail.(*structpb.Struct); ok {
				return st, detailSt
			}
		}
		return st, nil
	}
	return nil, nil
}
