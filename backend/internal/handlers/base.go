package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/konnikamii/svelte-go-task-app/backend/internal/apperrors"
	"github.com/sirupsen/logrus"
)

// Response wrapper for consistent API responses
type Response struct {
	Success   bool   `json:"success"`
	Data      any    `json:"data,omitempty"`
	Error     string `json:"error,omitempty"`
	ErrorCode string `json:"errorCode,omitempty"`
}

// BaseHandler provides common handler functionality
type BaseHandler struct{}

func (bh *BaseHandler) Read(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

// JSON sends a JSON response with the specified status code
func (bh *BaseHandler) JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logrus.WithError(err).Error("failed to encode response")
	}
}

// Error sends an error response
func (bh *BaseHandler) Error(w http.ResponseWriter, status int, message string) {
	bh.ErrorWithCode(w, status, message, "")
}

// ErrorWithCode sends an error response with a stable machine-readable code.
func (bh *BaseHandler) ErrorWithCode(w http.ResponseWriter, status int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success:   false,
		Error:     message,
		ErrorCode: code,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.WithError(err).Error("failed to encode error response")
	}
}

// AppError sends either a known app error or a generic internal error.
func (bh *BaseHandler) AppError(w http.ResponseWriter, err error) {
	appErr := apperrors.FromError(err)
	bh.ErrorWithCode(w, appErr.Status, appErr.Message, appErr.Code)
}
