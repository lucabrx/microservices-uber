package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/lukabrx/uber-clone/internal/user"
	"golang.org/x/oauth2"
)

type Service struct {
	pasetoMaker          *PasetoMaker
	googleOauth          *oauth2.Config
	googleUserInfoURL    string
	userRepo             *user.MemoryRepository
	refreshTokenRepo     *RefreshTokenRepository
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewService(pasetoMaker *PasetoMaker, googleOauth *oauth2.Config, userRepo *user.MemoryRepository, refreshTokenRepo *RefreshTokenRepository) *Service {
	return &Service{
		pasetoMaker:          pasetoMaker,
		googleOauth:          googleOauth,
		googleUserInfoURL:    "https://www.googleapis.com/oauth2/v2/userinfo",
		userRepo:             userRepo,
		refreshTokenRepo:     refreshTokenRepo,
		accessTokenDuration:  15 * time.Minute,
		refreshTokenDuration: 7 * 24 * time.Hour,
	}
}

func (s *Service) AuthenticateWithGoogle(ctx context.Context, code string) (string, string, *user.User, error) {
	token, err := s.googleOauth.Exchange(ctx, code)
	if err != nil {
		return "", "", nil, errors.New("failed to exchange code for token: " + err.Error())
	}

	client := s.googleOauth.Client(ctx, token)
	resp, err := client.Get(s.googleUserInfoURL)
	if err != nil {
		return "", "", nil, errors.New("failed to get user info from google: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", nil, errors.New("failed to read user info response body")
	}

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return "", "", nil, errors.New("failed to unmarshal user info")
	}

	persistedUser, err := s.userRepo.CreateOrUpdateUser(user.User{Email: userInfo.Email, Name: userInfo.Name})
	if err != nil {
		return "", "", nil, errors.New("failed to save user")
	}

	accessToken, err := s.pasetoMaker.CreateToken(persistedUser.ID, s.accessTokenDuration)
	if err != nil {
		return "", "", nil, errors.New("failed to create access token")
	}

	refreshToken, err := s.pasetoMaker.CreateToken(persistedUser.ID, s.refreshTokenDuration)
	if err != nil {
		return "", "", nil, errors.New("failed to create refresh token")
	}

	err = s.refreshTokenRepo.Store(RefreshToken{
		Token:     refreshToken,
		UserID:    persistedUser.ID,
		ExpiresAt: time.Now().Add(s.refreshTokenDuration),
	})
	if err != nil {
		return "", "", nil, errors.New("failed to store refresh token")
	}

	return accessToken, refreshToken, &persistedUser, nil
}

func (s *Service) VerifyToken(token string) (*Payload, error) {
	return s.pasetoMaker.VerifyToken(token)
}

func (s *Service) RefreshToken(ctx context.Context, oldRefreshToken string) (string, string, error) {
	payload, err := s.pasetoMaker.VerifyToken(oldRefreshToken)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	_, err = s.refreshTokenRepo.Get(oldRefreshToken)
	if err != nil {
		return "", "", err
	}
	s.refreshTokenRepo.Delete(oldRefreshToken)

	newAccessToken, err := s.pasetoMaker.CreateToken(payload.UserID, s.accessTokenDuration)
	if err != nil {
		return "", "", errors.New("failed to create new access token")
	}

	newRefreshToken, err := s.pasetoMaker.CreateToken(payload.UserID, s.refreshTokenDuration)
	if err != nil {
		return "", "", errors.New("failed to create new refresh token")
	}

	err = s.refreshTokenRepo.Store(RefreshToken{
		Token:     newRefreshToken,
		UserID:    payload.UserID,
		ExpiresAt: time.Now().Add(s.refreshTokenDuration),
	})
	if err != nil {
		return "", "", errors.New("failed to store new refresh token")
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *Service) GetUser(ctx context.Context, userID string) (*user.User, error) {
	foundUser, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return &foundUser, nil
}
