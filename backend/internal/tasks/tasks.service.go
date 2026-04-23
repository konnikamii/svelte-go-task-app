package tasks

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/apperrors"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/authorization"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/middleware"
	"github.com/sirupsen/logrus"
)

type Service struct {
	// repo Repository
	repo repo.Queries
	db   *pgxpool.Pool
}

// NewService creates a new task service
func NewService(repo repo.Queries, db *pgxpool.Pool) *Service {
	return &Service{repo: repo, db: db}
}

// GetTasks retrieves all tasks
func (s *Service) GetTasksPaginated(ctx context.Context, arg PaginatedParams) (PaginatedReturn, error) {
	actorID := middleware.UserIDFromContext(ctx)
	scopeSet, err := s.scopeSet(ctx, actorID, "task", "read")
	if err != nil {
		return PaginatedReturn{}, err
	}
	if !scopeSet.Any && !scopeSet.Own {
		return PaginatedReturn{}, apperrors.Forbidden("insufficient permission")
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

	// Build filters
	if arg.Filters.Search != "" {
		args = append(args, "%"+arg.Filters.Search+"%")

		where = append(where,
			fmt.Sprintf("(CAST(id AS TEXT) ILIKE $%d OR title ILIKE $%d OR description ILIKE $%d)", len(args), len(args), len(args)),
		)
	}

	if arg.Filters.Completed != nil {
		args = append(args, *arg.Filters.Completed)
		where = append(where, fmt.Sprintf("completed = $%d", len(args)))
	}

	if !scopeSet.Any {
		args = append(args, actorID)
		where = append(where, fmt.Sprintf("owner_id = $%d", len(args)))
	}

	// Base query
	baseQuery := `SELECT * FROM tasks`

	if len(where) > 0 {
		baseQuery += " WHERE " + strings.Join(where, " AND ")
	}

	// Sorting whitelist
	sortBy := "created_at"
	switch arg.SortBy {
	case "id", "title":
		sortBy = arg.SortBy
	case "createdAt", "created_at":
		sortBy = "created_at"
	}

	sortDir := "ASC"
	if strings.ToLower(arg.SortType) == "desc" {
		sortDir = "DESC"
	}

	baseQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortDir)

	// Pagination
	args = append(args, arg.PageSize)
	limitPos := len(args)

	args = append(args, offset)
	offsetPos := len(args)

	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", limitPos, offsetPos)
	rows, err := s.db.Query(ctx, baseQuery, args...)
	if err != nil {
		logrus.WithError(err).Error("failed to query tasks")
		return PaginatedReturn{}, err
	}
	defer rows.Close()

	tasks := make([]TaskResponse, 0)

	for rows.Next() {
		var task repo.Task
		if err := rows.Scan(
			&task.ID,
			&task.OwnerID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			logrus.WithError(err).Error("failed to scan task row")
			return PaginatedReturn{}, err
		}
		tasks = append(tasks, taskToResponse(task))
	}

	if err := rows.Err(); err != nil {
		logrus.WithError(err).Error("failed while iterating task rows")
		return PaginatedReturn{}, err
	}

	// Count query
	countQuery := "SELECT COUNT(*) FROM tasks"
	if len(where) > 0 {
		countQuery += " WHERE " + strings.Join(where, " AND ")
	}

	var count int64
	err = s.db.QueryRow(ctx, countQuery, args[:len(args)-2]...).Scan(&count)
	if err != nil {
		logrus.WithError(err).Error("failed to count tasks")
		return PaginatedReturn{}, err
	}

	return PaginatedReturn{
		TotalEntries: count,
		Entries:      tasks,
	}, nil
}

// GetTaskByID retrieves a single task by ID
func (s *Service) GetTaskByID(ctx context.Context, id int64) (TaskResponse, error) {
	actorID := middleware.UserIDFromContext(ctx)
	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TaskResponse{}, apperrors.NotFound("task not found")
		}
		return TaskResponse{}, err
	}

	allowed, err := s.canAccessTask(ctx, actorID, task.OwnerID, "read")
	if err != nil {
		return TaskResponse{}, err
	}
	if !allowed {
		return TaskResponse{}, apperrors.Forbidden("insufficient permission")
	}

	return taskToResponse(task), nil
}

// CreateTask creates a new task
func (s *Service) CreateTask(ctx context.Context, task *CreateTaskRequest) (TaskResponse, error) {
	actorID := middleware.UserIDFromContext(ctx)
	allowed, err := s.canAccessTask(ctx, actorID, actorID, "write")
	if err != nil {
		return TaskResponse{}, err
	}
	if !allowed {
		return TaskResponse{}, apperrors.Forbidden("insufficient permission")
	}

	dueDate, err := parseDueDate(task.DueDate)
	if err != nil {
		return TaskResponse{}, err
	}

	completed := false
	if task.Completed != nil {
		completed = *task.Completed
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return TaskResponse{}, err
	}
	defer tx.Rollback(ctx)
	qtx := s.repo.WithTx(tx)
	created, err := qtx.CreateTask(ctx, repo.CreateTaskParams{
		OwnerID:     actorID,
		Title:       task.Title,
		Description: textValue(task.Description),
		DueDate:     dueDate,
		Completed:   completed,
	})
	if err != nil {
		return TaskResponse{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return TaskResponse{}, err
	}

	return taskToResponse(created), nil
}

// UpdateTask updates an existing task
func (s *Service) UpdateTask(ctx context.Context, id int64, task *UpdateTaskRequest) (TaskResponse, error) {
	actorID := middleware.UserIDFromContext(ctx)
	current, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TaskResponse{}, apperrors.NotFound("task not found")
		}
		return TaskResponse{}, err
	}

	allowed, err := s.canAccessTask(ctx, actorID, current.OwnerID, "write")
	if err != nil {
		return TaskResponse{}, err
	}
	if !allowed {
		return TaskResponse{}, apperrors.Forbidden("insufficient permission")
	}

	dueDate, err := parseDueDate(task.DueDate)
	if err != nil {
		return TaskResponse{}, err
	}

	arg := repo.UpdateTaskParams{
		ID:          id,
		Title:       task.Title,
		Description: textValue(task.Description),
		DueDate:     dueDate,
		Completed:   task.Completed,
	}

	updated, err := s.repo.UpdateTask(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return TaskResponse{}, apperrors.NotFound("task not found")
		}
		return TaskResponse{}, err
	}

	return taskToResponse(updated), nil
}

// DeleteTask deletes a task
func (s *Service) DeleteTask(ctx context.Context, id int64) (int64, error) {
	actorID := middleware.UserIDFromContext(ctx)
	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, apperrors.NotFound("task not found")
		}
		return 0, err
	}

	allowed, err := s.canAccessTask(ctx, actorID, task.OwnerID, "write")
	if err != nil {
		return 0, err
	}
	if !allowed {
		return 0, apperrors.Forbidden("insufficient permission")
	}

	return s.repo.DeleteTask(ctx, id)
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

func (s *Service) canAccessTask(ctx context.Context, actorID, ownerID int32, action string) (bool, error) {
	scopeSet, err := s.scopeSet(ctx, actorID, "task", action)
	if err != nil {
		return false, err
	}

	return scopeSet.Allows(actorID, ownerID), nil
}

func textValue(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{}
	}

	return pgtype.Text{String: *value, Valid: true}
}

func parseDueDate(value *string) (pgtype.Timestamptz, error) {
	if value == nil {
		return pgtype.Timestamptz{}, nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return pgtype.Timestamptz{}, nil
	}

	parsed, err := time.Parse(time.RFC3339, trimmed)
	if err == nil {
		return pgtype.Timestamptz{Time: parsed, Valid: true}, nil
	}

	// Accept date-only payloads (YYYY-MM-DD) from the frontend.
	parsed, err = time.Parse("2006-01-02", trimmed)
	if err == nil {
		return pgtype.Timestamptz{Time: parsed, Valid: true}, nil
	}

	return pgtype.Timestamptz{}, apperrors.BadRequest("invalid dueDate")
}
