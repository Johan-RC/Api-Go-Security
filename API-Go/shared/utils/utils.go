package utils

import (
	"github.com/gin-gonic/gin"
)

// GetUserID obtiene el ID del usuario autenticado desde el contexto.
// El valor lo inyectará el middleware de autenticación en pasos posteriores.
func GetUserID(c *gin.Context) string {
	userID, _ := c.Get("user_id")
	id, _ := userID.(string)
	return id
}

// GetClientIP devuelve la IP del cliente, respetando proxies.
func GetClientIP(c *gin.Context) string {
	return c.ClientIP()
}
