package tasks

import (
	"context"
	"net/http"
	"time"
)

// Task represents a task in the system
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Error types for proper error handling
type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}

var (
	ErrTaskNotFound = Error{Code: http.StatusNotFound, Message: "task not found"}
	ErrInvalidInput = Error{Code: http.StatusBadRequest, Message: "invalid input"}
	ErrInternal     = Error{Code: http.StatusInternalServerError, Message: "internal server error"}
)

// Repository interface for data access abstraction
type Repository interface {
	GetTasks(ctx context.Context) ([]Task, error)
	GetTaskByID(ctx context.Context, id int) (*Task, error)
	CreateTask(ctx context.Context, task *Task) (*Task, error)
	UpdateTask(ctx context.Context, id int, task *Task) (*Task, error)
	DeleteTask(ctx context.Context, id int) error
}
