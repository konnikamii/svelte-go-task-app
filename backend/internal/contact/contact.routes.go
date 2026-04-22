package contact

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
)

func Routes(r chi.Router, db *pgxpool.Pool, mailConfig MailConfig) {
	handler := NewHandler(NewService(*repo.New(db), db, mailConfig))

	r.Route("/contact", func(r chi.Router) {
		r.Post("/", handler.CreateContactRequest)
	})
}
