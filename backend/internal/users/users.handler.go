package users

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/handlers"
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

// GetUsers handles GET /users
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to get users")
		h.Error(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}

	h.JSON(w, http.StatusOK, users)
}

// CreateUser handles POST /users
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logrus.WithError(err).Warn("failed to decode user")
		h.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.service.CreateUser(ctx, &user)
	if err != nil {
		appErr := err.(Error)
		h.Error(w, appErr.Code, appErr.Message)
		return
	}

	h.JSON(w, http.StatusCreated, created)
}

// GetUserByID handles GET /users/{id}
func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithField("id", idStr).Warn("invalid user id")
		h.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.service.GetUserByID(ctx, id)
	if err != nil {
		appErr := err.(Error)
		h.Error(w, appErr.Code, appErr.Message)
		return
	}

	h.JSON(w, http.StatusOK, user)
}

// UpdateUser handles PUT /users/{id}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithField("id", idStr).Warn("invalid user id")
		h.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logrus.WithError(err).Warn("failed to decode user")
		h.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updated, err := h.service.UpdateUser(ctx, id, &user)
	if err != nil {
		appErr := err.(Error)
		h.Error(w, appErr.Code, appErr.Message)
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
		h.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	err = h.service.DeleteUser(ctx, id)
	if err != nil {
		appErr := err.(Error)
		h.Error(w, appErr.Code, appErr.Message)
		return
	}

	h.JSON(w, http.StatusNoContent, nil)
}
