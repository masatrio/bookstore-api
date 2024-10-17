package postgresql

import (
	"context"
	"database/sql"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/masatrio/bookstore-api/internal/domain/repository"
	"github.com/masatrio/bookstore-api/utils"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

// NewPostgresOrderRepository creates a new instance of PostgresOrderRepository.
func NewPostgresOrderRepository(db *sql.DB) repository.OrderRepository {
	return &PostgresOrderRepository{
		db: db,
	}
}

// CreateOrder inserts a new order into the database and returns the inserted order's ID.
func (r *PostgresOrderRepository) CreateOrder(ctx context.Context, order *repository.Order) (int64, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "PostgresOrderRepository.CreateOrder")
	defer span.End()

	query := `INSERT INTO orders (user_id, status, created_at, updated_at) 
		      VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`

	id, err := utils.ExecContextWithPreparedReturningID(ctx, r.db, query, order.UserID, order.Status)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create order")
		return 0, err
	}

	span.SetStatus(codes.Ok, "Order created successfully")
	return id, nil
}

// GetOrderByID retrieves an order by its ID.
func (r *PostgresOrderRepository) GetOrderByID(ctx context.Context, orderID int64) (*repository.Order, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "PostgresOrderRepository.GetOrderByID")
	defer span.End()

	query := `SELECT id, user_id, status, created_at, updated_at 
		      FROM orders 
		      WHERE id = $1`

	row := utils.PrepareAndQueryRowContext(ctx, r.db, query, orderID)

	var order repository.Order
	err := row.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			span.SetStatus(codes.Ok, "Order not found")
			return nil, nil
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get order by ID")
		return nil, err
	}

	span.SetStatus(codes.Ok, "Order retrieved successfully")
	return &order, nil
}

// GetOrdersByUserID retrieves orders by user ID with pagination.
func (r *PostgresOrderRepository) GetOrdersByUserID(ctx context.Context, userID int64, limit, offset int) ([]*repository.Order, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "PostgresOrderRepository.GetOrdersByUserID")
	defer span.End()

	query := `SELECT id, user_id, status, created_at, updated_at 
              FROM orders 
              WHERE user_id = $1
              ORDER BY created_at DESC
              LIMIT $2 OFFSET $3`

	rows, err := utils.PrepareAndQueryContext(ctx, r.db, query, userID, limit, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get orders by user ID")
		return nil, err
	}
	defer rows.Close()

	var orders []*repository.Order
	for rows.Next() {
		var order repository.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetStatus(codes.Ok, "Orders retrieved successfully")
	return orders, nil
}
