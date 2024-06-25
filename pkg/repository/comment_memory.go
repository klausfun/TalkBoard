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

	if _, ok := m.memory.Comments[comment.PostId]; !ok {
		m.memory.Comments[comment.PostId] = make(map[int]models.Comment)
		comment.ParentCommentId = 0
		m.memory.Comments[comment.PostId][comment.Id] = comment

		return comment.Id, nil
	}

	if _, ok := m.memory.Comments[comment.PostId][comment.ParentCommentId]; !ok {
		comment.ParentCommentId = 0
	}
	m.memory.Comments[comment.PostId][comment.Id] = comment

	return comment.Id, nil
}

func (m *CommentMemory) GetByPostId(postId, limit, offset int) ([]models.Comment, error) {
	m.memory.mu.RLock()
	defer m.memory.mu.RUnlock()

	var result []models.Comment

	var findReplies func(parentId int) []models.Comment
	findReplies = func(parentId int) []models.Comment {
		var replies []models.Comment
		for _, comment := range m.memory.Comments[postId] {
			if comment.ParentCommentId == parentId {
				comment.Replies = findReplies(comment.Id)
				replies = append(replies, comment)
			}
		}
		return replies
	}

	for _, comment := range m.memory.Comments[postId] {
		if comment.ParentCommentId == 0 {
			comment.Replies = findReplies(comment.Id)
			result = append(result, comment)
		}
	}

	return result, nil
}
