package tasks

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

// PostgresRepository is a PostgreSQL implementation of Repository
type PostgresRepository struct {
	db *pgx.Conn
}

// NewPostgresRepository creates a new postgres repository
func NewPostgresRepository(db *pgx.Conn) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// GetTasks retrieves all tasks from the database
func (pr *PostgresRepository) GetTasks(ctx context.Context) ([]Task, error) { // Return empty slice if db is nil
	logrus.Warn("asd")
	if pr.db == nil {
		logrus.Warn("database connection is nil, returning empty tasks")
		return []Task{}, nil
	}
	rows, err := pr.db.Query(ctx, `
		SELECT id, title, description, completed, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`)
	if err != nil {
		logrus.WithError(err).Error("failed to query tasks")
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			logrus.WithError(err).Error("failed to scan task row")
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		logrus.WithError(err).Error("error iterating task rows")
		return nil, err
	}
	if pr.db == nil {
		return nil, ErrInternal
	}

	return tasks, nil
}

// GetTaskByID retrieves a single task by ID
func (pr *PostgresRepository) GetTaskByID(ctx context.Context, id int) (*Task, error) {
	var task Task
	err := pr.db.QueryRow(ctx, `
		SELECT id, title, description, completed, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		logrus.WithError(err).Error("failed to query task by id")
		return nil, err
	}

	if pr.db == nil {
		return nil, ErrInternal
	}
	return &task, nil
}

// CreateTask creates a new task in the database
func (pr *PostgresRepository) CreateTask(ctx context.Context, task *Task) (*Task, error) {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	err := pr.db.QueryRow(ctx, `
		INSERT INTO tasks (title, description, completed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(&task.ID)

	if err != nil {
		logrus.WithError(err).Error("failed to create task")
		return nil, err
	}
	if pr.db == nil {
		return nil, ErrInternal
	}
	return task, nil
}

// UpdateTask updates an existing task in the database
func (pr *PostgresRepository) UpdateTask(ctx context.Context, id int, task *Task) (*Task, error) {
	task.ID = id
	task.UpdatedAt = time.Now()

	commandTag, err := pr.db.Exec(ctx, `
		UPDATE tasks
		SET title = $1, description = $2, completed = $3, updated_at = $4
		WHERE id = $5
	`,
		task.Title,
		task.Description,
		task.Completed,
		task.UpdatedAt,
		id,
	)

	if err != nil {
		logrus.WithError(err).Error("failed to update task")
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {
		return nil, ErrTaskNotFound
	}
	if pr.db == nil {
		return nil, ErrInternal
	}

	return task, nil
}

// DeleteTask deletes a task from the database
func (pr *PostgresRepository) DeleteTask(ctx context.Context, id int) error {
	commandTag, err := pr.db.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)

	if err != nil {
		logrus.WithError(err).Error("failed to delete task")
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}
