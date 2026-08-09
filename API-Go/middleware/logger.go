package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger registra cada petición HTTP de forma estructurada con slog.
// Sustituye a gin.Logger() para mantener el formato del logger de la app.
func Logger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		if raw := c.Request.URL.RawQuery; raw != "" {
			path = path + "?" + raw
		}

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		args := []any{
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"latency", latency.String(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		}
		if userID, exists := c.Get(ContextUserID); exists {
			args = append(args, "user_id", userID)
		}
		if len(c.Errors) > 0 {
			args = append(args, "errors", c.Errors.String())
		}

		switch {
		case status >= 500:
			log.Error("request", args...)
		case status >= 400:
			log.Warn("request", args...)
		default:
			log.Info("request", args...)
		}
	}
}