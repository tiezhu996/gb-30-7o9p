package util

import (
	"log/slog"
	"os"
)

// NewLogger builds a structured JSON slog logger.
func NewLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

// LoggerWithRequest returns a logger bound with request metadata.
func LoggerWithRequest(logger *slog.Logger, requestID, method, path string) *slog.Logger {
	return logger.With("request_id", requestID, "method", method, "path", path)
}
