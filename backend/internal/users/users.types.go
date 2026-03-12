package users

import (
	"context"
	"net/http"
	"time"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Never expose password in JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
	ErrUserNotFound = Error{Code: http.StatusNotFound, Message: "user not found"}
	ErrInvalidInput = Error{Code: http.StatusBadRequest, Message: "invalid input"}
	ErrInternal     = Error{Code: http.StatusInternalServerError, Message: "internal server error"}
	ErrUserExists   = Error{Code: http.StatusBadRequest, Message: "user already exists"}
)

// Repository interface for data access abstraction
type Repository interface {
	GetUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, id int, user *User) (*User, error)
	DeleteUser(ctx context.Context, id int) error
}
