package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Response wrapper for consistent API responses
type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// BaseHandler provides common handler functionality
type BaseHandler struct{}

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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: false,
		Error:   message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.WithError(err).Error("failed to encode error response")
	}
}
