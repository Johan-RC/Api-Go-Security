package session

import (
	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/middleware"
	"gorm.io/gorm"
)

// RegisterRoutes registra las rutas del módulo session (protegidas).
func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB, secret string) {
	repo := NewRepository(db)
	service := NewService(repo)
	h := NewHandler(service)

	protected := middleware.RequireAuth(secret)

	sessions := api.Group("/sessions", protected)
	sessions.GET("/active", h.ListActiveByUser)
	sessions.DELETE("/:id", h.RevokeById)
}