package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/shared/response"
)

// Recovery captura pánicos, los registra y responde 500 sin matar el server.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic capturado",
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"stack", string(debug.Stack()),
				)
				response.Error(c, http.StatusInternalServerError, "INTERNAL", "error interno del servidor")
				c.Abort()
			}
		}()
		c.Next()
	}
}
