package repository

import (
	"context"
	"time"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *Order) (int64, error)
	GetOrderByID(ctx context.Context, orderID int64) (*Order, error)
	GetOrdersByUserID(ctx context.Context, userID int64, limit, offset int) ([]*Order, error)
}

type OrderItemRepository interface {
	CreateOrderItem(ctx context.Context, orderItem *OrderItem) (int64, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]*OrderItem, error)
}

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderItem struct {
	ID       int64 `json:"id"`
	OrderID  int64 `json:"order_id"`
	BookID   int64 `json:"book_id"`
	Quantity int   `json:"quantity"`
}
