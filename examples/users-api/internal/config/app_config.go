package config

import (
	"os"
	"strconv"

	httpconfig "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter/config"
	postgresgorm "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
	"github.com/phatnt199/go-infra/pkg/logger"
)

// AppConfig holds all application configuration
type AppConfig struct {
	HTTPConfig *httpconfig.FiberHttpOptions
	DBConfig   *postgresgorm.GormOptions
}

// LoadConfig loads all application configuration from environment variables
func LoadConfig(log logger.Logger) *AppConfig {
	// Load .env file (optional - will use system env if not found)
	// _ = godotenv.Load("../../.env")

	return &AppConfig{
		HTTPConfig: loadHTTPConfig(),
		DBConfig:   loadDatabaseConfig(),
	}
}

// loadHTTPConfig loads HTTP configuration from environment variables
func loadHTTPConfig() *httpconfig.FiberHttpOptions {
	return &httpconfig.FiberHttpOptions{
		Port:        getEnv("FIBER_HTTP_OPTIONS_PORT", ":8080"),
		Host:        getEnv("FIBER_HTTP_OPTIONS_HOST", "localhost"),
		BasePath:    getEnv("FIBER_HTTP_OPTIONS_BASE_PATH", "/api/v1"),
		Name:        getEnv("FIBER_HTTP_OPTIONS_NAME", "Users API"),
		Development: getEnvAsBool("FIBER_HTTP_OPTIONS_DEVELOPMENT", true),
		Timeout:     getEnvAsInt("FIBER_HTTP_OPTIONS_TIMEOUT", 30),
	}
}

// loadDatabaseConfig loads database configuration from environment variables
func loadDatabaseConfig() *postgresgorm.GormOptions {
	return &postgresgorm.GormOptions{
		Type:     postgresgorm.Postgres,
		Host:     getEnv("GORM_OPTIONS_HOST", "localhost"),
		Port:     getEnvAsInt("GORM_OPTIONS_PORT", 5432),
		User:     getEnv("GORM_OPTIONS_USER", "postgres"),
		Password: getEnv("GORM_OPTIONS_PASSWORD", "postgres"),
		DBName:   getEnv("GORM_OPTIONS_DBNAME", "usersdb"),
		SSLMode:  false,
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt gets an environment variable as integer with a default value
func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// getEnvAsBool gets an environment variable as boolean with a default value
func getEnvAsBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1" || value == "yes"
}
