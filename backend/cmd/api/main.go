package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/contact"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/env"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	// Config
	if err := godotenv.Load(); err != nil {
		logrus.WithError(err).Info("No .env file loaded, using environment variables")
	}

	// Logger
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	// Init env
	envCfg, err := env.InitEnv()
	if err != nil {
		logrus.Error("Error initializing env variables!")
		panic(err)
	}

	ctx := context.Background()
	cfg := config{
		addr: envCfg.ServerHost + ":" + envCfg.ServerPort,
		db: dbConfig{
			dsn: envCfg.DBString,
		},
		frontendURL: envCfg.FrontendURL,
		mail: contact.MailConfig{
			Host: envCfg.SMTPHost,
			Port: envCfg.SMTPPort,
			From: envCfg.SMTPFrom,
			To:   envCfg.SMTPTo,
		},
	}

	// Database
	pool, err := pgxpool.New(ctx, cfg.db.dsn)
	if err != nil {
		logrus.Error("Error connecting to database!")
		panic(err)
	}
	defer pool.Close()

	queries := repo.New(pool)

	// Session store (cookie token + server-side DB session records)
	middleware.InitStore(
		envCfg.SessionSecret,
		envCfg.CookieSecure,
		envCfg.SessionDurationMinutes,
		queries,
	)

	// Start
	app := application{
		config: cfg,
		db:     pool,
	}
	if err := app.run(app.mount()); err != nil {
		logrus.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
