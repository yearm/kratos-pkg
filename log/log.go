package log

import (
	"github.com/sirupsen/logrus"
	"github.com/yearm/kratos-pkg/util/debug"
	"runtime"
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			return "", debug.CallerByFrame(frame, 2)
		},
	})
}
