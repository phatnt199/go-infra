package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v8"
)

// Config holds the application configuration
type Config struct {
	HTTP     HTTPConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Logging  LoggingConfig
	Env      string `env:"ENVIRONMENT" envDefault:"development"`
}

// HTTPConfig holds HTTP server configuration
type HTTPConfig struct {
	Port     int    `env:"HTTP_PORT" envDefault:"8080"`
	BasePath string `env:"HTTP_BASE_PATH" envDefault:"/api/v1"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host         string `env:"DB_HOST" envDefault:"localhost"`
	Port         int    `env:"DB_PORT" envDefault:"5432"`
	User         string `env:"DB_USER" envDefault:"postgres"`
	Password     string `env:"DB_PASSWORD" envDefault:"postgres"`
	Name         string `env:"DB_NAME" envDefault:"authenticationdb"`
	SSLMode      string `env:"DB_SSL_MODE" envDefault:"disable"`
	MaxOpenConns int    `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns int    `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret             string        `env:"JWT_SECRET" envDefault:"change-this-secret"`
	Issuer             string        `env:"JWT_ISSUER" envDefault:"authentication-service"`
	Audience           string        `env:"JWT_AUDIENCE" envDefault:"authentication-api"`
	AccessTokenExpiry  time.Duration `env:"JWT_ACCESS_TOKEN_EXPIRY" envDefault:"15m"`
	RefreshTokenExpiry time.Duration `env:"JWT_REFRESH_TOKEN_EXPIRY" envDefault:"168h"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `env:"LOG_LEVEL" envDefault:"info"`
	Format string `env:"LOG_FORMAT" envDefault:"json"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// GetDSN returns the database connection string
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}
