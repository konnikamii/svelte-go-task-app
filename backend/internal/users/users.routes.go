package users

import (
	"github.com/go-chi/chi/v5"
)

// Routes sets up all user-related routes
func Routes(r *chi.Mux) {

	handler := NewHandler(NewService(NewMockRepository()))

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", handler.GetUsers)          // GET /users
		r.Post("/", handler.CreateUser)       // POST /users
		r.Get("/{id}", handler.GetUserByID)   // GET /users/{id}
		r.Put("/{id}", handler.UpdateUser)    // PUT /users/{id}
		r.Delete("/{id}", handler.DeleteUser) // DELETE /users/{id}
	})
}
