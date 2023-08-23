package repository

import (
	"context"

	"github.com/meedoed/auth-rest/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	SetSession(ctx context.Context, userID primitive.ObjectID, session domain.Session) error
}

type Repository struct {
	Users
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Users: NewUsersRepo(db),
	}
}
