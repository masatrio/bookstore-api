package repository

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) (int64, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
