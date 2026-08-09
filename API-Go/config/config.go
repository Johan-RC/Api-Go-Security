package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jesus-ariel/api-jeussairel/mailer"
	"github.com/joho/godotenv"
)

// Config concentra toda la configuración de la aplicación.
type Config struct {
	// Servidor
	AppEnv        string
	AppPort       string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
	ShutdownGrace time.Duration

	// CORS
	AllowedOrigins []string

	// Base de datos
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration

	// JWT
	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration

	// Correo (SMTP)
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
}

// Load lee las variables de entorno desde el archivo .env (si existe)
// y construye la Config con sus valores por defecto.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("advertencia: no se encontró archivo .env")
	}

	return &Config{
		AppEnv:        getEnv("APP_ENV", "develop"),
		AppPort:       getEnv("APP_PORT", "8080"),
		ReadTimeout:   getDuration("SERVER_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:  getDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:   getDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		ShutdownGrace: getDuration("SERVER_SHUTDOWN_GRACE", 10*time.Second),

		AllowedOrigins: getList("CORS_ALLOWED_ORIGINS", []string{"*"}),

		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBName:            getEnv("DB_NAME", "postgres"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		DBMaxOpenConns:    getInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    getInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime: getDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),

		JWTSecret:     getEnv("JWT_SECRET", ""),
		JWTAccessTTL:  getDuration("JWT_ACCESS_TTL", 15*time.Minute),
		JWTRefreshTTL: getDuration("JWT_REFRESH_TTL", 168*time.Hour),

		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", ""),
	}
}

// SMTP devuelve la configuración del envío de correos.
func (c *Config) SMTP() mailer.Config {
	return mailer.Config{
		Host:     c.SMTPHost,
		Port:     c.SMTPPort,
		Username: c.SMTPUser,
		Password: c.SMTPPassword,
		From:     c.SMTPFrom,
	}
}

// IsProduction indica si la app corre en producción.
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

// CorsAllowedOrigin devuelve el origen permitido si está en la lista.
func (c *Config) CorsAllowedOrigin(origin string) string {
	if origin == "" {
		return ""
	}
	for _, allowed := range c.AllowedOrigins {
		if allowed == "*" || allowed == origin {
			return origin
		}
	}
	return ""
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return fallback
}

func getInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	var n int
	if _, err := fmt.Sscanf(value, "%d", &n); err != nil {
		return fallback
	}
	return n
}

func getList(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	var out []string
	for _, item := range strings.Split(value, ",") {
		if trimmed := strings.TrimSpace(item); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
