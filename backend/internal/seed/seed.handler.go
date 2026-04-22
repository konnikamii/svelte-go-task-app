package seed

import (
	"net/http"

	"github.com/konnikamii/svelte-go-task-app/backend/internal/handlers"
)

type Handler struct {
	handlers.BaseHandler
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Seed(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.SeedDatabase(r.Context())
	if err != nil {
		h.AppError(w, err)
		return
	}

	h.JSON(w, http.StatusCreated, handlers.Response{Success: true, Data: result})
}
