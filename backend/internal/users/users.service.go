package users

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Service struct {
	repo Repository
}

// NewService creates a new user service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetUsers retrieves all users
func (s *Service) GetUsers(ctx context.Context) ([]User, error) {
	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch users from repository")
		return nil, ErrInternal
	}
	return users, nil
}

// GetUserByID retrieves a single user by ID
func (s *Service) GetUserByID(ctx context.Context, id int) (*User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("failed to fetch user from repository")
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetUserByEmail retrieves a user by email
func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		logrus.WithError(err).WithField("email", email).Error("failed to fetch user from repository")
		return nil, ErrUserNotFound
	}
	return user, nil
}

// CreateUser creates a new user
func (s *Service) CreateUser(ctx context.Context, user *User) (*User, error) {
	if err := validateUser(user); err != nil {
		logrus.WithError(err).Warn("invalid user data")
		return nil, ErrInvalidInput
	}

	// Check if user already exists
	existing, _ := s.repo.GetUserByEmail(ctx, user.Email)
	if existing != nil {
		logrus.WithField("email", user.Email).Warn("user already exists")
		return nil, ErrUserExists
	}

	created, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		logrus.WithError(err).Error("failed to create user in repository")
		return nil, ErrInternal
	}

	logrus.WithField("user_id", created.ID).Info("user created successfully")
	return created, nil
}

// UpdateUser updates an existing user
func (s *Service) UpdateUser(ctx context.Context, id int, user *User) (*User, error) {
	if err := validateUser(user); err != nil {
		logrus.WithError(err).Warn("invalid user data")
		return nil, ErrInvalidInput
	}

	updated, err := s.repo.UpdateUser(ctx, id, user)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("failed to update user in repository")
		return nil, ErrInternal
	}

	logrus.WithField("user_id", id).Info("user updated successfully")
	return updated, nil
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(ctx context.Context, id int) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		logrus.WithError(err).WithField("id", id).Error("failed to delete user from repository")
		return ErrInternal
	}

	logrus.WithField("user_id", id).Info("user deleted successfully")
	return nil
}

// validateUser validates user data
func validateUser(user *User) error {
	if user == nil {
		return ErrInvalidInput
	}
	if user.Name == "" || user.Email == "" {
		return ErrInvalidInput
	}
	return nil
}
