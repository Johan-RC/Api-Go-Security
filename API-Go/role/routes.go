package role

import (
	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/feature"
	"github.com/jesus-ariel/api-jeussairel/middleware"
	"github.com/jesus-ariel/api-jeussairel/user"
	"gorm.io/gorm"
)

// RegisterRoutes registra las rutas del módulo role (todas protegidas).
func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB, secret string) {
	repo := NewRepository(db)
	featuresRepo := feature.NewRepository(db)
	usersRepo := user.NewRepository(db)
	service := NewService(repo, featuresRepo, usersRepo)
	h := NewHandler(service)

	protected := middleware.RequireAuth(secret)

	roles := api.Group("/roles", protected)
	roles.GET("", middleware.RequirePermission(db, "IDENTITY_ROLE_VIEW", "READ"), h.List)
	roles.POST("", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.Create)
	roles.GET("/:id", middleware.RequirePermission(db, "IDENTITY_ROLE_VIEW", "READ"), h.GetById)
	roles.PUT("/:id", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.Update)
	roles.DELETE("/:id", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.Delete)
	roles.GET("/:id/features", middleware.RequirePermission(db, "IDENTITY_ROLE_VIEW", "READ"), h.ListFeatures)
	roles.PUT("/:id/features", middleware.RequirePermission(db, "IDENTITY_ROLE_MANAGE", "WRITE"), h.AssignFeatures)

	usersRoles := api.Group("/users", protected)
	usersRoles.POST("/:id/roles", middleware.RequirePermission(db, "IDENTITY_ROLE_ASSIGN", "WRITE"), h.AssignRoleToUser)
	usersRoles.GET("/:id/roles", middleware.RequirePermission(db, "IDENTITY_ROLE_VIEW", "READ"), h.ListRolesByUser)
	usersRoles.DELETE("/:id/roles/:userRoleId", middleware.RequirePermission(db, "IDENTITY_ROLE_ASSIGN", "WRITE"), h.RemoveRoleFromUser)

	overrides := api.Group("/scope-overrides", protected)
	overrides.POST("", middleware.RequirePermission(db, "IDENTITY_SCOPE_MANAGE", "WRITE"), h.CreateScopeOverride)
	overrides.GET("/users/:id", middleware.RequirePermission(db, "IDENTITY_SCOPE_MANAGE", "WRITE"), h.ListScopeOverrides)
	overrides.DELETE("/:id", middleware.RequirePermission(db, "IDENTITY_SCOPE_MANAGE", "WRITE"), h.RemoveScopeOverride)
}