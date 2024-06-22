package service

import (
	"TalkBoard/models"
	"TalkBoard/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(email, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Post interface {
	Create(userId int, post models.Post) (int, error)
	GetAll() ([]models.Post, error)
	GetByPostId(postId int) (models.Post, error)
}

type Comment interface {
	//Create(userId int, comment models.Comment) (int, error)
	Create(comment models.Comment) (int, error)
	GetByPostId(postId int) ([]models.Comment, error)
}

type Service struct {
	Authorization
	Post
	Comment
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Post:          NewPostService(repos.Post),
		Comment:       NewCommentService(repos.Comment),
	}
}
