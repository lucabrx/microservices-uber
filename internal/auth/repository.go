package auth

import (
	"errors"
	"sync"
	"time"
)

type RefreshToken struct {
	Token     string
	UserID    string
	ExpiresAt time.Time
}

type RefreshTokenRepository struct {
	tokens map[string]RefreshToken
	mu     sync.RWMutex
}

func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{
		tokens: make(map[string]RefreshToken),
	}
}

func (r *RefreshTokenRepository) Store(token RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[token.Token] = token
	return nil
}

func (r *RefreshTokenRepository) Get(token string) (RefreshToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tokens[token]
	if !ok || time.Now().After(t.ExpiresAt) {
		return RefreshToken{}, errors.New("invalid or expired refresh token")
	}
	return t, nil
}

func (r *RefreshTokenRepository) Delete(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tokens, token)
}
