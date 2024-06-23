package repository

import (
	"TalkBoard/models"
	"errors"
	"sync"
	"time"
)

type PostMemory struct {
	posts map[int][]models.Post
	mu    sync.RWMutex
}

func NewPostMemory() *PostMemory {
	return &PostMemory{
		posts: make(map[int][]models.Post),
	}
}

func (m *PostMemory) Create(userId int, post models.Post) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := int(time.Now().UnixNano() / 1e6)
	post.Id = id
	post.UserId = userId
	m.posts[userId] = append(m.posts[userId], post)

	return id, nil
}

func (m *PostMemory) GetAll() ([]models.Post, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var posts []models.Post
	for _, userPosts := range m.posts {
		for _, post := range userPosts {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

func (m *PostMemory) GetByPostId(postId int) (models.Post, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, userPosts := range m.posts {
		for _, post := range userPosts {
			if post.Id == postId {
				return post, nil
			}
		}
	}

	return models.Post{}, errors.New("there is no post with this id")
}
