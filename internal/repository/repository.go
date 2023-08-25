package repository

import (
	"context"

	"github.com/meedoed/auth-rest/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users interface {
	GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error)
	GetById(ctx context.Context, id primitive.ObjectID) (domain.Session, error)
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
