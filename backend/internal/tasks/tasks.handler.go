package tasks

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/apperrors"
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

// GetTasksPaginated handles GET /tasks
func (h *Handler) GetTasksPaginated(w http.ResponseWriter, r *http.Request) {
	var body PaginatedParams
	if err := h.Read(r, &body); err != nil {
		logrus.WithError(err).Error("failed to parse params")
		h.AppError(w, apperrors.BadRequest("invalid request body"))
		return
	}

	tasks, err := h.service.GetTasksPaginated(r.Context(), body)
	if err != nil {
		logrus.WithError(err).Error("failed to get tasks")
		h.AppError(w, err)
		return
	}

	h.JSON(w, http.StatusOK, tasks)
}

// CreateTask handles POST /tasks
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var body CreateTaskRequest
	if err := h.Read(r, &body); err != nil {
		logrus.WithError(err).Error("failed to parse params")
		h.AppError(w, apperrors.BadRequest("invalid request body"))
		return
	}

	created, err := h.service.CreateTask(r.Context(), &body)
	if err != nil {
		h.AppError(w, err)
		return
	}

	h.JSON(w, http.StatusCreated, created)
}

// GetTaskByID handles GET /tasks/{id}
func (h *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithField("id", idStr).Warn("invalid task id")
		h.AppError(w, apperrors.BadRequest("invalid task id"))
		return
	}

	task, err := h.service.GetTaskByID(r.Context(), int64(id))
	if err != nil {
		logrus.WithError(err).Error("failed to get task by id")
		h.AppError(w, err)
		return
	}

	h.JSON(w, http.StatusOK, task)
}

// UpdateTask handles PUT /tasks/{id}
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.WithField("id", idStr).Warn("invalid task id")
		h.AppError(w, apperrors.BadRequest("invalid task id"))
		return
	}

	var task UpdateTaskRequest
	if err := h.Read(r, &task); err != nil {
		logrus.WithError(err).Warn("failed to parse task")
		h.AppError(w, apperrors.BadRequest("invalid request body"))
		return
	}

	updated, err := h.service.UpdateTask(r.Context(), int64(id), &task)
	if err != nil {
		h.AppError(w, err)
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
		h.AppError(w, apperrors.BadRequest("invalid task id"))
		return
	}

	deleted, err := h.service.DeleteTask(ctx, int64(id))
	if err != nil {
		h.AppError(w, err)
		return
	}

	h.JSON(w, http.StatusNoContent, deleted)
}
