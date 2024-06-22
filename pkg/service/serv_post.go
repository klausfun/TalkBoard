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

func (s *PostService) GetAll() ([]models.Post, error) {
	return s.repo.GetAll()
}

func (s *PostService) GetByPostId(postId int) (models.Post, error) {
	post, err := s.repo.GetByPostId(postId)

	return post, err
}
