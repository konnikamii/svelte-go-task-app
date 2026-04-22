package contact

import (
	"time"

	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
)

type CreateContactRequest struct {
	Email   string `json:"email"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ContactRequestResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func contactRequestToResponse(item repo.ContactRequest) ContactRequestResponse {
	var createdAt time.Time
	if item.CreatedAt.Valid {
		createdAt = item.CreatedAt.Time
	}

	var updatedAt time.Time
	if item.UpdatedAt.Valid {
		updatedAt = item.UpdatedAt.Time
	}

	return ContactRequestResponse{
		ID:        item.ID,
		Email:     item.Email,
		Title:     item.Title,
		Message:   item.Message,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
