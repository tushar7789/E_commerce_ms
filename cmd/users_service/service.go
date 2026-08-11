package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	db "github.com/tushar7789/E_commerce_ms/cmd/users_service/data/db/out"
)

type Service interface {
	CreateUser(ctx context.Context, user db.CreateUserParams) (interface{}, TokenPair, error)
	Login(ctx context.Context, username, password string) (TokenPair, error)
	Refresh(refreshToken string) (TokenPair, error)
}

type svc struct {
	repo          db.Querier
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

func NewService(repo db.Querier, access_secret string, refresh_secret string, accessTTL time.Duration, refreshTTL time.Duration) Service {
	return &svc{
		repo:          repo,
		AccessSecret:  access_secret,
		RefreshSecret: refresh_secret,
		AccessTTL:     accessTTL,
		RefreshTTL:    refreshTTL,
	}
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func (s *svc) CreateUser(ctx context.Context, user db.CreateUserParams) (interface{}, TokenPair, error) {
	newCreatedUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		log.Print(err)
		return nil, TokenPair{}, nil
	}

	newUser, err := s.repo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return nil, TokenPair{}, ErrInvalidCredentials
	}

	userId := strconv.FormatInt(newUser.ID, 10)

	Tokens, err := GetTokens(userId, s.AccessSecret, s.AccessTTL, s.RefreshSecret, s.RefreshTTL)
	if err != nil {
		return nil, TokenPair{}, err
	}

	return newCreatedUser, Tokens, nil
}

func (s *svc) Login(ctx context.Context, username, password string) (TokenPair, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return TokenPair{}, ErrInvalidCredentials
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return TokenPair{}, ErrInvalidCredentials
	}

	userId := strconv.FormatInt(user.ID, 10)

	tokens, err := GetTokens(userId, s.AccessSecret, s.AccessTTL, s.RefreshSecret, s.RefreshTTL)
	if err != nil {
		return TokenPair{}, err
	}

	return tokens, nil
}

func (s *svc) Refresh(refreshToken string) (TokenPair, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &Claims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.RefreshSecret), nil
	})

	if err != nil || !token.Valid {
		return TokenPair{}, ErrInvalidCredentials
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return TokenPair{}, ErrInvalidCredentials
	}

	tokens, err := GetTokens(claims.UserID, s.AccessSecret, s.AccessTTL, s.RefreshSecret, s.RefreshTTL)
	if err != nil {
		return TokenPair{}, err
	}

	return tokens, nil
}
