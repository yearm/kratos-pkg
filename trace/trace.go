package trace

import (
	"github.com/sirupsen/logrus"
	"github.com/yearm/kratos-pkg/trace/tracer"
	"go.opentelemetry.io/otel"
)

func Init() {
	tp, err := tracer.New()
	if err != nil {
		logrus.Panicln("trace init error:", err)
	}
	otel.SetTracerProvider(tp)
}
