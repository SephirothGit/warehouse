package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/SephirothGit/warehouse/internal/repository"
	"github.com/SephirothGit/warehouse/internal/service"
)

func respondError(w http.ResponseWriter, err error, fallbackMsg string) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, repository.ErrNotFound):
		http.Error(w, "resource not found", http.StatusNotFound)
	case errors.Is(err, repository.ErrAlreadyExists):
		http.Error(w, "resource already exists", http.StatusConflict)
	case errors.Is(err, repository.ErrInsufficientStock):
		http.Error(w, "insufficient stock", http.StatusConflict)
	case errors.Is(err, context.DeadlineExceeded):
		http.Error(w, "request timed out", http.StatusGatewayTimeout)
	default:
		slog.Error("request failed", "error", err)
		http.Error(w, fallbackMsg, http.StatusInternalServerError)
	}
}
