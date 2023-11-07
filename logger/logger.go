package logger

import (
	"context"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/sirupsen/logrus"
	"github.com/yearm/kratos-pkg/config/env"
	alogger "github.com/yearm/kratos-pkg/logger/aliyun"
	llogger "github.com/yearm/kratos-pkg/logger/logrus"
	"github.com/yearm/kratos-pkg/util/debug"
	"time"
)

// Init ...
func Init() (cleanup func()) {
	var logger klog.Logger
	switch env.GetMode() {
	case env.ModeDevelop, env.ModeTest, env.ModeProduction:
		logger, cleanup = alogger.NewLogger()
	case env.ModeLocal:
		logger, cleanup = llogger.NewLogger()
	default:
		logrus.Panicln("new logger mode error")
	}

	ips, err := gipv4.GetIpArray()
	if err != nil {
		logrus.Panicln("get ip error:", err)
	}
	var source string
	if len(ips) > 0 {
		source = ips[0]
	}
	logger = klog.With(logger,
		"system", env.GetServiceName(),
		"version", env.GetServiceVersion(),
		"source", source,
		"traceId", tracing.TraceID(),
		"spanId", tracing.SpanID(),
		"caller", callerValue(),
		"timestamp", timestampValue(),
		"date", klog.Timestamp("2006-01-02 15:04:05"),
	)
	klog.SetLogger(logger)
	return
}

// callerValue ...
func callerValue() klog.Valuer {
	return func(ctx context.Context) interface{} {
		return debug.Caller(5, 3)
	}
}

// timestampValue ...
func timestampValue() klog.Valuer {
	return func(ctx context.Context) interface{} {
		return time.Now().Unix()
	}
}
