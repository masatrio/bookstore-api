package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/masatrio/bookstore-api/internal/domain/repository"
	"github.com/masatrio/bookstore-api/utils"
)

type PostgresBookRepository struct {
	db *sql.DB
}

// NewPostgresBookRepository creates a new instance of PostgresBookRepository.
func NewPostgresBookRepository(db *sql.DB) repository.BookRepository {
	return &PostgresBookRepository{
		db: db,
	}
}

// CreateBook inserts a new book into the database and returns the inserted book's ID.
func (r *PostgresBookRepository) CreateBook(ctx context.Context, book *repository.Book) (int64, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "PostgresBookRepository.CreateBook")
	defer span.End()

	query := `INSERT INTO books (title, author, price, created_at, updated_at) 
		      VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`

	id, err := utils.ExecContextWithPreparedReturningID(ctx, r.db, query, book.Title, book.Author, book.Price)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create book")
		return 0, err
	}

	span.SetStatus(codes.Ok, "Book created successfully")
	return id, nil
}

// GetBookByID retrieves a book by its ID.
func (r *PostgresBookRepository) GetBookByID(ctx context.Context, bookID int64) (*repository.Book, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "PostgresBookRepository.GetBookByID")
	defer span.End()

	query := `SELECT id, title, author, price, created_at, updated_at 
			  FROM books 
			  WHERE id = $1`

	book := &repository.Book{}
	err := r.db.QueryRowContext(ctx, query, bookID).Scan(
		&book.ID, &book.Title, &book.Author, &book.Price, &book.CreatedAt, &book.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			span.SetStatus(codes.Ok, "Book not found")
			return nil, nil
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get book by ID")
		return nil, err
	}

	span.SetStatus(codes.Ok, "Book retrieved successfully")
	return book, nil
}

// GetFiltered retrieves books with filters and pagination.
func (r *PostgresBookRepository) GetFiltered(ctx context.Context, filter repository.BookFilter) ([]repository.Book, int, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "PostgresBookRepository.GetFiltered")
	defer span.End()

	var conditions []string
	var params []interface{}
	var paramCounter = 1

	// Dynamic filters based on provided input
	if filter.Title != "" {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", paramCounter))
		params = append(params, "%"+filter.Title+"%")
		paramCounter++
	}
	if filter.Author != "" {
		conditions = append(conditions, fmt.Sprintf("author ILIKE $%d", paramCounter))
		params = append(params, "%"+filter.Author+"%")
		paramCounter++
	}
	if filter.MinPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("price >= $%d", paramCounter))
		params = append(params, filter.MinPrice)
		paramCounter++
	}
	if filter.MaxPrice > 0 {
		conditions = append(conditions, fmt.Sprintf("price <= $%d", paramCounter))
		params = append(params, filter.MaxPrice)
		paramCounter++
	}
	if !filter.StartDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", paramCounter))
		params = append(params, filter.StartDate)
		paramCounter++
	}
	if !filter.EndDate.IsZero() {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", paramCounter))
		params = append(params, filter.EndDate)
		paramCounter++
	}

	query := `SELECT id, title, author, price, created_at, updated_at FROM books`
	countQuery := `SELECT COUNT(*) FROM books`

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	if filter.Limit < 0 {
		err := fmt.Errorf("invalid limit: %d", filter.Limit)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid limit")
		return nil, 0, err
	}
	if filter.Offset < 0 {
		err := fmt.Errorf("invalid offset: %d", filter.Offset)
		span.RecordError(err)
		span.SetStatus(codes.Error, "Invalid offset")
		return nil, 0, err
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCounter, paramCounter+1)
	params = append(params, filter.Limit, filter.Offset)

	var total int
	countArgs := make([]interface{}, 0, paramCounter-1)

	for i := 0; i < len(params); i++ {
		if i < paramCounter-1 {
			countArgs = append(countArgs, params[i])
		}
	}

	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to count books")
		return nil, 0, err
	}

	// Fetch filtered books with limit/offset
	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to fetch filtered books")
		return nil, 0, err
	}
	defer rows.Close()

	var books []repository.Book
	for rows.Next() {
		var book repository.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price, &book.CreatedAt, &book.UpdatedAt); err != nil {
			span.RecordError(err)
			return nil, 0, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		span.RecordError(err)
		return nil, 0, err
	}

	span.SetStatus(codes.Ok, "Filtered books retrieved successfully")
	return books, total, nil
}
