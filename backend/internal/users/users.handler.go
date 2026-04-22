package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	repo "github.com/konnikamii/svelte-go-task-app/backend/internal/adapters/postgresql/sqlc/out"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/apperrors"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/handlers"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/middleware"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	*handlers.BaseHandler
	service *Service
}

// NewHandler creates a new user handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		BaseHandler: &handlers.BaseHandler{},
		service:     service,
	}
}

// GetUsersPaginated handles POST /users/list
func (h *Handler) GetUsersPaginated(w http.ResponseWriter, r *http.Request) {
	var body PaginatedRequest
	if err := h.Read(r, &body); err != nil {
		logrus.WithError(err).Error("failed to parse params")
		h.AppError(w, apperrors.BadRequest("invalid request body"))
		return
	}

	users, err := h.service.GetUsersPaginated(r.Context(), body)
	if err != nil {
		logrus.WithError(err).Error("failed to get users")
		h.AppError(w, err)
		return
	}

	h.JSON(w, http.StatusOK, users)
}

// CreateUser handles POST /users
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2000); err != nil {
		logrus.WithError(err).Error("failed to parse form")
		h.AppError(w, apperrors.BadRequest("invalid form data"))
		return
	}

	body := repo.CreateUserParams{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	created, err := h.service.CreateUser(r.Context(), &body)
	if err != nil {
		logrus.WithError(err).Error("failed to create user")
		h.AppError(w, err)
		return
	}

	if err := middleware.StartUserSession(r.Context(), w, r, created.ID); err != nil {
		logrus.WithError(err).Error("failed to create registration session")
		h.AppError(w, apperrors.Internal("could not create registration session"))
		return
	}

	h.JSON(w, http.StatusCreated, created)
}

// GetUserByID handles GET /users/{id}
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithField("id", idStr).Warn("invalid user id")
		h.AppError(w, apperrors.BadRequest("invalid user id"))
		return
	}

	user, err := h.service.GetUserByID(r.Context(), int32(id))
	if err != nil {
		logrus.WithError(err).Error("failed to get user")
		h.AppError(w, err)
		return
	}

	h.JSON(w, http.StatusOK, user)
}

// UpdateUser handles PUT /users/{id}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithField("id", idStr).Warn("invalid user id")
		h.AppError(w, apperrors.BadRequest("invalid user id"))
		return
	}

	var user UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logrus.WithError(err).Warn("failed to decode user")
		h.AppError(w, apperrors.BadRequest("invalid request body"))
		return
	}

	updated, err := h.service.UpdateUser(r.Context(), int32(id), &user)
	if err != nil {
		logrus.WithError(err).Error("failed to update user")
		h.AppError(w, err)
		return
	}

	h.JSON(w, http.StatusOK, updated)
}

// DeleteUser handles DELETE /users/{id}
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithField("id", idStr).Warn("invalid user id")
		h.AppError(w, apperrors.BadRequest("invalid user id"))
		return
	}

	deleted, err := h.service.DeleteUser(ctx, int32(id))
	if err != nil {
		logrus.WithError(err).Error("failed to delete user")
		h.AppError(w, err)
		return
	}

	h.JSON(w, http.StatusNoContent, deleted)
}
