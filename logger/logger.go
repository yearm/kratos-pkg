package logger

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/samber/lo"
	"github.com/yearm/kratos-pkg/config/gconfig"
	"github.com/yearm/kratos-pkg/env"
	"github.com/yearm/kratos-pkg/errors"
	"github.com/yearm/kratos-pkg/logger/aliyun"
	"github.com/yearm/kratos-pkg/logger/logrus"
	"github.com/yearm/kratos-pkg/logger/zap"
	"github.com/yearm/kratos-pkg/utils/debug"
	"github.com/yearm/kratos-pkg/utils/net"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Type is a logger type.
type Type uint8

const (
	Zap Type = iota
	Logrus
	Aliyun
)

// Init initializes the global logger with specified type and contextual metadata.
func Init(typ Type, localTyp ...Type) (func(), error) {
	var (
		logger  log.Logger
		cleanup func()
		err     error
	)
	switch typ {
	case Zap:
		logger, cleanup, err = NewZapLogger()
	case Logrus:
		logger, cleanup, err = NewLogrusLogger()
	case Aliyun:
		logger, cleanup, err = NewAliyunLogger(localTyp...)
	default:
		return nil, errors.Errorf("unknown log type[%v]", typ)
	}
	if err != nil {
		return nil, errors.Wrap(err, "create logger failed")
	}

	ips, err := net.GetIPArray()
	if err != nil {
		return nil, errors.Wrap(err, "net.GetIPArray failed")
	}
	log.NewHelper(log.NewFilter(log.GetLogger(), log.FilterKey("")))

	logger = log.With(logger,
		"system", env.GetServiceName(),
		"version", env.GetServiceVersion(),
		"source", lo.If(len(ips) <= 0, "").ElseF(func() string { return ips[0] }),
		"timestamp", timestampValue(),
		"caller", callerValue(),
		"traceId", tracing.TraceID(),
		"spanId", tracing.SpanID(),
	)
	log.SetLogger(logger)
	return cleanup, nil
}

// callerValue is a Valuer that returns the file and line.
func callerValue() log.Valuer {
	return func(ctx context.Context) any {
		return debug.Caller(5)
	}
}

// timestampValue is a Valuer that returns the current timestamp.
func timestampValue() log.Valuer {
	return func(ctx context.Context) any {
		return time.Now().Unix()
	}
}

// NewZapLogger creates a zap logger with file.
// in local environments, duplicates output to stdout.
func NewZapLogger() (log.Logger, func(), error) {
	c, err := gconfig.GetLogFileConfig()
	if err != nil {
		return nil, nil, errors.Wrap(err, "gconfig.GetLogFileConfig failed")
	}
	if err := os.MkdirAll(filepath.Dir(c.Path), 0755); err != nil {
		return nil, nil, errors.Wrap(err, "os.MkdirAll failed")
	}
	f, err := os.OpenFile(c.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "os.OpenFile failed, path = %v", c.Path)
	}
	ws := []io.Writer{f}
	isLocalMode, err := gconfig.IsLocalMode()
	if err != nil {
		return nil, nil, errors.Wrap(err, "gconfig.IsLocalMode failed")
	}
	if isLocalMode {
		ws = append(ws, os.Stdout)
	}
	w := io.MultiWriter(ws...)
	logger, cleanup := zap.NewLogger(w, zap.ParseLevel(log.ParseLevel(c.Level)))
	return logger, cleanup, nil
}

// NewLogrusLogger creates a logrus logger with file.
// in local environments, duplicates output to stdout.
func NewLogrusLogger() (log.Logger, func(), error) {
	c, err := gconfig.GetLogFileConfig()
	if err != nil {
		return nil, nil, errors.Wrap(err, "gconfig.GetLogFileConfig failed")
	}
	if err := os.MkdirAll(filepath.Dir(c.Path), 0755); err != nil {
		return nil, nil, errors.Wrap(err, "os.MkdirAll failed")
	}
	f, err := os.OpenFile(c.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "os.OpenFile failed, path = %v", c.Path)
	}
	ws := []io.Writer{f}
	isLocalMode, err := gconfig.IsLocalMode()
	if err != nil {
		return nil, nil, errors.Wrap(err, "gconfig.IsLocalMode failed")
	}
	if isLocalMode {
		ws = append(ws, os.Stdout)
	}
	w := io.MultiWriter(ws...)
	logger, cleanup := logrus.NewLogger(w, logrus.ParseLevel(log.ParseLevel(c.Level)))
	return logger, cleanup, nil
}

// NewAliyunLogger creates an aliyun logger.
// in local environments: Uses local logging implementations (Zap/Logrus).
func NewAliyunLogger(localType ...Type) (log.Logger, func(), error) {
	isLocalMode, err := gconfig.IsLocalMode()
	if err != nil {
		return nil, nil, errors.Wrap(err, "gconfig.IsLocalMode failed")
	}
	if isLocalMode && len(localType) > 0 {
		switch localType[0] {
		case Zap:
			return NewZapLogger()
		case Logrus:
			return NewLogrusLogger()
		default:
			return nil, nil, errors.Errorf("unsupported local log type[%v]", localType)
		}
	}
	c, err := gconfig.GetLogAliyunConfig()
	if err != nil {
		return nil, nil, errors.Wrap(err, "gconfig.GetLogAliyunConfig failed")
	}
	ac := aliyun.NewConfig(c.AccessKey, c.SecretKey, c.Endpoint, c.Project, c.Logstore, c.Level)
	return aliyun.NewLogger(ac)
}
