package http

import (
	"net/http"

	"github.com/masatrio/bookstore-api/internal/domain/delivery"
)

type Handler struct{}

// NewHandler creates a new HTTP Handler.
func NewHandler() delivery.HTTPHandler {
	return &Handler{}
}

// RegisterHandler handles user registration.
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User registered"))
}

// LoginHandler handles user login.
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("User logged in"))
}

// BooksHandler handles get books.
func (h *Handler) BooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Books retrieved"))
}

// GetOrdersHandler handles get user orders.
func (h *Handler) GetOrdersHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Orders retrieved"))
}

// CreateOrderHandler handles create order.
func (h *Handler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Orders retrieved"))
}
