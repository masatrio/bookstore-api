package delivery

import "net/http"

type HTTPHandler interface {
	RegisterHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	ListBooksHandler(w http.ResponseWriter, r *http.Request)
	GetOrdersHandler(w http.ResponseWriter, r *http.Request)
	CreateOrderHandler(w http.ResponseWriter, r *http.Request)
	HealthCheckHandler(w http.ResponseWriter, r *http.Request)
}
