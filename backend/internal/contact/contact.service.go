package contact

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/apperrors"
)

type Service struct {
	repo repo.Queries
	db   *pgxpool.Pool
	mail MailConfig
}

func NewService(repo repo.Queries, db *pgxpool.Pool, mail MailConfig) *Service {
	return &Service{repo: repo, db: db, mail: mail}
}

func (s *Service) CreateContactRequest(ctx context.Context, input *CreateContactRequest) (ContactRequestResponse, error) {
	email := strings.TrimSpace(input.Email)
	title := strings.TrimSpace(input.Title)
	message := strings.TrimSpace(input.Message)

	if email == "" || title == "" || message == "" {
		return ContactRequestResponse{}, apperrors.BadRequest("email, title, and message are required")
	}

	if !strings.Contains(email, "@") {
		return ContactRequestResponse{}, apperrors.BadRequest("invalid email address")
	}

	created, err := s.repo.CreateContactRequest(ctx, repo.CreateContactRequestParams{
		Email:   email,
		Title:   title,
		Message: message,
	})
	if err != nil {
		return ContactRequestResponse{}, err
	}

	response := contactRequestToResponse(created)
	if err := s.mail.SendContactNotification(response); err != nil {
		logContactMailFailure(err, response)
	}

	return response, nil
}
