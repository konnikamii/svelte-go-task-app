package auth

import (
	"net/http"

	"github.com/konnikamii/svelte-go-task-app/backend/internal/apperrors"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/handlers"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/middleware"
	"github.com/sirupsen/logrus"
)

// UserInfo is the client-facing user shape returned by login and /me.
type UserInfo struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Handler struct {
	*handlers.BaseHandler
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		BaseHandler: &handlers.BaseHandler{},
		service:     service,
	}
}

// Login handles POST /api/login — open
// Accepts multipart form with fields: email, password
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2000); err != nil {
		h.AppError(w, apperrors.BadRequest("invalid form data"))
		return
	}

	params := LoginParams{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	user, err := h.service.Login(r.Context(), params)
	if err != nil {
		logrus.WithError(err).Warn("login failed")
		h.AppError(w, err)
		return
	}

	// Derive the device ID from request headers.
	deviceID := middleware.DeviceIDFromRequest(r)

	// Revoke any existing active sessions for this user on this device.
	// This ensures one active session per user per device (supports multi-device logins).
	if err := middleware.RevokeActiveSessionsForUserOnDevice(r.Context(), user.ID, deviceID); err != nil {
		logrus.WithError(err).Warn("failed to revoke previous device sessions")
	}

	if err := middleware.CleanupStaleSessions(r.Context()); err != nil {
		logrus.WithError(err).Warn("failed to cleanup stale sessions")
	}

	token, err := middleware.CreateSession(r.Context(), user.ID, deviceID)
	if err != nil {
		h.AppError(w, apperrors.Internal("could not create login session"))
		return
	}
	middleware.SetSessionCookie(w, token)

	h.JSON(w, http.StatusOK, UserInfo{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

// Logout handles POST /api/logout — open (clears session)
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if token, err := middleware.SessionTokenFromRequest(r); err == nil {
		if err := middleware.RevokeSession(r.Context(), token); err != nil {
			logrus.WithError(err).Warn("failed to revoke session")
		}
	}

	middleware.ClearSessionCookie(w)
	h.JSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

// Me handles GET /api/me — requires RequireAuth middleware
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.UserIDFromContext(ctx)

	user, err := h.service.GetMe(ctx, userID)
	if err != nil {
		logrus.WithError(err).Error("failed to get current user")
		h.AppError(w, err)
		return
	}

	h.JSON(w, http.StatusOK, UserInfo{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}
