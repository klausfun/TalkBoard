package repository

import (
	"TalkBoard/models"
	"errors"
	"time"
)

type UserMemory struct {
	memory *Memory
}

func NewUserMemory(memory *Memory) *UserMemory {
	return &UserMemory{
		memory: memory,
	}
}

func (m *UserMemory) CreateUser(user models.User) (int, error) {
	m.memory.mu.Lock()
	defer m.memory.mu.Unlock()

	id := int(time.Now().UnixNano() / 1e9)
	user.Id = id
	m.memory.Users[id] = user

	return id, nil
}

func (m *UserMemory) GetUser(email, password string) (models.User, error) {
	var user models.User
	for _, user := range m.memory.Users {
		if user.Email == email && user.Password == password {
			return user, nil
		}
	}

	return user, errors.New("no such user has been found")
}
