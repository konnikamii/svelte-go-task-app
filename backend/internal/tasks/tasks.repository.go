package tasks

import (
	"context"
)

// MockRepository is a mock implementation of Repository for testing
type MockRepository struct {
	tasks map[int]*Task
}

// NewMockRepository creates a new mock repository
func NewMockRepository() *MockRepository {
	return &MockRepository{
		tasks: map[int]*Task{
			1: {ID: 1, Title: "Sample Task 1", Description: "This is a sample task", Completed: false},
			2: {ID: 2, Title: "Sample Task 2", Description: "Another sample task", Completed: true},
		},
	}
}

func (mr *MockRepository) GetTasks(ctx context.Context) ([]Task, error) {
	var tasks []Task
	for _, task := range mr.tasks {
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func (mr *MockRepository) GetTaskByID(ctx context.Context, id int) (*Task, error) {
	if task, exists := mr.tasks[id]; exists {
		return task, nil
	}
	return nil, ErrTaskNotFound
}

func (mr *MockRepository) CreateTask(ctx context.Context, task *Task) (*Task, error) {
	// Simple ID generation for mock
	task.ID = len(mr.tasks) + 1
	mr.tasks[task.ID] = task
	return task, nil
}

func (mr *MockRepository) UpdateTask(ctx context.Context, id int, task *Task) (*Task, error) {
	if _, exists := mr.tasks[id]; !exists {
		return nil, ErrTaskNotFound
	}
	task.ID = id
	mr.tasks[id] = task
	return task, nil
}

func (mr *MockRepository) DeleteTask(ctx context.Context, id int) error {
	if _, exists := mr.tasks[id]; !exists {
		return ErrTaskNotFound
	}
	delete(mr.tasks, id)
	return nil
}
