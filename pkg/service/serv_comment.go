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

func (s *CommentService) GetByPostId(postId int) ([]models.Comment, error) {
	comments, err := s.repo.GetByPostId(postId)
	if err != nil {
		return nil, err
	}

	return comments, nil
}
