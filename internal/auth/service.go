package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"golang.org/x/oauth2"
)

type Service struct {
	pasetoMaker       *PasetoMaker
	googleOauth       *oauth2.Config
	googleUserInfoURL string
}

func NewService(pasetoMaker *PasetoMaker, googleOauth *oauth2.Config) *Service {
	return &Service{
		pasetoMaker:       pasetoMaker,
		googleOauth:       googleOauth,
		googleUserInfoURL: "https://www.googleapis.com/oauth2/v2/userinfo",
	}
}

func (s *Service) AuthenticateWithGoogle(ctx context.Context, code string) (string, error) {
	token, err := s.googleOauth.Exchange(ctx, code)
	if err != nil {
		return "", errors.New("failed to exchange code for token: " + err.Error())
	}

	client := s.googleOauth.Client(ctx, token)
	resp, err := client.Get(s.googleUserInfoURL)
	if err != nil {
		return "", errors.New("failed to get user info from google: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("failed to read user info response body")
	}

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return "", errors.New("failed to unmarshal user info")
	}

	userID := userInfo.ID

	pasetoToken, err := s.pasetoMaker.CreateToken(userID, 24*time.Hour)
	if err != nil {
		return "", errors.New("failed to create paseto token")
	}

	return pasetoToken, nil
}

func (s *Service) VerifyToken(token string) (*Payload, error) {
	return s.pasetoMaker.VerifyToken(token)
}
