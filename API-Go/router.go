package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jesus-ariel/api-jeussairel/audit"
	"github.com/jesus-ariel/api-jeussairel/auth"
	"github.com/jesus-ariel/api-jeussairel/config"
	"github.com/jesus-ariel/api-jeussairel/feature"
	authmw "github.com/jesus-ariel/api-jeussairel/middleware"
	"github.com/jesus-ariel/api-jeussairel/role"
	"github.com/jesus-ariel/api-jeussairel/session"
	"github.com/jesus-ariel/api-jeussairel/user"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	_ "github.com/jesus-ariel/api-jeussairel/docs"
)

// setupRouter construye el router principal con sus middlewares globales
// y monta las rutas de cada módulo.
func setupRouter(db *gorm.DB, cfg *config.Config, log *slog.Logger) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Logger estructurado, recuperación de pánicos y CORS.
	r.Use(authmw.Logger(log))
	r.Use(authmw.Recovery())
	r.Use(authmw.CORS(cfg))

	// Respuestas JSON para rutas no encontradas.
	r.NoRoute(authmw.NotFound())
	r.NoMethod(authmw.MethodNotAllowed())

	api := r.Group("/api/v1")

	// Swagger UI.
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check.
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Registro de rutas por módulo.
	auth.RegisterRoutes(api, db, cfg)
	user.RegisterRoutes(api, db, cfg.JWTSecret)
	role.RegisterRoutes(api, db, cfg.JWTSecret)
	feature.RegisterRoutes(api, db, cfg.JWTSecret)
	session.RegisterRoutes(api, db, cfg.JWTSecret)
	audit.RegisterRoutes(api, db, cfg.JWTSecret)

	return r
}