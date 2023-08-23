package service

import (
	"context"
	"time"

	"github.com/meedoed/auth-rest/internal/repository"
	"github.com/meedoed/auth-rest/pkg/auth"
	"github.com/meedoed/auth-rest/pkg/hash"
)

type UserSignUpInput struct {
	Name     string
	Email    string
	Phone    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
}

type Services struct {
	Users
}

type Deps struct {
	Repos           *repository.Repository
	Hasher          hash.PasswordHasher
	TokenManager    auth.TokenManager
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewService(deps Deps) *Services {

	usersService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL)

	return &Services{
		Users: usersService,
	}
}
