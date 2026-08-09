// @title API JeussAirel (IAM)
// @version 1.0
// @description API REST de Identity and Access Management (IAM): autenticación JWT, RBAC, usuarios, roles, funcionalidades, sesiones y auditoría.

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Token JWT de acceso. Uso: `Authorization: Bearer <access_token>`
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jesus-ariel/api-jeussairel/config"
	"github.com/jesus-ariel/api-jeussairel/database"
	"github.com/jesus-ariel/api-jeussairel/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.AppEnv, os.Getenv("LOG_LEVEL"))

	db, err := database.Connect(cfg)
	if err != nil {
		log.Error("error conectando a la base de datos", "error", err)
		os.Exit(1)
	}
	log.Info("conexión a la base de datos establecida")

	r := setupRouter(db, cfg, log)

	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Info("servidor iniciado", "port", cfg.AppPort, "env", cfg.AppEnv)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("error iniciando el servidor", "error", err)
			os.Exit(1)
		}
	}()

	// Espera señales de cierre para apagar de forma ordenada.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("apagando servidor...")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownGrace)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("error en el apagado del servidor", "error", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}

	log.Info("servidor detenido")
}
