package auth

import (
	"context"
	"errors"
	repo "server/internal/adapters/postgresql/sqlc"
	"server/internal/security"
)

type Service interface {
	Signup(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthService struct {
	repo repo.Queries
}

func NewAuthService(repo repo.Queries) Service {
	return &AuthService{repo: repo}
}

func (s *AuthService) Signup(ctx context.Context, email, password string) error {
	hash, err := security.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.repo.CreateUser(ctx, repo.CreateUserParams{
		Email:    email,
		Password: hash,
	})
	return err
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := security.CheckPassword(user.Password, password); err != nil {
		return "", errors.New("invalid credentials")
	}

	return security.GenerateJWT(user.ID)
}
