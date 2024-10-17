package repository

import (
	"context"

	"github.com/masatrio/bookstore-api/utils"
)

type Repository interface {
	BookRepository() BookRepository
	OrderRepository() OrderRepository
	OrderItemRepository() OrderItemRepository
	UserRepository() UserRepository
	WithTransaction(TransactionFunc) utils.CustomError
}

type TransactionFunc func(ctx context.Context) utils.CustomError
