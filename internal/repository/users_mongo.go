package repository

import (
	"context"
	"errors"
	"time"

	"github.com/meedoed/auth-rest/internal/domain"
	"github.com/meedoed/auth-rest/pkg/hash"
	"go.mongodb.org/mongo-driver/bson"
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

func (r *UsersRepo) GetById(ctx context.Context, id primitive.ObjectID) (domain.Session, error) {
	var user domain.User
	if err := r.db.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Session{}, domain.ErrUserNotFound
		}

		return domain.Session{}, err
	}

	return user.Session, nil
}

func (r *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User

	tokenHash, err := hash.HashToken(refreshToken)
	if err != nil {
		return domain.User{}, err
	}

	if err := r.db.FindOne(ctx, bson.M{
		"session.refreshToken": tokenHash,
		"session.expiresAt":    bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}
func (r *UsersRepo) SetSession(ctx context.Context, id primitive.ObjectID, session domain.Session) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})

	return err
}
