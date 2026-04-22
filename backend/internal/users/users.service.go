package users

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/apperrors"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/authorization"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/middleware"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo repo.Queries
	db   *pgxpool.Pool
}

// NewService creates a new user service
func NewService(repo repo.Queries, db *pgxpool.Pool) *Service {
	return &Service{repo: repo, db: db}
}

// Retrieves users with filtering, sorting and pagination.
func (s *Service) GetUsersPaginated(ctx context.Context, arg PaginatedRequest) (PaginatedResponce, error) {
	actorID := middleware.UserIDFromContext(ctx)
	scopeSet, err := s.scopeSet(ctx, actorID, "user", "read")
	if err != nil {
		return PaginatedResponce{}, err
	}
	if !scopeSet.Any && !scopeSet.Own {
		return PaginatedResponce{}, apperrors.Forbidden("insufficient permission")
	}

	if arg.Page <= 0 {
		arg.Page = 1
	}
	if arg.PageSize <= 0 {
		arg.PageSize = 10
	}
	offset := (arg.Page - 1) * arg.PageSize

	var (
		args  []any
		where []string
	)

	if arg.Filters.Search != "" {
		args = append(args, "%"+arg.Filters.Search+"%")
		where = append(where,
			fmt.Sprintf("(CAST(id AS TEXT) ILIKE $%d OR username ILIKE $%d OR email ILIKE $%d)", len(args), len(args), len(args)),
		)
	}

	if !scopeSet.Any {
		args = append(args, actorID)
		where = append(where, fmt.Sprintf("id = $%d", len(args)))
	}

	baseQuery := `SELECT * FROM users`
	if len(where) > 0 {
		baseQuery += " WHERE " + strings.Join(where, " AND ")
	}

	sortBy := "created_at"
	switch arg.SortBy {
	case "id", "username", "email", "created_at":
		sortBy = arg.SortBy
	}

	sortDir := "ASC"
	if strings.ToLower(arg.SortType) == "desc" {
		sortDir = "DESC"
	}

	baseQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortDir)

	args = append(args, arg.PageSize)
	limitPos := len(args)

	args = append(args, offset)
	offsetPos := len(args)

	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", limitPos, offsetPos)

	rows, err := s.db.Query(ctx, baseQuery, args...)
	if err != nil {
		logrus.WithError(err).Error("failed to query users")
		return PaginatedResponce{}, err
	}
	defer rows.Close()

	users := make([]UserResponse, 0)
	for rows.Next() {
		var user repo.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			logrus.WithError(err).Error("failed to scan user row")
			return PaginatedResponce{}, err
		}
		users = append(users, userToResponse(user))
	}

	if err := rows.Err(); err != nil {
		logrus.WithError(err).Error("failed while iterating user rows")
		return PaginatedResponce{}, err
	}

	countQuery := "SELECT COUNT(*) FROM users"
	if len(where) > 0 {
		countQuery += " WHERE " + strings.Join(where, " AND ")
	}

	var count int64
	err = s.db.QueryRow(ctx, countQuery, args[:len(args)-2]...).Scan(&count)
	if err != nil {
		logrus.WithError(err).Error("failed to count users")
		return PaginatedResponce{}, err
	}

	return PaginatedResponce{
		TotalEntries: count,
		Entries:      users,
	}, nil
}

// Retrieves a single user by ID
func (s *Service) GetUserByID(ctx context.Context, id int32) (UserResponse, error) {
	actorID := middleware.UserIDFromContext(ctx)
	allowed, err := s.canAccess(ctx, actorID, id, "user", "read")
	if err != nil {
		return UserResponse{}, err
	}
	if !allowed {
		return UserResponse{}, apperrors.Forbidden("insufficient permission")
	}

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return UserResponse{}, err
	}
	return userToResponse(user), nil
}

// Retrieves a user by email
func (s *Service) GetUserByEmail(ctx context.Context, email string) (UserResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return UserResponse{}, err
	}
	return userToResponse(user), nil
}

// Creates a new user
func (s *Service) CreateUser(ctx context.Context, user *repo.CreateUserParams) (UserResponse, error) {
	_, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return UserResponse{}, apperrors.Conflict("email is already taken")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return UserResponse{}, err
	}
	_, err = s.repo.GetUserByUsername(ctx, user.Username)
	if err == nil {
		return UserResponse{}, apperrors.Conflict("username is already taken")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return UserResponse{}, err
	}
	if err := ValidatePassword(user.Password); err != nil {
		return UserResponse{}, err
	}
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return UserResponse{}, err
	}
	defer tx.Rollback(ctx)
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return UserResponse{}, err
	}
	qtx := s.repo.WithTx(tx)
	created, err := qtx.CreateUser(ctx, repo.CreateUserParams{Username: user.Username, Email: user.Email, Password: string(hash)})
	if err != nil {
		return UserResponse{}, err
	}

	assigned, err := qtx.AssignRoleByNameToUser(ctx, repo.AssignRoleByNameToUserParams{
		UserID: created.ID,
		Name:   "user",
	})
	if err != nil {
		return UserResponse{}, err
	}
	if assigned == 0 {
		return UserResponse{}, apperrors.Internal("default role assignment failed")
	}

	if err := tx.Commit(ctx); err != nil {
		return UserResponse{}, err
	}

	return userToResponse(created), nil
}

// Updates an existing user
func (s *Service) UpdateUser(ctx context.Context, id int32, user *UpdateUserRequest) (UserResponse, error) {
	actorID := middleware.UserIDFromContext(ctx)
	allowed, err := s.canAccess(ctx, actorID, id, "user", "write")
	if err != nil {
		return UserResponse{}, err
	}
	if !allowed {
		return UserResponse{}, apperrors.Forbidden("insufficient permission")
	}

	currentUser, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserResponse{}, apperrors.NotFound("user not found")
		}
		return UserResponse{}, err
	}

	if user.Email != "" && user.Email != currentUser.Email {
		return UserResponse{}, apperrors.BadRequest("email cannot be updated")
	}

	passwordToSave := currentUser.Password
	if user.NewPassword != "" {
		if user.OldPassword == "" {
			return UserResponse{}, apperrors.BadRequest("old password is required to set a new password")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(user.OldPassword)); err != nil {
			return UserResponse{}, apperrors.BadRequest("old password is incorrect")
		}
		if err := ValidatePassword(user.NewPassword); err != nil {
			return UserResponse{}, err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.NewPassword), 12)
		if err != nil {
			return UserResponse{}, err
		}
		passwordToSave = string(hash)
	}

	nameToSave := currentUser.Username
	if user.Username != "" {
		nameToSave = user.Username
	}

	arg := repo.UpdateUserParams{
		ID:       id,
		Username: nameToSave,
		Email:    currentUser.Email,
		Password: passwordToSave,
	}

	updated, err := s.repo.UpdateUser(ctx, arg)
	if err != nil {
		return UserResponse{}, err
	}

	return userToResponse(updated), nil
}

// Deletes a user
func (s *Service) DeleteUser(ctx context.Context, id int32) (int64, error) {
	actorID := middleware.UserIDFromContext(ctx)
	allowed, err := s.canAccess(ctx, actorID, id, "user", "write")
	if err != nil {
		return 0, err
	}
	if !allowed {
		return 0, apperrors.Forbidden("insufficient permission")
	}

	return s.repo.DeleteUser(ctx, id)
}

func (s *Service) scopeSet(ctx context.Context, actorID int32, resource, action string) (authorization.ScopeSet, error) {
	scopes, err := s.repo.GetPermissionScopesForUser(ctx, repo.GetPermissionScopesForUserParams{
		UserID:   actorID,
		Resource: resource,
		Action:   action,
	})
	if err != nil {
		return authorization.ScopeSet{}, err
	}

	return authorization.BuildScopeSet(scopes), nil
}

func (s *Service) canAccess(ctx context.Context, actorID, ownerID int32, resource, action string) (bool, error) {
	scopeSet, err := s.scopeSet(ctx, actorID, resource, action)
	if err != nil {
		return false, err
	}

	return scopeSet.Allows(actorID, ownerID), nil
}

// -------------------------- Helpers --------------------------

// Validates input password
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return apperrors.BadRequest("password must be at least 6 characters")
	}
	if len(password) > 60 {
		return apperrors.BadRequest("password must be at most 60 characters")
	}
	hasNumber := false
	for _, c := range password {
		if c >= '0' && c <= '9' {
			hasNumber = true
			break
		}
	}
	if !hasNumber {
		return apperrors.BadRequest("password must contain at least one number")
	}
	return nil
}
