package seed

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
)

func Routes(r chi.Router, db *pgxpool.Pool) {
	handler := NewHandler(NewService(*repo.New(db), db))
	r.Post("/seed", handler.Seed)
}
