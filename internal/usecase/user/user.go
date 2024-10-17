package user

import (
	"context"
	"time"

	"github.com/masatrio/bookstore-api/config"
	"github.com/masatrio/bookstore-api/internal/domain/repository"
	"github.com/masatrio/bookstore-api/internal/domain/usecase"
	"github.com/masatrio/bookstore-api/utils"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	repo      repository.Repository
	jwtSecret string
	jwtExpiry time.Duration
}

// NewUserUseCase creates a new instance of userUseCase.
func NewUserUseCase(repo repository.Repository, jwtSecret string, jwtExpiry time.Duration) usecase.UserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

// Register handles user registration.
func (u *userUseCase) Register(ctx context.Context, input usecase.RegisterInput) (*usecase.RegisterOutput, utils.CustomError) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "userUseCase.Register")
	defer span.End()

	existingUser, err := u.repo.UserRepository().GetByEmail(ctx, input.Email)
	if err != nil {
		span.RecordError(err)
		return nil, utils.NewCustomSystemError("Database Error")
	}
	if existingUser != nil {
		span.SetStatus(codes.Error, "Email is already registered")
		return nil, utils.NewCustomUserError("email is already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		span.RecordError(err)
		return nil, utils.NewCustomSystemError("System Error")
	}

	userID, err := u.repo.UserRepository().Create(ctx, &repository.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		span.RecordError(err)
		return nil, utils.NewCustomSystemError("Database Error")
	}

	token, err := utils.GenerateJWT(userID, input.Email, config.LoadConfig().JWT.Secret, config.LoadConfig().JWT.Expiry)
	if err != nil {
		span.RecordError(err)
		return nil, utils.NewCustomSystemError("System Error")
	}

	span.SetStatus(codes.Ok, "Registration successful")
	return &usecase.RegisterOutput{
		Token: token,
		User: usecase.User{
			ID:    userID,
			Name:  input.Name,
			Email: input.Email,
		},
	}, nil
}

// Login handles user login.
func (u *userUseCase) Login(ctx context.Context, input usecase.LoginInput) (*usecase.LoginOutput, utils.CustomError) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("").Start(ctx, "userUseCase.Login")
	defer span.End()

	user, err := u.repo.UserRepository().GetByEmail(ctx, input.Email)
	if err != nil {
		span.RecordError(err)
		return nil, utils.NewCustomSystemError("Database Error")
	}

	if user == nil {
		span.SetStatus(codes.Error, "Email does not exist")
		return nil, utils.NewCustomUserError("email not exist")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		span.SetStatus(codes.Error, "Invalid email or password")
		return nil, utils.NewCustomUserError("invalid email or password")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, config.LoadConfig().JWT.Secret, config.LoadConfig().JWT.Expiry)
	if err != nil {
		span.RecordError(err)
		return nil, utils.NewCustomSystemError("System Error")
	}

	span.SetStatus(codes.Ok, "Login successful")
	return &usecase.LoginOutput{
		Token: token,
		User: usecase.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
