package config

import (
	"fmt"

	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	postgresgorm "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
	"github.com/phatnt199/go-infra/pkg/migration"
)

// EnvConfig holds all configuration loaded from environment variables only
type EnvConfig struct {
	// Database configuration
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     int    `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER" envDefault:"postgres"`
	DBPassword string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName     string `env:"DB_NAME" envDefault:"auth_db"`
	DBSSLMode  bool   `env:"DB_SSL_MODE" envDefault:"false"`

	// Migration configuration
	MigrationDir string `env:"MIGRATION_DIR" envDefault:"migrations"`

	// Server configuration
	ServerPort int    `env:"SERVER_PORT" envDefault:"8081"`
	BasePath   string `env:"BASE_PATH" envDefault:"/api/v1"`

	// Logging
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

// LoadConfigFromEnv loads configuration from environment variables only
func LoadConfigFromEnv() (*EnvConfig, error) {
	// Load .env file if it exists (optional)
	_ = godotenv.Load()

	cfg := &EnvConfig{}

	// Parse environment variables into config struct
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment config: %w", err)
	}

	return cfg, nil
}

// ToGormOptions converts the env config to GormOptions
func (c *EnvConfig) ToGormOptions() *postgresgorm.GormOptions {
	return &postgresgorm.GormOptions{
		Type:          postgresgorm.Postgres,
		Host:          c.DBHost,
		Port:          c.DBPort,
		User:          c.DBUser,
		Password:      c.DBPassword,
		DBName:        c.DBName,
		SSLMode:       c.DBSSLMode,
		EnableTracing: false, // Can be added to env config if needed
	}
}

// ToGooseOptions converts the env config to GooseOptions
func (c *EnvConfig) ToMigrationOptions() *migration.MigrationOptions {
	return &migration.MigrationOptions{
		Host:          c.DBHost,
		Port:          c.DBPort,
		User:          c.DBUser,
		Password:      c.DBPassword,
		DBName:        c.DBName,
		SSLMode:       c.DBSSLMode,
		MigrationsDir: c.MigrationDir,
		VersionTable:  "goose_db_version", // Default goose version table
		SkipMigration: false,
	}
}

// PrintConfig prints the current configuration (for debugging)
func (c *EnvConfig) PrintConfig() {
	fmt.Printf("Configuration loaded from environment:\n")
	fmt.Printf("  Database: %s:%d/%s (user: %s)\n", c.DBHost, c.DBPort, c.DBName, c.DBUser)
	fmt.Printf("  Server: port %d, base path %s\n", c.ServerPort, c.BasePath)
	fmt.Printf("  Migrations: %s\n", c.MigrationDir)
	fmt.Printf("  Log Level: %s\n", c.LogLevel)
}
