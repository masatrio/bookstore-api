package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/masatrio/bookstore-api/internal/delivery/http/middleware"
	"go.opentelemetry.io/otel/trace"
)

// BasicHandler applies the necessary middlewares to a public handler.
func BasicHandler(handlerFunc http.HandlerFunc, tracer trace.Tracer) http.Handler {
	return middleware.PanicRecoveryMiddleware(middleware.OTelMiddleware(tracer)(handlerFunc))
}

// ProtectedHandler applies JWT authentication and other middlewares to protected handlers.
func ProtectedHandler(handlerFunc http.HandlerFunc, tracer trace.Tracer) http.Handler {
	return BasicHandler(http.HandlerFunc(middleware.JWTMiddleware(handlerFunc).ServeHTTP), tracer)
}

// InitRoutes initializes the routes for the bookstore service.
func InitRoutes(tracer trace.Tracer) http.Handler {
	r := mux.NewRouter()
	handler := NewHandler()

	// Public routes
	r.HandleFunc("/api/v1/auth/register", BasicHandler(handler.RegisterHandler, tracer).ServeHTTP).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/auth/login", BasicHandler(handler.LoginHandler, tracer).ServeHTTP).Methods(http.MethodPost)

	// Private routes with JWT middleware
	r.HandleFunc("/api/v1/books", ProtectedHandler(handler.BooksHandler, tracer).ServeHTTP).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/orders", ProtectedHandler(handler.GetOrdersHandler, tracer).ServeHTTP).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/orders", ProtectedHandler(handler.CreateOrderHandler, tracer).ServeHTTP).Methods(http.MethodPost)

	return r
}
