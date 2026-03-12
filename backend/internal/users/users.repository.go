package users

import (
	"context"
	"time"
)

// MockRepository is a mock implementation of Repository for testing
type MockRepository struct {
	users map[int]*User
}

// NewMockRepository creates a new mock repository
func NewMockRepository() *MockRepository {
	return &MockRepository{
		users: map[int]*User{
			1: {
				ID:        1,
				Name:      "John Doe",
				Email:     "john@example.com",
				Password:  "hashed_password_1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			2: {
				ID:        2,
				Name:      "Jane Smith",
				Email:     "jane@example.com",
				Password:  "hashed_password_2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
}

func (mr *MockRepository) GetUsers(ctx context.Context) ([]User, error) {
	var users []User
	for _, user := range mr.users {
		users = append(users, *user)
	}
	return users, nil
}

func (mr *MockRepository) GetUserByID(ctx context.Context, id int) (*User, error) {
	if user, exists := mr.users[id]; exists {
		return user, nil
	}
	return nil, ErrUserNotFound
}

func (mr *MockRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	for _, user := range mr.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (mr *MockRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	user.ID = len(mr.users) + 1
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	mr.users[user.ID] = user
	return user, nil
}

func (mr *MockRepository) UpdateUser(ctx context.Context, id int, user *User) (*User, error) {
	if _, exists := mr.users[id]; !exists {
		return nil, ErrUserNotFound
	}
	user.ID = id
	user.UpdatedAt = time.Now()
	mr.users[id] = user
	return user, nil
}

func (mr *MockRepository) DeleteUser(ctx context.Context, id int) error {
	if _, exists := mr.users[id]; !exists {
		return ErrUserNotFound
	}
	delete(mr.users, id)
	return nil
}
