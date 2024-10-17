package order

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/masatrio/bookstore-api/internal/domain/repository" // Adjust this import based on your repository structure
	"github.com/masatrio/bookstore-api/internal/domain/usecase"
	"github.com/masatrio/bookstore-api/utils"
)

const successOrderStatus string = "success"

type orderUseCase struct {
	repo repository.Repository
}

// NewOrderUseCase creates a new instance of orderUseCase.
func NewOrderUseCase(repo repository.Repository) usecase.OrderUseCase {
	return &orderUseCase{
		repo: repo,
	}
}

// CreateOrder handles order creation.
func (o *orderUseCase) CreateOrder(ctx context.Context, input usecase.CreateOrderInput, userID int64) (*usecase.CreateOrderOutput, utils.CustomError) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "orderUseCase.CreateOrder")
	defer span.End()

	var orderID int64
	err := o.repo.WithTransaction(func(txCtx context.Context) utils.CustomError {

		var err error
		orderID, err = o.repo.OrderRepository().CreateOrder(txCtx, &repository.Order{
			UserID: userID,
			Status: successOrderStatus,
		})

		if err != nil {
			span.RecordError(err)
			return utils.NewCustomSystemError("Database 3 Error")
		}

		for _, item := range input.Items {
			book, err := o.repo.BookRepository().GetBookByID(txCtx, item.BookID)
			if err != nil {
				span.RecordError(err)
				return utils.NewCustomSystemError("Database 2 Error")
			}

			if book == nil {
				return utils.NewCustomUserError("Book ID Not Found")
			}

			if _, err := o.repo.OrderItemRepository().CreateOrderItem(txCtx, &repository.OrderItem{
				OrderID:  orderID,
				BookID:   item.BookID,
				Quantity: item.Quantity,
			}); err != nil {
				span.RecordError(err)
				return utils.NewCustomSystemError("Database 3 Error")
			}
		}

		return nil
	})

	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return &usecase.CreateOrderOutput{
		OrderID:   orderID,
		Items:     input.Items,
		Status:    successOrderStatus,
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

// GetOrders retrieves user orders by userID with pagination.
func (o *orderUseCase) GetOrders(ctx context.Context, userID int64, limit, offset int) ([]usecase.GetOrderOutput, utils.CustomError) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "orderUseCase.GetOrders")
	defer span.End()

	orders, err := o.repo.OrderRepository().GetOrdersByUserID(ctx, userID, limit, offset)
	if err != nil {
		span.RecordError(err)
		return nil, utils.NewCustomSystemError("Database Error")
	}

	var output []usecase.GetOrderOutput
	for _, order := range orders {

		items, err := o.repo.OrderItemRepository().GetOrderItemsByOrderID(ctx, order.ID)
		if err != nil {
			span.RecordError(err)
			return nil, utils.NewCustomSystemError("Database Error")
		}

		var orderItems []usecase.OrderItem
		for _, item := range items {
			orderItems = append(orderItems, usecase.OrderItem{
				BookID:   item.BookID,
				Quantity: item.Quantity,
			})
		}

		output = append(output, usecase.GetOrderOutput{
			OrderID:   order.ID,
			Items:     orderItems,
			Status:    order.Status,
			CreatedAt: order.CreatedAt.Format(time.RFC3339),
		})
	}

	return output, nil
}
