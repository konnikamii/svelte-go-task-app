package tasks

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/apperrors"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/authorization"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/middleware"
	"github.com/sirupsen/logrus"
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
	TotalEntries int64       `json:"totalEntries" bson:"totalEntries"`
	Entries      []repo.Task `json:"entries" bson:"entries"`
}

type Service struct {
	// repo Repository
	repo repo.Queries
	db   *pgx.Conn
}

// NewService creates a new task service
func NewService(repo repo.Queries, db *pgx.Conn) *Service {
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
	case "id", "title", "created_at":
		sortBy = arg.SortBy
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

	tasks := make([]repo.Task, 0)

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
		tasks = append(tasks, task)
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
func (s *Service) GetTaskByID(ctx context.Context, id int64) (repo.Task, error) {
	actorID := middleware.UserIDFromContext(ctx)
	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.Task{}, apperrors.NotFound("task not found")
		}
		return repo.Task{}, err
	}

	allowed, err := s.canAccessTask(ctx, actorID, task.OwnerID, "read")
	if err != nil {
		return repo.Task{}, err
	}
	if !allowed {
		return repo.Task{}, apperrors.Forbidden("insufficient permission")
	}

	return task, nil
}

// CreateTask creates a new task
func (s *Service) CreateTask(ctx context.Context, task *repo.CreateTaskParams) (*repo.Task, error) {
	actorID := middleware.UserIDFromContext(ctx)
	allowed, err := s.canAccessTask(ctx, actorID, actorID, "write")
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, apperrors.Forbidden("insufficient permission")
	}

	task.OwnerID = actorID

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := s.repo.WithTx(tx)
	created, err := qtx.CreateTask(ctx, *task)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &created, nil
}

// UpdateTask updates an existing task
func (s *Service) UpdateTask(ctx context.Context, id int64, task *repo.Task) (*repo.Task, error) {
	actorID := middleware.UserIDFromContext(ctx)
	current, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound("task not found")
		}
		return nil, err
	}

	allowed, err := s.canAccessTask(ctx, actorID, current.OwnerID, "write")
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, apperrors.Forbidden("insufficient permission")
	}

	arg := repo.UpdateTaskParams{
		ID:          id,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Completed:   task.Completed,
	}

	updated, err := s.repo.UpdateTask(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound("task not found")
		}
		return nil, err
	}

	return &updated, nil
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
