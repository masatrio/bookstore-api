package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/masatrio/bookstore-api/internal/domain/repository"
	"github.com/masatrio/bookstore-api/utils"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type RepositoryImpl struct {
	bookRepo      repository.BookRepository
	orderRepo     repository.OrderRepository
	orderItemRepo repository.OrderItemRepository
	userRepo      repository.UserRepository
	db            *sql.DB
}

// NewRepository creates a new instance of RepositoryImpl.
func NewRepository(
	db *sql.DB,
	bookRepo repository.BookRepository,
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	userRepo repository.UserRepository,
) repository.Repository {
	return &RepositoryImpl{
		bookRepo:      bookRepo,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		userRepo:      userRepo,
		db:            db,
	}
}

// BookRepository returns the BookRepository instance.
func (r *RepositoryImpl) BookRepository() repository.BookRepository {
	return r.bookRepo
}

// OrderRepository returns the OrderRepository instance.
func (r *RepositoryImpl) OrderRepository() repository.OrderRepository {
	return r.orderRepo
}

// OrderItemRepository returns the OrderItemRepository instance.
func (r *RepositoryImpl) OrderItemRepository() repository.OrderItemRepository {
	return r.orderItemRepo
}

// UserRepository returns the UserRepository instance.
func (r *RepositoryImpl) UserRepository() repository.UserRepository {
	return r.userRepo
}

// WithTransaction wraps the database operation in a transaction.
func (r *RepositoryImpl) WithTransaction(fn repository.TransactionFunc) utils.CustomError {
	ctx, span := trace.SpanFromContext(context.Background()).TracerProvider().Tracer("").Start(context.Background(), "PostgresUserRepository.WithTransaction")
	defer span.End()

	tx, err := r.db.Begin()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to begin transaction")
		return utils.NewCustomSystemError(err.Error())
	}

	ctx = context.WithValue(ctx, utils.TransactionContextKey, tx)

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("rollback failed: %v; original error: %w", rollbackErr, err)
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, "Transaction rolled back")
		}
	}()

	funcErr := fn(ctx)
	if funcErr != nil {
		err = funcErr
		span.RecordError(err)
		span.SetStatus(codes.Error, "Transaction function error")
		return funcErr
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to commit transaction")
		return utils.NewCustomSystemError(err.Error())
	}

	span.SetStatus(codes.Ok, "Transaction committed successfully")
	return nil
}
