package repository

import (
	"context"
	"time"
)

type BookRepository interface {
	CreateBook(ctx context.Context, book *Book) (int64, error)
	GetBookByID(ctx context.Context, bookID int64) (*Book, error)
	GetFiltered(ctx context.Context, filter BookFilter) ([]Book, int, error)
}

type Book struct {
	ID        int64
	Title     string
	Author    string
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BookFilter struct {
	Title     string
	Author    string
	MinPrice  float64
	MaxPrice  float64
	StartDate time.Time
	EndDate   time.Time
	Limit     int
	Offset    int
}
