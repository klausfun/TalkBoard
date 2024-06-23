package repository

import (
	"TalkBoard/models"
	"sync"
)

type Memory struct {
	Users    map[int]models.User
	Posts    map[int][]models.Post
	Comments map[int][]models.Comment
	mu       sync.RWMutex
}
