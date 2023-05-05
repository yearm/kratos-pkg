package status

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/yearm/kratos-pkg/config/env"
	"github.com/yearm/kratos-pkg/ecode"
	"github.com/yearm/kratos-pkg/util/debug"
	gstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// Error ...
func Error(err error, status ecode.Status, levels ...log.Level) error {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	_struct, _ := structpb.NewStruct(map[string]interface{}{
		"status": status.String(),
		"msg":    status.Message(),
		"errorf": debug.Caller(2, 3),
		"level":  level(levels...).String(),
	})
	st, _ := gstatus.New(ecode.RPCBusinessError, fmt.Sprintf("[%s]%s", env.GetServiceName(), errMsg)).WithDetails(_struct)
	return st.Err()
}

// ErrorWithMsg ...
func ErrorWithMsg(err error, status ecode.Status, msg string, levels ...log.Level) error {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	_struct, _ := structpb.NewStruct(map[string]interface{}{
		"status": status.String(),
		"msg":    msg,
		"errorf": debug.Caller(2, 3),
		"level":  level(levels...).String(),
	})
	st, _ := gstatus.New(ecode.RPCBusinessError, fmt.Sprintf("[%s]%s", env.GetServiceName(), errMsg)).WithDetails(_struct)
	return st.Err()
}

// FromError ...
func FromError(err error) (*gstatus.Status, *structpb.Struct) {
	if err == nil {
		return nil, nil
	}
	if st, ok := gstatus.FromError(err); ok {
		for _, detail := range st.Details() {
			if _struct, ok := detail.(*structpb.Struct); ok {
				return st, _struct
			}
		}
		return st, nil
	}
	return nil, nil
}

// level ...
func level(levels ...log.Level) log.Level {
	level := log.LevelError
	if len(levels) > 0 {
		level = levels[0]
	}
	return level
}
