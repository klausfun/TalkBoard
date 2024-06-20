package repository

import (
	"TalkBoard/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(email, password string) (models.User, error)
}

type Post interface{}

type Subscription interface {
}

type Repository struct {
	Authorization
	Post
	Subscription
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
