package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New crea un logger estructurado.
// - En producción usa formato JSON.
// - En desarrollo usa formato de texto legible.
// - El nivel se controla con LOG_LEVEL (debug, info, warn, error).
func New(appEnv, logLevel string) *slog.Logger {
	var handler slog.Handler

	level := parseLevel(logLevel)

	if appEnv == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}

	return slog.New(handler)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}