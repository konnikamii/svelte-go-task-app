package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/env"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	// Config
	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file")
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
	}

	// Database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		logrus.Error("Error connecting to database!")
		panic(err)
	}
	defer conn.Close(ctx)

	queries := repo.New(conn)

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
		db:     conn,
	}
	if err := app.run(app.mount()); err != nil {
		logrus.Error("server failed to start", "error", err)
		os.Exit(1)
	}
}
