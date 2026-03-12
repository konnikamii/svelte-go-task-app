package users

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

// GetUsers retrieves all users from the database
func (pr *PostgresRepository) GetUsers(ctx context.Context) ([]User, error) {
	if pr.db == nil {
		logrus.Warn("database connection is nil, returning empty users")
		return []User{}, nil
	}

	rows, err := pr.db.Query(ctx, `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`)
	if err != nil {
		logrus.WithError(err).Error("failed to query users")
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			logrus.WithError(err).Error("failed to scan user row")
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		logrus.WithError(err).Error("error iterating user rows")
		return nil, err
	}

	return users, nil
}

// GetUserByID retrieves a single user by ID
func (pr *PostgresRepository) GetUserByID(ctx context.Context, id int) (*User, error) {
	if pr.db == nil {
		return nil, ErrUserNotFound
	}

	var user User
	err := pr.db.QueryRow(ctx, `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		logrus.WithError(err).Error("failed to query user by id")
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (pr *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := pr.db.QueryRow(ctx, `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		logrus.WithError(err).Error("failed to query user by email")
		return nil, err
	}
	if pr.db == nil {
		return nil, ErrInternal
	}
	return &user, nil
}

// CreateUser creates a new user in the database
func (pr *PostgresRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := pr.db.QueryRow(ctx, `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		logrus.WithError(err).Error("failed to create user")
		return nil, err
	}
	if pr.db == nil {
		return nil, ErrInternal
	}
	return user, nil
}

// UpdateUser updates an existing user in the database
func (pr *PostgresRepository) UpdateUser(ctx context.Context, id int, user *User) (*User, error) {
	user.ID = id
	user.UpdatedAt = time.Now()

	commandTag, err := pr.db.Exec(ctx, `
		UPDATE users
		SET name = $1, email = $2, password = $3, updated_at = $4
		WHERE id = $5
	`,
		user.Name,
		user.Email,
		user.Password,
		user.UpdatedAt,
		id,
	)

	if err != nil {
		logrus.WithError(err).Error("failed to update user")
		return nil, err
	}

	if commandTag.RowsAffected() == 0 {

		return nil, ErrUserNotFound
	}
	if pr.db == nil {
		return nil, ErrInternal
	}
	return user, nil
}

// DeleteUser deletes a user from the database
func (pr *PostgresRepository) DeleteUser(ctx context.Context, id int) error {
	commandTag, err := pr.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)

	if err != nil {
		logrus.WithError(err).Error("failed to delete user")
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
