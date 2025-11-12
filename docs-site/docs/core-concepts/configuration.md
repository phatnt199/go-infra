---
sidebar_position: 4
---

# Configuration Management

Learn how to manage application configuration in go-infra using environment-based settings.

## Overview

go-infra provides a flexible configuration system that supports:

- **Environment-based configs** - Different settings for dev, staging, production
- **Multiple sources** - Files, environment variables, command-line flags
- **Type-safe** - Strongly typed configuration structs
- **Validation** - Built-in validation support
- **Hot reload** - Update config without restart (optional)

## Quick Start

### Define Configuration Structure

```go
// config/config.go
package config

type Config struct {
    Environment string         `json:"environment" env:"APP_ENV" default:"development"`
    Server      ServerConfig   `json:"server"`
    Database    DatabaseConfig `json:"database"`
    JWT         JWTConfig      `json:"jwt"`
    Email       EmailConfig    `json:"email"`
}

type ServerConfig struct {
    Host string `json:"host" env:"SERVER_HOST" default:"localhost"`
    Port int    `json:"port" env:"SERVER_PORT" default:"3000"`
}

type DatabaseConfig struct {
    Host     string `json:"host" env:"DB_HOST" default:"localhost"`
    Port     int    `json:"port" env:"DB_PORT" default:"5432"`
    Username string `json:"username" env:"DB_USER" default:"postgres"`
    Password string `json:"password" env:"DB_PASSWORD"`
    Database string `json:"database" env:"DB_NAME" default:"myapp"`
    SSLMode  string `json:"ssl_mode" env:"DB_SSL_MODE" default:"disable"`
}

type JWTConfig struct {
    Secret     string `json:"secret" env:"JWT_SECRET"`
    Expiration int    `json:"expiration" env:"JWT_EXPIRATION" default:"24"` // hours
}

type EmailConfig struct {
    Host     string `json:"host" env:"EMAIL_HOST"`
    Port     int    `json:"port" env:"EMAIL_PORT" default:"587"`
    Username string `json:"username" env:"EMAIL_USER"`
    Password string `json:"password" env:"EMAIL_PASSWORD"`
    From     string `json:"from" env:"EMAIL_FROM"`
}
```

### Load Configuration

```go
// config/loader.go
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "github.com/joho/godotenv"
)

func Load() (*Config, error) {
    // Load .env file if exists
    _ = godotenv.Load()

    // Get environment
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "development"
    }

    // Load config file
    configFile := fmt.Sprintf("config.%s.json", env)
    cfg, err := loadFromFile(configFile)
    if err != nil {
        return nil, err
    }

    // Override with environment variables
    overrideFromEnv(cfg)

    // Validate
    if err := validate(cfg); err != nil {
        return nil, err
    }

    return cfg, nil
}

func loadFromFile(filename string) (*Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

func overrideFromEnv(cfg *Config) {
    if val := os.Getenv("SERVER_PORT"); val != "" {
        fmt.Sscanf(val, "%d", &cfg.Server.Port)
    }
    if val := os.Getenv("DB_HOST"); val != "" {
        cfg.Database.Host = val
    }
    // ... more overrides
}

func validate(cfg *Config) error {
    if cfg.JWT.Secret == "" {
        return fmt.Errorf("JWT secret is required")
    }
    if cfg.Database.Password == "" {
        return fmt.Errorf("database password is required")
    }
    return nil
}
```

### Use Configuration

```go
// main.go
package main

import (
    "log"
    "myapp/config"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Use configuration
    log.Printf("Starting server on %s:%d", cfg.Server.Host, cfg.Server.Port)

    // Setup database with config
    db := setupDatabase(cfg.Database)

    // Start server
    app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
}
```

## Configuration Files

### Development Config

```json
// config.development.json
{
	"environment": "development",
	"server": {
		"host": "localhost",
		"port": 3000
	},
	"database": {
		"host": "localhost",
		"port": 5432,
		"username": "postgres",
		"password": "postgres",
		"database": "myapp_dev",
		"ssl_mode": "disable"
	},
	"jwt": {
		"secret": "dev-secret-key-change-in-production",
		"expiration": 24
	},
	"email": {
		"host": "smtp.mailtrap.io",
		"port": 2525,
		"username": "your-mailtrap-username",
		"password": "your-mailtrap-password",
		"from": "noreply@myapp.dev"
	}
}
```

### Production Config

```json
// config.production.json
{
	"environment": "production",
	"server": {
		"host": "0.0.0.0",
		"port": 8080
	},
	"database": {
		"host": "${DB_HOST}",
		"port": 5432,
		"username": "${DB_USER}",
		"password": "${DB_PASSWORD}",
		"database": "${DB_NAME}",
		"ssl_mode": "require"
	},
	"jwt": {
		"secret": "${JWT_SECRET}",
		"expiration": 1
	},
	"email": {
		"host": "${EMAIL_HOST}",
		"port": 587,
		"username": "${EMAIL_USER}",
		"password": "${EMAIL_PASSWORD}",
		"from": "${EMAIL_FROM}"
	}
}
```

## Environment Variables

### .env File

```bash
# .env
APP_ENV=development

# Server
SERVER_HOST=localhost
SERVER_PORT=3000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=myapp_dev
DB_SSL_MODE=disable

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=24

# Email
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USER=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
EMAIL_FROM=noreply@myapp.com
```

### Loading .env File

```go
import "github.com/joho/godotenv"

func init() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
}
```

## Using Viper

For more advanced configuration management, use Viper:

```go
// config/viper.go
package config

import (
    "github.com/spf13/viper"
)

func LoadWithViper() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("json")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")

    // Enable environment variable override
    viper.AutomaticEnv()

    // Read config file
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    // Unmarshal into struct
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

// Watch for config changes
func WatchConfig(onChange func(*Config)) {
    viper.WatchConfig()
    viper.OnConfigChange(func(e fsnotify.Event) {
        var cfg Config
        viper.Unmarshal(&cfg)
        onChange(&cfg)
    })
}
```

## go-infra Environment Package

Use the built-in environment package:

```go
import (
    "github.com/phatnt199/go-infra/pkg/application/environment"
    "github.com/phatnt199/go-infra/pkg/application/config"
)

func main() {
    // Load environment
    env.Load()

    // Get environment value
    dbHost := env.GetString("DB_HOST", "localhost")
    dbPort := env.GetInt("DB_PORT", 5432)

    // Get required value (panics if not found)
    jwtSecret := env.MustGetString("JWT_SECRET")

    // Check environment
    if env.IsDevelopment() {
        log.Println("Running in development mode")
    }

    if env.IsProduction() {
        log.Println("Running in production mode")
    }
}
```

## Configuration Validation

### Using struct tags

```go
import "github.com/go-playground/validator/v10"

type Config struct {
    Server   ServerConfig   `validate:"required"`
    Database DatabaseConfig `validate:"required"`
}

type ServerConfig struct {
    Host string `validate:"required,hostname"`
    Port int    `validate:"required,min=1,max=65535"`
}

func validate(cfg *Config) error {
    validate := validator.New()
    return validate.Struct(cfg)
}
```

### Custom validation

```go
func validate(cfg *Config) error {
    if cfg.JWT.Secret == "" || len(cfg.JWT.Secret) < 32 {
        return fmt.Errorf("JWT secret must be at least 32 characters")
    }

    if cfg.Database.Password == "" {
        return fmt.Errorf("database password is required")
    }

    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
    }

    return nil
}
```

## Secrets Management

### Using Environment Variables

```go
// Never commit secrets to git
// Use environment variables in production

// .env.example (commit this)
JWT_SECRET=your-secret-here
DB_PASSWORD=your-password-here

// .env (do NOT commit this)
JWT_SECRET=actual-production-secret
DB_PASSWORD=actual-production-password
```

### Using AWS Secrets Manager

```go
import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/secretsmanager"
)

func loadSecretsFromAWS(secretName string) (map[string]string, error) {
    sess := session.Must(session.NewSession())
    svc := secretsmanager.New(sess)

    input := &secretsmanager.GetSecretValueInput{
        SecretId: aws.String(secretName),
    }

    result, err := svc.GetSecretValue(input)
    if err != nil {
        return nil, err
    }

    var secrets map[string]string
    json.Unmarshal([]byte(*result.SecretString), &secrets)

    return secrets, nil
}
```

## Feature Flags

Implement feature flags in configuration:

```go
type Config struct {
    Features FeatureFlags `json:"features"`
}

type FeatureFlags struct {
    EnableNewUI      bool `json:"enable_new_ui" env:"FEATURE_NEW_UI"`
    EnableBetaAPI    bool `json:"enable_beta_api" env:"FEATURE_BETA_API"`
    EnableRateLimiting bool `json:"enable_rate_limiting" env:"FEATURE_RATE_LIMITING"`
}

// Usage
func (h *Handler) GetUser(c *fiber.Ctx) error {
    if h.config.Features.EnableNewUI {
        return h.getUserNewUI(c)
    }
    return h.getUserLegacy(c)
}
```

## Multi-Environment Setup

### Directory Structure

```
config/
├── config.go              # Config structs
├── loader.go              # Loading logic
├── config.development.json
├── config.staging.json
├── config.production.json
└── .env.example
```

### Environment Detection

```go
func getEnvironment() string {
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = os.Getenv("GO_ENV")
    }
    if env == "" {
        env = "development"
    }
    return env
}

func isDevelopment() bool {
    return getEnvironment() == "development"
}

func isProduction() bool {
    return getEnvironment() == "production"
}
```

## Best Practices

### 1. Never Commit Secrets

```bash
# .gitignore
.env
config.local.json
*.secret.json
```

### 2. Use Defaults

```go
type Config struct {
    Port int `json:"port" default:"3000"`
}
```

### 3. Document Configuration

```go
// Config holds application configuration
type Config struct {
    // Server configuration
    Server ServerConfig `json:"server"`

    // Database connection settings
    Database DatabaseConfig `json:"database"`

    // JWT token configuration
    // Secret must be at least 32 characters
    JWT JWTConfig `json:"jwt"`
}
```

### 4. Validate Early

```go
func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Configuration error:", err)
    }

    // Validate immediately
    if err := cfg.Validate(); err != nil {
        log.Fatal("Invalid configuration:", err)
    }

    // Now safe to use
}
```

### 5. Use Type-Safe Access

```go
// ✅ Good - type-safe
cfg.Server.Port  // int

// ❌ Bad - string parsing
port, _ := strconv.Atoi(os.Getenv("PORT"))
```

## Example: Complete Configuration Setup

```go
// config/config.go
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    Environment string         `json:"environment"`
    Server      ServerConfig   `json:"server"`
    Database    DatabaseConfig `json:"database"`
    Redis       RedisConfig    `json:"redis"`
    JWT         JWTConfig      `json:"jwt"`
    Email       EmailConfig    `json:"email"`
    Features    FeatureFlags   `json:"features"`
}

type ServerConfig struct {
    Host         string `json:"host" env:"SERVER_HOST" default:"localhost"`
    Port         int    `json:"port" env:"SERVER_PORT" default:"3000"`
    ReadTimeout  int    `json:"read_timeout" default:"30"`
    WriteTimeout int    `json:"write_timeout" default:"30"`
}

type DatabaseConfig struct {
    Host            string `json:"host" env:"DB_HOST" default:"localhost"`
    Port            int    `json:"port" env:"DB_PORT" default:"5432"`
    Username        string `json:"username" env:"DB_USER" default:"postgres"`
    Password        string `json:"password" env:"DB_PASSWORD"`
    Database        string `json:"database" env:"DB_NAME" default:"myapp"`
    SSLMode         string `json:"ssl_mode" env:"DB_SSL_MODE" default:"disable"`
    MaxOpenConns    int    `json:"max_open_conns" default:"25"`
    MaxIdleConns    int    `json:"max_idle_conns" default:"5"`
    ConnMaxLifetime int    `json:"conn_max_lifetime" default:"300"`
}

type RedisConfig struct {
    Host     string `json:"host" env:"REDIS_HOST" default:"localhost"`
    Port     int    `json:"port" env:"REDIS_PORT" default:"6379"`
    Password string `json:"password" env:"REDIS_PASSWORD"`
    DB       int    `json:"db" env:"REDIS_DB" default:"0"`
}

type JWTConfig struct {
    Secret          string `json:"secret" env:"JWT_SECRET"`
    AccessExpiration  int  `json:"access_expiration" default:"15"`   // minutes
    RefreshExpiration int  `json:"refresh_expiration" default:"168"` // hours
}

type EmailConfig struct {
    Host     string `json:"host" env:"EMAIL_HOST"`
    Port     int    `json:"port" env:"EMAIL_PORT" default:"587"`
    Username string `json:"username" env:"EMAIL_USER"`
    Password string `json:"password" env:"EMAIL_PASSWORD"`
    From     string `json:"from" env:"EMAIL_FROM"`
}

type FeatureFlags struct {
    EnableNewUI     bool `json:"enable_new_ui" env:"FEATURE_NEW_UI"`
    EnableBetaAPI   bool `json:"enable_beta_api" env:"FEATURE_BETA_API"`
    EnableCaching   bool `json:"enable_caching" env:"FEATURE_CACHING"`
}

var instance *Config

func Load() (*Config, error) {
    if instance != nil {
        return instance, nil
    }

    // Load .env
    _ = godotenv.Load()

    // Get environment
    env := os.Getenv("APP_ENV")
    if env == "" {
        env = "development"
    }

    // Load config file
    filename := fmt.Sprintf("config.%s.json", env)
    cfg, err := loadFromFile(filename)
    if err != nil {
        return nil, err
    }

    cfg.Environment = env

    // Override with env vars
    overrideFromEnv(cfg)

    // Validate
    if err := cfg.Validate(); err != nil {
        return nil, err
    }

    instance = cfg
    return cfg, nil
}

func (c *Config) Validate() error {
    if c.JWT.Secret == "" || len(c.JWT.Secret) < 32 {
        return fmt.Errorf("JWT secret must be at least 32 characters")
    }
    if c.Database.Password == "" {
        return fmt.Errorf("database password is required")
    }
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }
    return nil
}

func (c *Config) IsProduction() bool {
    return c.Environment == "production"
}

func (c *Config) IsDevelopment() bool {
    return c.Environment == "development"
}
```

## Next Steps

- Learn about [Error Handling](./error-handling)
- Explore [Module System](./modules)
- See [Deployment Guide](../deployment/production)
