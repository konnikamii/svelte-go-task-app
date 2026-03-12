package tasks

import (
	"github.com/go-chi/chi/v5"
)

func Routes(r *chi.Mux) {

	handler := NewHandler(NewService(NewMockRepository()))

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", handler.GetTasks)          // GET /tasks
		r.Post("/", handler.CreateTask)       // POST /tasks
		r.Get("/{id}", handler.GetTaskByID)   // GET /tasks/{id}
		r.Put("/{id}", handler.UpdateTask)    // PUT /tasks/{id}
		r.Delete("/{id}", handler.DeleteTask) // DELETE /tasks/{id}
	})

}
