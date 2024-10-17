package http

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/masatrio/bookstore-api/config"
	"github.com/masatrio/bookstore-api/internal/delivery/http/middleware"
	"github.com/masatrio/bookstore-api/internal/domain/usecase"
	"github.com/masatrio/bookstore-api/internal/repository/db/postgresql"
	"github.com/masatrio/bookstore-api/internal/usecase/book"
	"github.com/masatrio/bookstore-api/internal/usecase/order"
	"github.com/masatrio/bookstore-api/internal/usecase/user"
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

// NewApp initializes the app with the necessary dependencies and starts the server.
func InitAPP(config *config.Config, tracer trace.Tracer) http.Handler {
	db, err := postgresql.NewDatabase(config.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	bookRepo := postgresql.NewPostgresBookRepository(db)
	userRepo := postgresql.NewPostgresUserRepository(db)
	orderRepo := postgresql.NewPostgresOrderRepository(db)
	orderItemRepo := postgresql.NewPostgresOrderItemRepository(db)

	repo := postgresql.NewRepository(db, bookRepo, orderRepo, orderItemRepo, userRepo)

	userUsecase := user.NewUserUseCase(repo, config.JWT.Secret, time.Duration(config.JWT.Expiry)*time.Second)
	bookUsecase := book.NewBookUseCase(repo)
	orderUsecase := order.NewOrderUseCase(repo)

	return InitRoutes(tracer, config, userUsecase, bookUsecase, orderUsecase)
}

// InitRoutes initializes the routes for the bookstore service.
func InitRoutes(
	tracer trace.Tracer,
	config *config.Config,
	userUsecase usecase.UserUseCase,
	bookUsecase usecase.BookUseCase,
	orderUsecase usecase.OrderUseCase,
) http.Handler {
	r := mux.NewRouter()

	handler := NewHandler(userUsecase, bookUsecase, orderUsecase)

	// Public routes
	authRoutes := r.PathPrefix("/api/v1/auth").Subrouter()
	authRoutes.HandleFunc("/register", BasicHandler(handler.RegisterHandler, tracer).ServeHTTP).Methods(http.MethodPost)
	authRoutes.HandleFunc("/login", BasicHandler(handler.LoginHandler, tracer).ServeHTTP).Methods(http.MethodPost)

	// Private routes with JWT middleware
	bookRoutes := r.PathPrefix("/api/v1/books").Subrouter()
	bookRoutes.HandleFunc("", ProtectedHandler(handler.ListBooksHandler, tracer).ServeHTTP).Methods(http.MethodGet)

	orderRoutes := r.PathPrefix("/api/v1/orders").Subrouter()
	orderRoutes.HandleFunc("", ProtectedHandler(handler.GetOrdersHandler, tracer).ServeHTTP).Methods(http.MethodGet)
	orderRoutes.HandleFunc("", ProtectedHandler(handler.CreateOrderHandler, tracer).ServeHTTP).Methods(http.MethodPost)

	// Health check route
	r.HandleFunc("/health", BasicHandler(handler.HealthCheckHandler, tracer).ServeHTTP).Methods(http.MethodGet)

	return r
}
