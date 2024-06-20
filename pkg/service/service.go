package service

import "TalkBoard/pkg/repository"

type Authorization interface {
}

type Post interface{}

type Subscription interface {
}

type Service struct {
	Authorization
	Post
	Subscription
}

func NewService(repos *repository.Repository) *Service {
	return &Service{}
}
