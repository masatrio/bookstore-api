package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPanicRecoveryMiddleware tests the PanicRecoveryMiddleware for handling panics.
func TestPanicRecoveryMiddleware(t *testing.T) {
	panickingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	recoveryMiddleware := PanicRecoveryMiddleware(panickingHandler)

	req := httptest.NewRequest(http.MethodGet, "http://bookstore", nil)
	rr := httptest.NewRecorder()

	recoveryMiddleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	assert.Contains(t, rr.Body.String(), "Internal Server Error")
}
