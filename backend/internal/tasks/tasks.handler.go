package tasks

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

// NewHandler creates a new task handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		BaseHandler: &handlers.BaseHandler{},
		service:     service,
	}
}

// GetTasks handles GET /tasks
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logrus.Error("asdasd")
	tasks, err := h.service.GetTasks(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to get tasks")
		h.Error(w, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}

	h.JSON(w, http.StatusOK, tasks)
}

// CreateTask handles POST /tasks
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		logrus.WithError(err).Warn("failed to decode task")
		h.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	created, err := h.service.CreateTask(ctx, &task)
	if err != nil {
		appErr := err.(Error)
		h.Error(w, appErr.Code, appErr.Message)
		return
	}

	h.JSON(w, http.StatusCreated, created)
}

// GetTaskByID handles GET /tasks/{id}
func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithField("id", idStr).Warn("invalid task id")
		h.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}

	task, err := h.service.GetTaskByID(ctx, id)
	if err != nil {
		appErr := err.(Error)
		h.Error(w, appErr.Code, appErr.Message)
		return
	}

	h.JSON(w, http.StatusOK, task)
}

// UpdateTask handles PUT /tasks/{id}
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithField("id", idStr).Warn("invalid task id")
		h.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		logrus.WithError(err).Warn("failed to decode task")
		h.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updated, err := h.service.UpdateTask(ctx, id, &task)
	if err != nil {
		appErr := err.(Error)
		h.Error(w, appErr.Code, appErr.Message)
		return
	}

	h.JSON(w, http.StatusOK, updated)
}

// DeleteTask handles DELETE /tasks/{id}
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithField("id", idStr).Warn("invalid task id")
		h.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}

	err = h.service.DeleteTask(ctx, id)
	if err != nil {
		appErr := err.(Error)
		h.Error(w, appErr.Code, appErr.Message)
		return
	}

	h.JSON(w, http.StatusNoContent, nil)
}
