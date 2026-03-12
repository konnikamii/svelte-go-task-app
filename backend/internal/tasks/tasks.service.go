package tasks

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Service struct {
	repo Repository
}

// NewService creates a new task service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetTasks retrieves all tasks
func (s *Service) GetTasks(ctx context.Context) ([]Task, error) {
	tasks, err := s.repo.GetTasks(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch tasks from repository")
		return nil, ErrInternal
	}
	return tasks, nil
}

// GetTaskByID retrieves a single task by ID
func (s *Service) GetTaskByID(ctx context.Context, id int) (*Task, error) {
	task, err := s.repo.GetTaskByID(ctx, id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("failed to fetch task from repository")
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// CreateTask creates a new task
func (s *Service) CreateTask(ctx context.Context, task *Task) (*Task, error) {
	if err := validateTask(task); err != nil {
		logrus.WithError(err).Warn("invalid task data")
		return nil, ErrInvalidInput
	}

	created, err := s.repo.CreateTask(ctx, task)
	if err != nil {
		logrus.WithError(err).Error("failed to create task in repository")
		return nil, ErrInternal
	}

	logrus.WithField("task_id", created.ID).Info("task created successfully")
	return created, nil
}

// UpdateTask updates an existing task
func (s *Service) UpdateTask(ctx context.Context, id int, task *Task) (*Task, error) {
	if err := validateTask(task); err != nil {
		logrus.WithError(err).Warn("invalid task data")
		return nil, ErrInvalidInput
	}

	updated, err := s.repo.UpdateTask(ctx, id, task)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("failed to update task in repository")
		return nil, ErrInternal
	}

	logrus.WithField("task_id", id).Info("task updated successfully")
	return updated, nil
}

// DeleteTask deletes a task
func (s *Service) DeleteTask(ctx context.Context, id int) error {
	err := s.repo.DeleteTask(ctx, id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("failed to delete task from repository")
		return ErrInternal
	}

	logrus.WithField("task_id", id).Info("task deleted successfully")
	return nil
}

// validateTask validates task data
func validateTask(task *Task) error {
	if task == nil {
		return ErrInvalidInput
	}
	if task.Title == "" {
		return ErrInvalidInput
	}
	return nil
}
