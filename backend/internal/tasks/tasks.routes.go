package tasks

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/middleware"
)

// Routes sets up all task-related routes. All routes require a valid session.
func Routes(r chi.Router, db *pgxpool.Pool) {
	handler := NewHandler(NewService(*repo.New(db), db))

	r.Route("/tasks", func(r chi.Router) {
		r.Use(middleware.RequireAuth)
		r.Get("/{id}", handler.GetTaskByID)        // GET /tasks/{id}
		r.Post("/", handler.CreateTask)            // POST /tasks
		r.Post("/list", handler.GetTasksPaginated) // POST /tasks/list
		r.Put("/{id}", handler.UpdateTask)         // PUT /tasks/{id}
		r.Delete("/{id}", handler.DeleteTask)      // DELETE /tasks/{id}
	})
}
