package repository

import (
	"TalkBoard/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(email, password string) (models.User, error)
}

type Post interface {
	Create(userId int, post models.Post) (int, error)
	GetAll() ([]models.Post, error)
	GetByPostId(postId int) (models.Post, error)
}

type Comment interface {
	Create(comment models.Comment) (int, error)
	GetByPostId(postId int) ([]models.Comment, error)
}

type Repository struct {
	Authorization
	Post
	Comment
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Post:          NewPostPostgres(db),
		Comment:       NewCommentPostgres(db),
	}
}
