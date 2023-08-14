package repository

import "github.com/meedoed/auth-rest/models"

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository() *Repository {
	return &Repository{}
}
