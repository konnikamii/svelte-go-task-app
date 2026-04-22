package contact

import (
	"net/http"

	"github.com/konnikamii/svelte-go-task-app/backend/internal/apperrors"
	"github.com/konnikamii/svelte-go-task-app/backend/internal/handlers"
	"github.com/sirupsen/logrus"
)

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

func (h *Handler) CreateContactRequest(w http.ResponseWriter, r *http.Request) {
	var body CreateContactRequest
	if err := h.Read(r, &body); err != nil {
		logrus.WithError(err).Warn("failed to parse contact request")
		h.AppError(w, apperrors.BadRequest("invalid request body"))
		return
	}

	created, err := h.service.CreateContactRequest(r.Context(), &body)
	if err != nil {
		h.AppError(w, err)
		return
	}

	h.JSON(w, http.StatusCreated, created)
}
