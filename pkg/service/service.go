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
}

type Subscription interface {
}

type Service struct {
	Authorization
	Post
	Subscription
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Post:          NewPostService(repos.Post),
	}
}
