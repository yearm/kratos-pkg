package slstracer

import (
	"github.com/yearm/kratos-pkg/config/env"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"os"
)

// New ...
func New() (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(env.GetAliyunSlsTraceEndpoint())))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(
			resource.NewSchemaless(
				semconv.ServiceNameKey.String(env.GetServiceName()),
				semconv.ServiceVersionKey.String(env.GetServiceVersion()),
				semconv.ServiceInstanceIDKey.String(env.GetServiceInstanceID()),
				attribute.String("hostname", os.Getenv("HOSTNAME")),
				attribute.String("sls.otel.project", env.GetAliyunSlsTraceProject()),
				attribute.String("sls.otel.instanceid", env.GetAliyunSlsTraceInstanceID()),
				attribute.String("sls.otel.akid", env.GetAliyunAccessKey()),
				attribute.String("sls.otel.aksecret", env.GetAliyunSecretKey()),
			),
		),
	)
	return tp, nil
}
