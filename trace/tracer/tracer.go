package tracer

import (
	"github.com/yearm/kratos-pkg/config/env"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.14.0"
	"os"
	"strings"
)

// New ...
func New() (*tracesdk.TracerProvider, error) {
	traceUrl := strings.ReplaceAll(env.GetAliyunTraceEndpoint(), "{{env}}", env.GetMode())
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(traceUrl)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(env.GetServiceName()),
			semconv.ServiceVersionKey.String(env.GetServiceVersion()),
			semconv.ServiceInstanceIDKey.String(env.GetServiceInstanceID()),
			attribute.String("hostname", os.Getenv("HOSTNAME")),
		)),
	)
	return tp, nil
}
