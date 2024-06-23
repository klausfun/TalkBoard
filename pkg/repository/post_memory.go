package repository

import (
	"TalkBoard/models"
	"errors"
	"time"
)

type PostMemory struct {
	memory *Memory
}

func NewPostMemory(memory *Memory) *PostMemory {
	return &PostMemory{
		memory: memory,
	}
}

func (m *PostMemory) Create(userId int, post models.Post) (int, error) {
	m.memory.mu.Lock()
	defer m.memory.mu.Unlock()

	id := int(time.Now().UnixNano() / 1e9)
	post.Id = id
	post.UserId = userId
	m.memory.Posts[userId] = append(m.memory.Posts[userId], post)

	return id, nil
}

func (m *PostMemory) GetAll() ([]models.Post, error) {
	m.memory.mu.Lock()
	defer m.memory.mu.Unlock()

	var posts []models.Post
	for _, userPosts := range m.memory.Posts {
		for _, post := range userPosts {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func (m *PostMemory) GetByPostId(postId int) (models.Post, error) {
	m.memory.mu.Lock()
	defer m.memory.mu.Unlock()

	for _, userPosts := range m.memory.Posts {
		for _, post := range userPosts {
			if post.Id == postId {
				return post, nil
			}
		}
	}

	return models.Post{}, errors.New("there is no post with this id")
}
