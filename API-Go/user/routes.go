package user

import (
	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/middleware"
	"gorm.io/gorm"
)

// RegisterRoutes registra las rutas del módulo user (todas protegidas).
func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB, secret string) {
	repo := NewRepository(db)
	service := NewService(repo)
	h := NewHandler(service)

	users := api.Group("/users", middleware.RequireAuth(secret))
	users.GET("", middleware.RequirePermission(db, "IDENTITY_USER_VIEW", "READ"), h.List)
	users.POST("", middleware.RequirePermission(db, "IDENTITY_USER_MANAGE", "WRITE"), h.Create)
	users.GET("/:id", middleware.RequirePermission(db, "IDENTITY_USER_VIEW", "READ"), h.GetById)
	users.PUT("/:id", middleware.RequirePermission(db, "IDENTITY_USER_MANAGE", "WRITE"), h.Update)
	users.DELETE("/:id", middleware.RequirePermission(db, "IDENTITY_USER_MANAGE", "WRITE"), h.Delete)
}