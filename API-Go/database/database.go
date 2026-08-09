package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jesus-ariel/api-jeussairel/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// DB es el handle global de GORM. Se inicializa en Connect.
var DB *gorm.DB

// Connect abre la conexión a PostgreSQL usando GORM y configura el pool.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(sqlLogMode(cfg)),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	if err := pingWithRetry(sqlDB); err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}

// PingWithRetry intenta verificar la conexión varias veces.
func pingWithRetry(sqlDB interface {
	PingContext(ctx context.Context) error
}) error {
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		lastErr = sqlDB.PingContext(ctx)
		cancel()
		if lastErr == nil {
			return nil
		}
		slog.Warn("reintentando conexión a la base de datos", "attempt", attempt, "error", lastErr)
		time.Sleep(2 * time.Second)
	}
	return lastErr
}

func sqlLogMode(cfg *config.Config) gormlogger.LogLevel {
	switch cfg.AppEnv {
	case "production":
		return gormlogger.Error
	case "develop":
		return gormlogger.Info
	default:
		return gormlogger.Warn
	}
}
