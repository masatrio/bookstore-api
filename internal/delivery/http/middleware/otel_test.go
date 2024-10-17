package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// TestOTelMiddleware tests the OTelMiddleware function.
func TestOTelMiddleware(t *testing.T) {
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint("localhost:4317"),
		),
	)
	if err != nil {
		t.Fatalf("failed to create exporter: %v", err)
	}
	res, err := resource.New(context.Background(), resource.WithAttributes(
		semconv.ServiceNameKey.String("test-service"),
	))
	if err != nil {
		t.Fatalf("failed to create resource: %v", err)
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)

	tracer := tp.Tracer("test-tracer")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := OTelMiddleware(tracer)(handler)

	req := httptest.NewRequest(http.MethodGet, "http://www.google.com", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if status := w.Result().StatusCode; status != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, status)
	}

}
