package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/apperrors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo repo.Queries
}

func NewService(repo repo.Queries) *Service {
	return &Service{repo: repo}
}

type LoginParams struct {
	Email    string
	Password string
}

// Login verifies credentials and returns the user on success.
func (s *Service) Login(ctx context.Context, params LoginParams) (repo.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.User{}, apperrors.Unauthorized("invalid credentials")
		}
		return repo.User{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return repo.User{}, apperrors.Unauthorized("invalid credentials")
	}

	return user, nil
}

// GetMe returns the user record for the given ID.
func (s *Service) GetMe(ctx context.Context, userID int32) (repo.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}
