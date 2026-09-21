package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		start := time.Now()

		slog.Info("request started",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path)

		next.ServeHTTP(w, r)

		slog.Info("request completed",
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(start).Milliseconds())
	})
}
