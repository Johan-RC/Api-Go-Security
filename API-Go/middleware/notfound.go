package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/shared/response"
)

// NotFound responde un JSON estándar para rutas inexistentes.
func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "recurso no encontrado")
	}
}

// MethodNotAllowed responde un JSON estándar para métodos no permitidos.
func MethodNotAllowed() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Error(c, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "método no permitido")
	}
}
