package repository

import (
	"context"

	"github.com/meedoed/auth-rest/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const usersCollection = "users"

type UsersRepo struct {
	db *mongo.Collection
}

func NewUsersRepo(db *mongo.Database) *UsersRepo {
	return &UsersRepo{
		db: db.Collection(usersCollection),
	}
}

func (r *UsersRepo) Create(ctx context.Context, user domain.User) error {
	return nil
	// TODO:Implement me!
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User

	return user, nil
	// TODO:Implement me!
}
func (r *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User

	return user, nil
	// TODO:Implement me!
}
func (r *UsersRepo) SetSession(ctx context.Context, userID primitive.ObjectID, session domain.Session) error {
	return nil
	// TODO:Implement me!
}
