package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/masatrio/bookstore-api/internal/delivery/http/middleware"
	"github.com/masatrio/bookstore-api/internal/domain/delivery"
	"github.com/masatrio/bookstore-api/internal/domain/usecase"
	"github.com/masatrio/bookstore-api/utils"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Handler struct {
	userUseCase  usecase.UserUseCase
	bookUseCase  usecase.BookUseCase
	orderUseCase usecase.OrderUseCase
}

// NewHandler creates a new HTTP Handler.
func NewHandler(
	userUseCase usecase.UserUseCase,
	bookUseCase usecase.BookUseCase,
	orderUseCase usecase.OrderUseCase,
) delivery.HTTPHandler {
	return &Handler{
		userUseCase:  userUseCase,
		bookUseCase:  bookUseCase,
		orderUseCase: orderUseCase,
	}
}

// RegisterHandler handles user registration.
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.SpanFromContext(r.Context()).TracerProvider().Tracer("").Start(r.Context(), "RegisterHandler")
	defer span.End()

	var input usecase.RegisterInput
	if err := parseAndValidate(r, &input); err != nil {
		span.SetStatus(codes.Error, "Invalid request data")
		errorResponse(w, utils.NewCustomUserError("Invalid request data"))
		return
	}

	if input.Name == "" || input.Email == "" || input.Password == "" {
		span.SetStatus(codes.Error, "Missing required fields")
		errorResponse(w, utils.NewCustomUserError("Name, email, and password are required"))
		return
	}

	output, err := h.userUseCase.Register(ctx, input)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		errorResponse(w, err)
		return
	}

	span.SetStatus(codes.Ok, "User registered successfully")
	jsonResponse(w, http.StatusCreated, output)
}

// LoginHandler handles user login.
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.SpanFromContext(r.Context()).TracerProvider().Tracer("").Start(r.Context(), "LoginHandler")
	defer span.End()

	var input usecase.LoginInput
	if err := parseAndValidate(r, &input); err != nil {
		span.SetStatus(codes.Error, "Invalid request data")
		errorResponse(w, utils.NewCustomUserError("Invalid request data"))
		return
	}

	if input.Email == "" || input.Password == "" {
		span.SetStatus(codes.Error, "Missing required fields")
		errorResponse(w, utils.NewCustomUserError("Email and password are required"))
		return
	}

	output, err := h.userUseCase.Login(ctx, input)
	if err != nil {
		span.SetStatus(codes.Error, "Unauthorized")
		errorResponse(w, err)
		return
	}

	span.SetStatus(codes.Ok, "Login successful")
	jsonResponse(w, http.StatusOK, output)
}

// ListBooksHandler handles listing books with optional filtering.
func (h *Handler) ListBooksHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.SpanFromContext(r.Context()).TracerProvider().Tracer("").Start(r.Context(), "ListBooksHandler")
	defer span.End()

	input := usecase.ListBooksInput{
		Title:     r.URL.Query().Get("title"),
		Author:    r.URL.Query().Get("author"),
		MinPrice:  parseFloatOrDefault(r.URL.Query().Get("min_price"), 0),
		MaxPrice:  parseFloatOrDefault(r.URL.Query().Get("max_price"), 0),
		StartDate: parseDateOrDefault(r.URL.Query().Get("start_date")),
		EndDate:   parseDateOrDefault(r.URL.Query().Get("end_date")),
		Limit:     parseIntOrDefault(r.URL.Query().Get("limit"), 10),
		Offset:    parseIntOrDefault(r.URL.Query().Get("offset"), 0),
	}

	output, err := h.bookUseCase.ListBooks(ctx, input)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		errorResponse(w, err)
		return
	}

	span.SetStatus(codes.Ok, "Books retrieved successfully")
	jsonResponse(w, http.StatusOK, output)
}

// CreateOrderHandler handles creating a new order.
func (h *Handler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.SpanFromContext(r.Context()).TracerProvider().Tracer("").Start(r.Context(), "CreateOrderHandler")
	defer span.End()

	var input usecase.CreateOrderInput
	if err := parseAndValidate(r, &input); err != nil {
		span.SetStatus(codes.Error, "Invalid request data")
		errorResponse(w, utils.NewCustomUserError("Invalid request data"))
		return
	}

	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		span.SetStatus(codes.Error, "User ID not found in context")
		errorResponse(w, utils.NewCustomSystemError("System Error"))
		return
	}

	// Validate the input
	if err := validateCreateOrderInput(input); err != nil {
		span.SetStatus(codes.Error, err.Error())
		errorResponse(w, err)
		return
	}

	output, err := h.orderUseCase.CreateOrder(ctx, input, userID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		errorResponse(w, err)
		return
	}

	span.SetStatus(codes.Ok, "Order created successfully")
	jsonResponse(w, http.StatusCreated, output)
}

// GetOrdersHandler handles retrieving orders with pagination.
func (h *Handler) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.SpanFromContext(r.Context()).TracerProvider().Tracer("").Start(r.Context(), "GetOrdersHandler")
	defer span.End()

	userID, ok := middleware.GetUserIDFromContext(ctx)
	if !ok {
		span.SetStatus(codes.Error, "User ID not found in context")
		errorResponse(w, utils.NewCustomSystemError("System Error"))
		return
	}

	limit := parseIntOrDefault(r.URL.Query().Get("limit"), 10)
	offset := parseIntOrDefault(r.URL.Query().Get("offset"), 0)

	orders, err := h.orderUseCase.GetOrders(ctx, userID, limit, offset)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		errorResponse(w, err)
		return
	}

	span.SetStatus(codes.Ok, "Orders retrieved successfully")
	jsonResponse(w, http.StatusOK, map[string]interface{}{"orders": orders})
}

// HealthCheckHandler handles health check requests.
func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

// parseFloatOrDefault parses float64 or returns default value.
func parseFloatOrDefault(value string, defaultValue float64) float64 {
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return parsed
}

// parseIntOrDefault parses int or returns default value.
func parseIntOrDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}

// parseDateOrDefault parses date or returns zero time.
func parseDateOrDefault(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	date, _ := time.Parse("2006-01-02", value)
	return date
}

// parseAndValidate decodes JSON request body and validates the required fields.
func parseAndValidate(r *http.Request, input interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return utils.NewCustomUserError("invalid request data")
	}
	return nil
}

// jsonResponse writes a JSON response with a given status code.
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// errorResponse writes an error response based on custom errors.
func errorResponse(w http.ResponseWriter, err utils.CustomError) {
	if err.IsUserError() {
		jsonResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err.IsSystemError() {
		jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	jsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}

// validateCreateOrderInput validates the input for creating an order.
func validateCreateOrderInput(input usecase.CreateOrderInput) utils.CustomError {
	if len(input.Items) == 0 {
		return utils.NewCustomUserError("At least one order item is required")
	}

	for _, item := range input.Items {
		if item.BookID <= 0 {
			return utils.NewCustomUserError("Invalid Book ID")
		}
		if item.Quantity <= 0 {
			return utils.NewCustomUserError("Quantity must be greater than zero")
		}
	}

	return nil
}
