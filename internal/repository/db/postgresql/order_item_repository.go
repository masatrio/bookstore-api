package postgresql

import (
	"context"
	"database/sql"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/masatrio/bookstore-api/internal/domain/repository"
	"github.com/masatrio/bookstore-api/utils"
)

type PostgresOrderItemRepository struct {
	db *sql.DB
}

// NewPostgresOrderItemRepository creates a new instance of PostgresOrderItemRepository.
func NewPostgresOrderItemRepository(db *sql.DB) repository.OrderItemRepository {
	return &PostgresOrderItemRepository{
		db: db,
	}
}

// CreateOrderItem inserts a new order item into the database and returns the inserted item's ID.
func (r *PostgresOrderItemRepository) CreateOrderItem(ctx context.Context, orderItem *repository.OrderItem) (int64, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "PostgresOrderItemRepository.CreateOrderItem")
	defer span.End()

	query := `INSERT INTO order_items (order_id, book_id, quantity) 
		      VALUES ($1, $2, $3) RETURNING id`

	id, err := utils.ExecContextWithPreparedReturningID(ctx, r.db, query, orderItem.OrderID, orderItem.BookID, orderItem.Quantity)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create order item")
		return 0, err
	}

	span.SetStatus(codes.Ok, "Order item created successfully")
	return id, nil
}

// GetOrderItemsByOrderID retrieves all order items for a specific order.
func (r *PostgresOrderItemRepository) GetOrderItemsByOrderID(ctx context.Context, orderID int64) ([]*repository.OrderItem, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "PostgresOrderItemRepository.GetOrderItemsByOrderID")
	defer span.End()

	query := `SELECT id, order_id, book_id, quantity 
		      FROM order_items 
		      WHERE order_id = $1`

	rows, err := utils.PrepareAndQueryContext(ctx, r.db, query, orderID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get order items by order ID")
		return nil, err
	}
	defer rows.Close()

	var orderItems []*repository.OrderItem
	for rows.Next() {
		var orderItem repository.OrderItem
		err := rows.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.BookID, &orderItem.Quantity)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		orderItems = append(orderItems, &orderItem)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "Order items retrieved successfully")
	return orderItems, nil
}
