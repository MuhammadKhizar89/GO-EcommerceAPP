package auth

import (
	"context"
	"errors"
	repo "server/internal/adapters/postgresql/sqlc"
	"server/internal/security"

	"github.com/jackc/pgx/v5/pgconn"
)

type Service interface {
	Signup(ctx context.Context, email, password string) (string, error)
	Login(ctx context.Context, email, password string) (string,
		error)
}

type AuthService struct {
	repo repo.Querier
}

func NewAuthService(repo repo.Querier) Service {
	return &AuthService{repo: repo}
}

func (s *AuthService) Signup(ctx context.Context, email, password string) (string, error) {
	if email == "" {
		return "", errors.New("email is required")
	}
	if password == "" {
		return "", errors.New("password is required")
	}

	hash, err := security.HashPassword(password)
	if err != nil {
		return "", err
	}

	user, err := s.repo.CreateUser(ctx, repo.CreateUserParams{
		Email:    email,
		Password: hash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", errors.New("user with this email already exists")
		}
		return "", err
	}

	token, err := security.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email, password string,
) (string, error) {

	if email == "" {
		return "", errors.New("email is required")
	}
	if password == "" {
		return "", errors.New("password is required")
	}
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := security.CheckPassword(user.Password, password); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := security.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
