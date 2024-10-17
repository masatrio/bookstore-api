package utils

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestNewTracer(t *testing.T) {
	ctx := context.Background()

	tracer := NewTracer(ctx, "test-service")

	if tracer == nil {
		t.Fatalf("expected tracer to be non-nil, got nil")
	}

	provider := otel.GetTracerProvider()
	if provider == nil {
		t.Fatalf("expected tracer provider to be set, got nil")
	}

	if tp, ok := provider.(*trace.TracerProvider); ok {
		err := tp.Shutdown(ctx)
		if err != nil {
			t.Errorf("failed to shut down tracer provider: %v", err)
		}
	}
}
