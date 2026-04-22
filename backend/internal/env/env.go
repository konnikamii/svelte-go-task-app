package env

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerHost             string
	ServerPort             string
	DBString               string
	FrontendURL            string
	SessionSecret          string
	SessionDurationMinutes int
	CookieSecure           bool
	SMTPHost               string
	SMTPPort               int
	SMTPFrom               string
	SMTPTo                 string
}

func GetString(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return fallback
}

func InitEnv() (*Config, error) {
	requiredVars := []string{
		"GOOSE_DBSTRING",
		"SESSION_SECRET",
	}

	for _, varName := range requiredVars {
		if os.Getenv(varName) == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", varName)
		}
	}

	sessionDurationMinutes := 1440
	if val := os.Getenv("SESSION_DURATION_MINUTES"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			sessionDurationMinutes = parsed
		}
	}

	smtpPort := 1025
	if val := os.Getenv("SMTP_PORT"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			smtpPort = parsed
		}
	}

	return &Config{
		ServerHost:             GetString("SERVER_HOST", "localhost"),
		ServerPort:             GetString("SERVER_PORT", "8000"),
		DBString:               os.Getenv("GOOSE_DBSTRING"),
		FrontendURL:            GetString("FRONTEND_URL", "http://localhost:5173"),
		SessionSecret:          os.Getenv("SESSION_SECRET"),
		SessionDurationMinutes: sessionDurationMinutes,
		CookieSecure:           GetString("COOKIE_SECURE", "false") == "true",
		SMTPHost:               GetString("SMTP_HOST", "mailhog"),
		SMTPPort:               smtpPort,
		SMTPFrom:               GetString("SMTP_FROM", "no-reply@taskify.local"),
		SMTPTo:                 GetString("SMTP_TO", "contact@taskify.local"),
	}, nil
}
