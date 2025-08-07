package trace

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hoisie/mustache"
	"github.com/yearm/kratos-pkg/config/gconfig"
	"github.com/yearm/kratos-pkg/env"
	"github.com/yearm/kratos-pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Init initializes the OpenTelemetry tracer provider.
// if a trace endpoint is configured, it sets up an OTLP HTTP exporter.
func Init(opts ...tracesdk.TracerProviderOption) error {
	baseOptions := []tracesdk.TracerProviderOption{
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.AlwaysSample())),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceInstanceID(env.GetServiceID()),
			semconv.ServiceName(env.GetServiceName()),
			semconv.ServiceVersion(env.GetServiceVersion()),
		)),
	}

	c, _ := gconfig.GetTraceConfig()
	if c != nil && c.Exporter != nil && c.Exporter.Endpoint != "" {
		mode, err := gconfig.GetMode()
		if err != nil {
			return errors.Wrap(err, "gconfig.GetMode failed")
		}
		c.Exporter.Endpoint = mustache.Render(c.Exporter.Endpoint, map[string]string{"mode": mode.String()})
		exporter, err := otlptracehttp.New(context.Background(),
			otlptracehttp.WithEndpointURL(c.Exporter.Endpoint),
			otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
		)
		if err != nil {
			return errors.Wrap(err, "otlptracehttp.New failed")
		}
		baseOptions = append(baseOptions, tracesdk.WithBatcher(exporter))
	}

	options := append(baseOptions, opts...)
	tp := tracesdk.NewTracerProvider(options...)
	otel.SetTracerProvider(tp)
	otel.SetErrorHandler(&errorHandler{})
	return nil
}

type errorHandler struct{}

func (e *errorHandler) Handle(err error) {
	log.Warnf("otel handler error: %v", err)
}
