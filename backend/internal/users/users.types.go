package users

import (
	"time"

	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
)

// -------------------- User --------------------
type UserResponse struct {
	ID        int32  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type UpdateUserRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}

func userToResponse(u repo.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Time.Format(time.RFC3339),
	}
}

// -------------------- Pagination --------------------
type PaginatedFilters struct {
	Search string `json:"search,omitempty" bson:"search,omitempty"`
}

type PaginatedRequest struct {
	Page     int32            `json:"page,omitempty" bson:"page,omitempty"`
	PageSize int32            `json:"pageSize,omitempty" bson:"pageSize,omitempty"`
	SortBy   string           `json:"sortBy,omitempty" bson:"sortBy,omitempty"`
	SortType string           `json:"sortType,omitempty" bson:"sortType,omitempty"`
	Filters  PaginatedFilters `json:"filters,omitempty" bson:"filters,omitempty"`
}

type PaginatedResponce struct {
	TotalEntries int64          `json:"totalEntries" bson:"totalEntries"`
	Entries      []UserResponse `json:"entries" bson:"entries"`
}
