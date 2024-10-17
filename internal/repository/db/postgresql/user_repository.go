package postgresql

import (
	"context"
	"database/sql"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/masatrio/bookstore-api/internal/domain/repository"
	"github.com/masatrio/bookstore-api/utils"
)

type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new instance of PostgresUserRepository.
func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

// Create inserts a new user into the database and returns the inserted user's ID.
func (r *PostgresUserRepository) Create(ctx context.Context, user *repository.User) (int64, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "PostgresUserRepository.Create")
	defer span.End()

	query := `INSERT INTO users (name, email, password, created_at, updated_at) 
		      VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`

	id, err := utils.ExecContextWithPreparedReturningID(ctx, r.db, query, user.Name, user.Email, user.Password)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to create user")
		return 0, err
	}

	span.SetStatus(codes.Ok, "User created successfully")
	return id, nil
}

// GetByID retrieves a user by their ID.
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int64) (*repository.User, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "PostgresUserRepository.GetByID")
	defer span.End()

	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE id = $1`

	user := &repository.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			span.SetStatus(codes.Ok, "User not found")
			return nil, nil
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get user by ID")
		return nil, err
	}

	span.SetStatus(codes.Ok, "User retrieved successfully")
	return user, nil
}

// GetByEmail retrieves a user by their email address.
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*repository.User, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "PostgresUserRepository.GetByEmail")
	defer span.End()

	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`

	user := &repository.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			span.SetStatus(codes.Ok, "User not found")
			return nil, nil
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to get user by email")
		return nil, err
	}

	span.SetStatus(codes.Ok, "User retrieved successfully")
	return user, nil
}
