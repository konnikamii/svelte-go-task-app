package seed

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/apperrors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo repo.Queries
	db   *pgxpool.Pool
}

type demoUserSeed struct {
	Username string
	Email    string
	Roles    []string
	Tasks    []demoTaskSeed
}

type demoTaskSeed struct {
	Title       string
	Description string
	DueDate     *time.Time
	Completed   bool
}

func NewService(repo repo.Queries, db *pgxpool.Pool) *Service {
	return &Service{repo: repo, db: db}
}

func (s *Service) SeedDatabase(ctx context.Context) (SeedResult, error) {
	existingUsers, err := s.countUsers(ctx)
	if err != nil {
		return SeedResult{}, err
	}
	if existingUsers > 0 {
		return SeedResult{}, apperrors.Conflict("database already contains users")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return SeedResult{}, err
	}
	defer tx.Rollback(ctx)

	if err := s.seedAuthorization(ctx, tx); err != nil {
		return SeedResult{}, err
	}

	sharedPassword := "Taskify123"
	hash, err := bcrypt.GenerateFromPassword([]byte(sharedPassword), 12)
	if err != nil {
		return SeedResult{}, err
	}

	qtx := s.repo.WithTx(tx)
	seedUsers := demoUsers(time.Now().UTC())
	credentials := make([]SeedCredential, 0, len(seedUsers))

	for _, seedUser := range seedUsers {
		createdUser, err := qtx.CreateUser(ctx, repo.CreateUserParams{
			Username: seedUser.Username,
			Email:    seedUser.Email,
			Password: string(hash),
		})
		if err != nil {
			return SeedResult{}, err
		}

		for _, roleName := range seedUser.Roles {
			assigned, err := qtx.AssignRoleByNameToUser(ctx, repo.AssignRoleByNameToUserParams{
				UserID: createdUser.ID,
				Name:   roleName,
			})
			if err != nil {
				return SeedResult{}, err
			}
			if assigned == 0 {
				return SeedResult{}, apperrors.Internal("seed role assignment failed")
			}
		}

		for _, taskSeed := range seedUser.Tasks {
			_, err := qtx.CreateTask(ctx, repo.CreateTaskParams{
				OwnerID:     createdUser.ID,
				Title:       taskSeed.Title,
				Description: nullableText(taskSeed.Description),
				DueDate:     nullableTime(taskSeed.DueDate),
				Completed:   taskSeed.Completed,
			})
			if err != nil {
				return SeedResult{}, err
			}
		}

		credentials = append(credentials, SeedCredential{
			Username: seedUser.Username,
			Email:    seedUser.Email,
			Password: sharedPassword,
			Roles:    append([]string(nil), seedUser.Roles...),
		})
	}

	totals, err := seededTotals(ctx, tx)
	if err != nil {
		return SeedResult{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return SeedResult{}, err
	}

	return SeedResult{Totals: totals, Credentials: credentials}, nil
}

func (s *Service) countUsers(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

func (s *Service) seedAuthorization(ctx context.Context, tx pgx.Tx) error {
	if _, err := tx.Exec(ctx, `
INSERT INTO roles (name, description)
VALUES ('admin', 'Full access to all resources'),
  ('user', 'Can read and write only own resources')
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description,
  updated_at = CURRENT_TIMESTAMP;
`); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
INSERT INTO permissions (resource, action, scope, description)
VALUES ('user', 'read', 'any', 'Read any user'),
  ('user', 'write', 'any', 'Write any user'),
  ('task', 'read', 'any', 'Read any task'),
  ('task', 'write', 'any', 'Write any task'),
  ('user', 'read', 'own', 'Read own user'),
  ('user', 'write', 'own', 'Write own user'),
  ('task', 'read', 'own', 'Read own task'),
  ('task', 'write', 'own', 'Write own task')
ON CONFLICT (resource, action, scope) DO UPDATE
SET description = EXCLUDED.description,
  updated_at = CURRENT_TIMESTAMP;
`); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, `
INSERT INTO roles_permissions (role_id, permission_id)
SELECT r.id,
  p.id
FROM roles r
  JOIN permissions p ON (
    (r.name = 'admin' AND p.scope = 'any')
    OR (r.name = 'user' AND p.scope = 'own')
  )
ON CONFLICT (role_id, permission_id) DO NOTHING;
`); err != nil {
		return err
	}

	return nil
}

func seededTotals(ctx context.Context, tx pgx.Tx) (SeedTotals, error) {
	users, err := countTable(ctx, tx, `SELECT COUNT(*) FROM users`)
	if err != nil {
		return SeedTotals{}, err
	}
	tasks, err := countTable(ctx, tx, `SELECT COUNT(*) FROM tasks`)
	if err != nil {
		return SeedTotals{}, err
	}
	roles, err := countTable(ctx, tx, `SELECT COUNT(*) FROM roles`)
	if err != nil {
		return SeedTotals{}, err
	}
	permissions, err := countTable(ctx, tx, `SELECT COUNT(*) FROM permissions`)
	if err != nil {
		return SeedTotals{}, err
	}
	rolePermissions, err := countTable(ctx, tx, `SELECT COUNT(*) FROM roles_permissions`)
	if err != nil {
		return SeedTotals{}, err
	}
	userRoles, err := countTable(ctx, tx, `SELECT COUNT(*) FROM users_roles`)
	if err != nil {
		return SeedTotals{}, err
	}

	return SeedTotals{
		Users:           users,
		Tasks:           tasks,
		Roles:           roles,
		Permissions:     permissions,
		RolePermissions: rolePermissions,
		UserRoles:       userRoles,
	}, nil
}

func countTable(ctx context.Context, tx pgx.Tx, query string) (int64, error) {
	var count int64
	err := tx.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func nullableText(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}

	return pgtype.Text{String: value, Valid: true}
}

func nullableTime(value *time.Time) pgtype.Timestamptz {
	if value == nil {
		return pgtype.Timestamptz{}
	}

	return pgtype.Timestamptz{Time: *value, Valid: true}
}

func demoUsers(now time.Time) []demoUserSeed {
	soon := now.Add(36 * time.Hour)
	later := now.Add(72 * time.Hour)
	overdue := now.Add(-18 * time.Hour)
	nextWeek := now.Add(7 * 24 * time.Hour)

	return []demoUserSeed{
		{
			Username: "demo-admin",
			Email:    "admin@taskify.local",
			Roles:    []string{"admin"},
			Tasks: []demoTaskSeed{
				{Title: "Review platform health", Description: "Check authentication, task throughput, and API readiness.", DueDate: &soon, Completed: false},
				{Title: "Approve onboarding copy", Description: "Confirm the public onboarding flow reads clearly.", DueDate: &later, Completed: false},
				{Title: "Archive stale incidents", Description: "Close out completed issues from the previous sprint.", DueDate: nil, Completed: true},
			},
		},
		{
			Username: "demo-olivia",
			Email:    "olivia@taskify.local",
			Roles:    []string{"user"},
			Tasks: []demoTaskSeed{
				{Title: "Prepare launch checklist", Description: "Finalize the release checklist for the upcoming client rollout.", DueDate: &overdue, Completed: false},
				{Title: "Write task detail notes", Description: "Expand task descriptions so handoffs are easier.", DueDate: &nextWeek, Completed: false},
				{Title: "Clean up backlog labels", Description: "Normalize labels across the current task backlog.", DueDate: nil, Completed: true},
			},
		},
		{
			Username: "demo-mateo",
			Email:    "mateo@taskify.local",
			Roles:    []string{"user"},
			Tasks: []demoTaskSeed{
				{Title: "Follow up with design", Description: "Resolve the last open comments on the dashboard polish pass.", DueDate: &soon, Completed: false},
				{Title: "Record handoff video", Description: "Capture a short walkthrough for the seeded workspace demo.", DueDate: &later, Completed: false},
				{Title: "Verify sample accounts", Description: "Make sure the seeded demo users can log in successfully.", DueDate: nil, Completed: true},
			},
		},
	}
}
