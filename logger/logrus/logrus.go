package logrus

import (
	klogrus "github.com/go-kratos/kratos/contrib/log/logrus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

// NewLogger creates a logrus logger.
func NewLogger(w io.Writer, level logrus.Level) (log.Logger, func()) {
	l := logrus.New()
	l.SetLevel(level)
	l.SetOutput(w)
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat:   time.RFC3339,
		DisableHTMLEscape: true,
	})
	logger := klogrus.NewLogger(l)
	return logger, func() {}
}

// ParseLevel parse level to logrus level.
func ParseLevel(level log.Level) logrus.Level {
	switch level {
	case log.LevelDebug:
		return logrus.DebugLevel
	case log.LevelInfo:
		return logrus.InfoLevel
	case log.LevelWarn:
		return logrus.WarnLevel
	case log.LevelError:
		return logrus.ErrorLevel
	case log.LevelFatal:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
