package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
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
	return &Repository{}
}
