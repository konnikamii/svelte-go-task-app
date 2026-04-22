package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/middleware"
)

// Routes registers auth endpoints directly on the router (no /auth prefix).
// POST /login  — open; accepts multipart form fields: email, password
// POST /logout — open
// GET  /me     — requires valid session (RequireAuth)
func Routes(r chi.Router, db *pgxpool.Pool) {
	handler := NewHandler(NewService(*repo.New(db)))

	r.Post("/login", handler.Login)
	r.Post("/logout", handler.Logout)

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.Get("/me", handler.Me)
	})
}
