package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/config"
	"github.com/jesus-ariel/api-jeussairel/middleware"
	"gorm.io/gorm"
)

// RegisterRoutes registra las rutas del módulo auth.
func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	repo := NewRepository(db)
	service := NewService(repo, cfg)
	h := NewHandler(service)

	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/verify-email", h.VerifyEmail)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.POST("/reset-password", h.ResetPassword)
		auth.GET("/me", middleware.RequireAuth(cfg.JWTSecret), h.Me)
	}
}