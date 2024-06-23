package repository

import (
	"TalkBoard/models"
	"errors"
	"time"
)

type CommentMemory struct {
	memory *Memory
}

func NewCommentMemory(memory *Memory) *CommentMemory {
	return &CommentMemory{
		memory: memory,
	}
}

func (m *CommentMemory) Create(comment models.Comment) (int, error) {
	m.memory.mu.Lock()
	defer m.memory.mu.Unlock()

	id := int(time.Now().UnixNano() / 1e9)
	comment.Id = id

	accessToComments, flag := false, false
	for _, userPosts := range m.memory.Posts {
		for _, post := range userPosts {
			if post.Id == comment.PostId {
				accessToComments = post.AccessToComments
				flag = true
			}
		}
	}
	if !flag {
		return 0, errors.New("there is no post with this id!")
	}
	if !accessToComments {
		return 0, errors.New("no comments are available under this post!")
	}

	comments, ok := m.memory.Comments[comment.ParentCommentId]
	if !ok {
		m.memory.Comments[0] = append(m.memory.Comments[0], comments...)
		return comment.Id, nil
	}

	flag = false
	for _, com := range comments {
		if com.PostId == comment.PostId &&
			com.ParentCommentId == comment.ParentCommentId {
			flag = true
			break
		}
	}
	if !flag {
		return 0, errors.New("postId and commentId do not match each other!")
	}

	m.memory.Comments[comment.ParentCommentId] = append(m.memory.Comments[comment.ParentCommentId], comments...)
	return comment.Id, nil
}

func (m *CommentMemory) GetByPostId(postId, limit, offset int) ([]models.Comment, error) {
	return []models.Comment{}, nil
}
