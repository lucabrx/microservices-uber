package user

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

type User struct {
	ID    string
	Email string
	Name  string
}

type MemoryRepository struct {
	users map[string]*User
	mu    sync.RWMutex
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: make(map[string]*User),
	}
}

func (r *MemoryRepository) CreateOrUpdateUser(user User) (User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = uuid.New().String()
	r.users[user.ID] = &user

	return user, nil
}

func (r *MemoryRepository) GetUserByID(id string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return User{}, errors.New("user not found")
	}

	return *user, nil
}
