package users

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/middleware"
)

// Routes sets up all user-related routes.
// POST /users is open (registration). All other routes require a valid session.
func Routes(r chi.Router, db *pgxpool.Pool) {
	handler := NewHandler(NewService(*repo.New(db), db))

	r.Route("/users", func(r chi.Router) {
		// Open
		r.Post("/", handler.CreateUser) // POST /users — register

		// Protected
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth)
			r.Get("/{id}", handler.GetUserByID)        // GET /users/{id}
			r.Post("/list", handler.GetUsersPaginated) // POST /users/list
			r.Put("/{id}", handler.UpdateUser)         // PUT /users/{id}
			r.Delete("/{id}", handler.DeleteUser)      // DELETE /users/{id}
		})
	})
}
