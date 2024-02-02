// tracer_manager.go

package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type TracerManager struct {
	tp *sdktrace.TracerProvider
}

func NewTracerManager(address string) (*TracerManager, error) {
	exp, err := NewExporter(address)
	if err != nil {
		return nil, err
	}

	tp := NewTracerProvider(exp)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return &TracerManager{tp: tp}, nil
}

func (tm *TracerManager) Shutdown(ctx context.Context) error {
	return tm.tp.Shutdown(ctx)
}

func NewExporter(address string) (*jaeger.Exporter, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(address)))
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func NewTracerProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("auth-service"),
			semconv.DeploymentEnvironmentKey.String("development"),
		)),
	)
}

func (tm *TracerManager) Tracer() trace.Tracer {
	return tm.tp.Tracer("catalogue-service")
}
