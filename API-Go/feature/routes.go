package feature

import (
	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/middleware"
	"gorm.io/gorm"
)

// RegisterRoutes registra las rutas del módulo feature (todas protegidas).
func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB, secret string) {
	repo := NewRepository(db)
	service := NewService(repo)
	h := NewHandler(service)

	protected := middleware.RequireAuth(secret)

	modules := api.Group("/modules", protected)
	modules.GET("", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.ListModules)
	modules.POST("", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.CreateModule)
	modules.GET("/:id", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.GetModuleById)
	modules.PUT("/:id", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.UpdateModule)
	modules.DELETE("/:id", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.DeleteModule)
	modules.GET("/:id/features", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.ListFeaturesByModule)

	features := api.Group("/features", protected)
	features.GET("", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.ListFeatures)
	features.POST("", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.CreateFeature)
	features.GET("/:id", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.GetFeatureById)
	features.PUT("/:id", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.UpdateFeature)
	features.DELETE("/:id", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.DeleteFeature)
}