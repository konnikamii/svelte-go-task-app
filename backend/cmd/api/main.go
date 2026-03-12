package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/env"
	"github.com/sirupsen/logrus"
)

func main() {
	// Config
	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	// Logger
	// logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})

	ctx := context.Background()
	cfg := config{
		addr: env.GetString("SERVER_HOST", "") + ":" + env.GetString("SERVER_PORT", ""),
		db: dbConfig{
			dsn: env.GetString("GOOSE_DBSTRING", ""),
		},
	}

	// Database
	conn, err := pgx.Connect(ctx, cfg.db.dsn)
	if err != nil {
		logrus.Error("Error connecting to database!")
		panic(err)
	}
	defer conn.Close(ctx)

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
