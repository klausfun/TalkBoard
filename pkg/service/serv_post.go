package service

import (
	"TalkBoard/models"
	"TalkBoard/pkg/repository"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(userId int, post models.Post) (int, error) {
	return s.repo.Create(userId, post)
}
