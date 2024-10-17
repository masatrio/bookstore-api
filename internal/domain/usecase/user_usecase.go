package usecase

import (
	"context"

	"github.com/masatrio/bookstore-api/utils"
)

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterOutput struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UserUseCase interface {
	Register(ctx context.Context, input RegisterInput) (*RegisterOutput, utils.CustomError)
	Login(ctx context.Context, input LoginInput) (*LoginOutput, utils.CustomError)
}
