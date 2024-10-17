package usecase

import (
	"context"

	"github.com/masatrio/bookstore-api/utils"
)

type OrderItem struct {
	BookID   int64 `json:"book_id"`
	Quantity int   `json:"quantity"`
}

type CreateOrderInput struct {
	Items []OrderItem `json:"items"`
}

type CreateOrderOutput struct {
	OrderID   int64       `json:"order_id"`
	Items     []OrderItem `json:"items"`
	Status    string      `json:"status"`
	CreatedAt string      `json:"created_at"`
}

type GetOrderOutput struct {
	OrderID   int64       `json:"order_id"`
	Items     []OrderItem `json:"items"`
	Status    string      `json:"status"`
	CreatedAt string      `json:"created_at"`
}

type OrderUseCase interface {
	CreateOrder(ctx context.Context, input CreateOrderInput, userID int64) (*CreateOrderOutput, utils.CustomError)
	GetOrders(ctx context.Context, userID int64, limit, offset int) ([]GetOrderOutput, utils.CustomError)
}
