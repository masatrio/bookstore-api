package usecase

import (
	"context"
	"time"

	"github.com/masatrio/bookstore-api/utils"
)

type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListBooksInput struct {
	Title     string    `json:"title,omitempty"`
	Author    string    `json:"author,omitempty"`
	MinPrice  float64   `json:"min_price,omitempty"`
	MaxPrice  float64   `json:"max_price,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Offset    int       `json:"offset,omitempty"`
}

type ListBooksOutput struct {
	Books      []Book `json:"books"`
	TotalCount int    `json:"total_count"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

type BookUseCase interface {
	CreateBook(ctx context.Context, input Book) (*Book, utils.CustomError)
	GetBook(ctx context.Context, id int64) (*Book, utils.CustomError)
	ListBooks(ctx context.Context, input ListBooksInput) (*ListBooksOutput, utils.CustomError)
}
