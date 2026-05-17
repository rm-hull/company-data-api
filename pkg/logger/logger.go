package logger

import (
	"log/slog"
	"os"
)

// SetupLogger configures the global logger to use a JSON handler.
func SetupLogger() {
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
