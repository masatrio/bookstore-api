package utils

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	tracer "go.opentelemetry.io/otel/trace"
)

// NewTracer creates and returns a new OpenTelemetry tracer.
func NewTracer(ctx context.Context, serviceName string) tracer.Tracer {
	exporter, err := otlptracehttp.New(ctx)
	if err != nil {
		log.Fatalf("failed to create OTLP HTTP trace exporter: %v", err)
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return tp.Tracer(serviceName)
}
