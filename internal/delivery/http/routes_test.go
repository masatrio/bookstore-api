package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/masatrio/bookstore-api/config"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

func TestInitAPP(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{},
		JWT: config.JWTConfig{
			Secret: "testsecret",
			Expiry: 3600,
		},
	}

	tracer := trace.NewNoopTracerProvider().Tracer("test")

	handler := InitAPP(cfg, tracer)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("failed to make a request: %v", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
