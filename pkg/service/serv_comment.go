package service

import (
	"TalkBoard/models"
	"TalkBoard/pkg/repository"
)

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) Create(comment models.Comment) (int, error) {
	return s.repo.Create(comment)
}
