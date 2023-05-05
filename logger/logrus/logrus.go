package llogger

import (
	klogrus "github.com/go-kratos/kratos/contrib/log/logrus/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/sirupsen/logrus"
	"os"
)

// NewLogger ...
func NewLogger() (klog.Logger, func()) {
	logger := logrus.New()

	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp:  true,
		DisableHTMLEscape: true,
		// PrettyPrint:       true,
	})
	return klogrus.NewLogger(logger), func() { logrus.Println("logrus logger graceful close") }
}
