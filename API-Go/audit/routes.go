package audit

import (
	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/middleware"
	"gorm.io/gorm"
)

// RegisterRoutes registra las rutas del módulo audit (protegidas).
func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB, secret string) {
	repo := NewRepository(db)
	service := NewService(repo)
	h := NewHandler(service)

	protected := middleware.RequireAuth(secret)

	audit := api.Group("/audit", protected)
	audit.GET("/logins", middleware.RequirePermission(db, "AUDIT_LOG_VIEW", "READ"), h.List)
	audit.GET("/logins/users/:userId", middleware.RequirePermission(db, "AUDIT_LOG_VIEW", "READ"), h.ListByUser)
}