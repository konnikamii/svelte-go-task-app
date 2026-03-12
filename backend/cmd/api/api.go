package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/handlers"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/tasks"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/users"
	"github.com/sirupsen/logrus"
)

type application struct {
	config config
	// logger
	db *pgx.Conn
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID) // rate limiting
	r.Use(middleware.RealIP)    // rate limiting / analytics / tracing
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags), NoColor: false}))
	r.Use(middleware.Recoverer) // recover from crashes

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	// Common routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good"))
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(""))
	})

	// Tasks routes
	tasks.Routes(r)
	users.Routes(r)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)

		if err := json.NewEncoder(w).Encode(handlers.Response{Success: false}); err != nil {
			logrus.WithError(err).Error("failed to encode response")
		}
	})

	return r
}

func (app *application) run(h http.Handler) error {
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	logrus.Printf("Server started on: http://%s", app.config.addr)

	return server.ListenAndServe()
}
