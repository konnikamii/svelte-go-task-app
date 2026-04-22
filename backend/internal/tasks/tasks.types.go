package tasks

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
)

type Filters struct {
	Search    string `json:"search,omitempty" bson:"search,omitempty"`
	Completed *bool  `json:"completed,omitempty" bson:"completed,omitempty"`
}

type PaginatedParams struct {
	Page     int32   `json:"page,omitempty" bson:"page,omitempty"`
	PageSize int32   `json:"pageSize,omitempty" bson:"pageSize,omitempty"`
	SortBy   string  `json:"sortBy,omitempty" bson:"sortBy,omitempty"`
	SortType string  `json:"sortType,omitempty" bson:"sortType,omitempty"`
	Filters  Filters `json:"filters,omitempty" bson:"filters,omitempty"`
}

type PaginatedReturn struct {
	TotalEntries int64          `json:"totalEntries" bson:"totalEntries"`
	Entries      []TaskResponse `json:"entries" bson:"entries"`
}

type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	DueDate     *string `json:"dueDate"`
	Completed   *bool   `json:"completed"`
}

type UpdateTaskRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	DueDate     *string `json:"dueDate"`
	Completed   bool    `json:"completed"`
}

type TaskResponse struct {
	ID          int64   `json:"id"`
	OwnerID     int32   `json:"ownerId"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	DueDate     *string `json:"dueDate"`
	Completed   bool    `json:"completed"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

func taskToResponse(task repo.Task) TaskResponse {
	return TaskResponse{
		ID:          task.ID,
		OwnerID:     task.OwnerID,
		Title:       task.Title,
		Description: nullableString(task.Description),
		DueDate:     nullableTime(task.DueDate),
		Completed:   task.Completed,
		CreatedAt:   task.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:   task.UpdatedAt.Time.Format(time.RFC3339),
	}
}

func nullableString(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}

	text := value.String
	return &text
}

func nullableTime(value pgtype.Timestamptz) *string {
	if !value.Valid {
		return nil
	}

	formatted := value.Time.Format(time.RFC3339)
	return &formatted
}
