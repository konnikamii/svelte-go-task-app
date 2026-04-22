package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
)

const (
	SessionCookieName = "sid"
)

var (
	sessionSecretKey []byte
	sessionOptions   *http.Cookie
	sessionTTL       time.Duration
	sessionRepo      *repo.Queries
)

// InitStore initializes session settings and database access.
// Must be called once at startup before any requests are handled.
func InitStore(secret string, secure bool, durationMinutes int, queries *repo.Queries) {
	sessionSecretKey = []byte(secret)
	sessionTTL = time.Duration(durationMinutes) * time.Minute
	sessionRepo = queries
	sessionOptions = &http.Cookie{
		Name:     SessionCookieName,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
		MaxAge:   durationMinutes * 60,
	}
}

func ensureSessionStoreInitialized() error {
	if len(sessionSecretKey) == 0 {
		return errors.New("session secret is not initialized")
	}
	if sessionRepo == nil {
		return errors.New("session repository is not initialized")
	}
	if sessionOptions == nil {
		return errors.New("session options are not initialized")
	}
	return nil
}

// DeviceIDFromRequest derives a stable device identifier from request headers.
// This allows distinguishing between different browsers/devices/test contexts.
// Uses User-Agent, Accept-Language, and Accept-Encoding as a fingerprint.
func DeviceIDFromRequest(r *http.Request) string {
	// Combine stable request headers for device fingerprinting.
	input := r.Header.Get("User-Agent") + "|" +
		r.Header.Get("Accept-Language") + "|" +
		r.Header.Get("Accept-Encoding")

	// Hash with session secret for additional security.
	mac := hmac.New(sha256.New, sessionSecretKey)
	_, _ = mac.Write([]byte(input))
	return hex.EncodeToString(mac.Sum(nil))[:24] // Truncate to 24 chars for VARCHAR(128)
}

func CreateSession(ctx context.Context, userID int32, deviceID string) (string, error) {
	if err := ensureSessionStoreInitialized(); err != nil {
		return "", err
	}

	token, err := newSessionToken()
	if err != nil {
		return "", err
	}

	_, err = sessionRepo.CreateSession(ctx, repo.CreateSessionParams{
		UserID:    userID,
		DeviceID:  deviceID,
		TokenHash: hashSessionToken(token),
		ExpiresAt: pgtype.Timestamptz{Time: time.Now().Add(sessionTTL), Valid: true},
	})
	if err != nil {
		return "", err
	}

	return token, nil
}

func SetSessionCookie(w http.ResponseWriter, token string) {
	cookie := *sessionOptions
	cookie.Value = token
	http.SetCookie(w, &cookie)
}

func ClearSessionCookie(w http.ResponseWriter) {
	cookie := *sessionOptions
	cookie.Value = ""
	cookie.MaxAge = -1
	cookie.Expires = time.Unix(0, 0)
	http.SetCookie(w, &cookie)
}

func SessionTokenFromRequest(r *http.Request) (string, error) {
	if err := ensureSessionStoreInitialized(); err != nil {
		return "", err
	}

	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", errors.New("missing session token")
	}

	return cookie.Value, nil
}

func SessionUserID(ctx context.Context, token string) (int32, error) {
	if err := ensureSessionStoreInitialized(); err != nil {
		return 0, err
	}

	session, err := sessionRepo.GetActiveSessionByTokenHash(ctx, hashSessionToken(token))
	if err != nil {
		return 0, err
	}

	return session.UserID, nil
}

func RevokeSession(ctx context.Context, token string) error {
	if err := ensureSessionStoreInitialized(); err != nil {
		return err
	}

	_, err := sessionRepo.RevokeSessionByTokenHash(ctx, hashSessionToken(token))
	return err
}

// RevokeActiveSessionsForUserOnDevice revokes all active sessions for this user on this device.
// Used at login to ensure only one session per user per device.
func RevokeActiveSessionsForUserOnDevice(ctx context.Context, userID int32, deviceID string) error {
	if err := ensureSessionStoreInitialized(); err != nil {
		return err
	}

	_, err := sessionRepo.RevokeSessionsByUserAndDevice(ctx, repo.RevokeSessionsByUserAndDeviceParams{
		UserID:   userID,
		DeviceID: deviceID,
	})
	return err
}

func CleanupStaleSessions(ctx context.Context) error {
	if err := ensureSessionStoreInitialized(); err != nil {
		return err
	}

	_, err := sessionRepo.DeleteStaleSessions(ctx)
	return err
}

func StartUserSession(ctx context.Context, w http.ResponseWriter, r *http.Request, userID int32) error {
	deviceID := DeviceIDFromRequest(r)

	if err := RevokeActiveSessionsForUserOnDevice(ctx, userID, deviceID); err != nil {
		return err
	}

	if err := CleanupStaleSessions(ctx); err != nil {
		return err
	}

	token, err := CreateSession(ctx, userID, deviceID)
	if err != nil {
		return err
	}

	SetSessionCookie(w, token)
	return nil
}

func newSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate session token: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashSessionToken(token string) string {
	mac := hmac.New(sha256.New, sessionSecretKey)
	_, _ = mac.Write([]byte(token))
	return hex.EncodeToString(mac.Sum(nil))
}
