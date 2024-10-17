package book

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/masatrio/bookstore-api/internal/domain/repository"
	"github.com/masatrio/bookstore-api/internal/domain/usecase"
	"github.com/masatrio/bookstore-api/utils"
)

type bookUseCase struct {
	repo repository.Repository
}

// NewBookUseCase creates a new instance of bookUseCase.
func NewBookUseCase(repo repository.Repository) usecase.BookUseCase {
	return &bookUseCase{
		repo: repo,
	}
}

// ListBooks retrieves a list of books based on the input criteria.
func (b *bookUseCase) ListBooks(ctx context.Context, input usecase.ListBooksInput) (*usecase.ListBooksOutput, utils.CustomError) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "bookUseCase.ListBooks")
	defer span.End()

	filter := repository.BookFilter{
		Title:     input.Title,
		Author:    input.Author,
		MinPrice:  input.MinPrice,
		MaxPrice:  input.MaxPrice,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		Limit:     input.Limit,
		Offset:    input.Offset,
	}

	books, totalCount, err := b.repo.BookRepository().GetFiltered(ctx, filter)
	if err != nil {
		span.RecordError(err)
		return nil, utils.NewCustomSystemError("Database Error")
	}

	return &usecase.ListBooksOutput{
		Books:      convertToUsecaseBooks(books),
		TotalCount: totalCount,
		Limit:      input.Limit,
		Offset:     input.Offset,
	}, nil
}

// CreateBook handles creating a new book.
func (b *bookUseCase) CreateBook(ctx context.Context, input usecase.Book) (*usecase.Book, utils.CustomError) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "bookUseCase.CreateBook")
	defer span.End()

	bookID, err := b.repo.BookRepository().CreateBook(ctx, &repository.Book{
		Title:     input.Title,
		Author:    input.Author,
		Price:     input.Price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		span.RecordError(err)
		return nil, utils.NewCustomSystemError("Database Error")
	}

	return &usecase.Book{
		ID:     bookID,
		Title:  input.Title,
		Author: input.Author,
		Price:  input.Price,
	}, nil
}

// GetBook retrieves a book by its ID.
func (b *bookUseCase) GetBook(ctx context.Context, id int64) (*usecase.Book, utils.CustomError) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "bookUseCase.GetBook")
	defer span.End()

	book, err := b.repo.BookRepository().GetBookByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, utils.NewCustomSystemError("Database Error")
	}
	if book == nil {
		return nil, utils.NewCustomUserError("Book ID Not Found")
	}

	return &usecase.Book{
		ID:     book.ID,
		Title:  book.Title,
		Author: book.Author,
		Price:  book.Price,
	}, nil
}

// convertToUsecaseBooks converts a slice of repository books to usecase books.
func convertToUsecaseBooks(repoBooks []repository.Book) []usecase.Book {
	usecaseBooks := make([]usecase.Book, len(repoBooks))
	for i, book := range repoBooks {
		usecaseBooks[i] = usecase.Book{
			ID:        book.ID,
			Title:     book.Title,
			Author:    book.Author,
			Price:     book.Price,
			CreatedAt: book.CreatedAt,
			UpdatedAt: book.UpdatedAt,
		}
	}
	return usecaseBooks
}
