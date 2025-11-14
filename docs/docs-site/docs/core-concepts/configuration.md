---
sidebar_position: 4
---

# Configuration Management

Learn how to manage application configuration in go-infra using Viper-powered, environment-aware configuration system.

## Overview

go-infra provides a powerful, production-ready configuration system built on Viper that supports:

- **Environment-aware configs** - Load different `config.<env>.json` files for development, staging, production
- **Multiple sources** - JSON/YAML files + environment variables with override capability
- **Automatic discovery** - Recursively finds `.env` files and config files in your project
- **Type-safe** - Strongly typed configuration structs with validation
- **Smart defaults** - Uses `go-defaults` package for sensible fallbacks
- **Flexible binding** - Load entire config or specific sections with `BindConfigKey`

## Key Concepts

go-infra provides **two complementary approaches** to configuration:

### 1. File-Based Configuration (Recommended for Most Apps)

Uses **Viper** to load structured `config.<env>.json` files with environment variable overrides. This is what most go-infra examples and modules use internally.

**Best for:**

- Applications with complex configuration needs
- Multiple deployment environments (dev, staging, prod)
- Configuration shared across modules (HTTP, database, auth, etc.)

### 2. Environment-Only Configuration

Uses only environment variables without config files. Useful for containerized deployments following 12-factor principles.

**Best for:**

- Simple microservices
- Docker/Kubernetes deployments
- When you prefer env vars over files

## How It Works

### Environment Bootstrap (`pkg/application/environment`)

The configuration system starts by bootstrapping the environment:

```go
import "github.com/phatnt199/go-infra/pkg/application/environment"

// Called automatically by fxapp.NewApplicationBuilder()
env := environment.ConfigAppEnv()
// Returns: environment.Development, environment.Production, or environment.Staging
```

**What it does:**

1. **Loads `.env` files** - Recursively searches from current directory upward using `godotenv`
2. **Resolves `APP_ENV`** - Reads from environment or defaults to `development`
3. **Sets `APP_ROOT_PATH`** - Finds project root by locating `go.mod` or using `APP_NAME`
4. **Fixes working directory** - Calls `FixProjectRootWorkingDirectoryPath()` for consistent file access

### Config File Binding (`pkg/application/config`)

Once environment is bootstrapped, bind configuration from files:

```go
import "github.com/phatnt199/go-infra/pkg/application/config"

// Load specific config section
type FiberHttpOptions struct {
    Port        string `mapstructure:"port" default:":8080"`
    Host        string `mapstructure:"host" default:"localhost"`
    BasePath    string `mapstructure:"basePath" default:"/api/v1"`
    Development bool   `mapstructure:"development" default:"true"`
}

func ProvideFiberConfig(env environment.Environment) (*FiberHttpOptions, error) {
    // Automatically loads from config.development.json (or config.production.json, etc.)
    return config.BindConfigKey[*FiberHttpOptions]("fiberHttpOptions", env)
}
```

**Loading sequence:**

1. **Set defaults** - Uses `go-defaults` package to apply `default` struct tags
2. **Find config file** - Searches for `config.<env>.json` starting from `APP_ROOT_PATH`
3. **Load with Viper** - Unmarshals file into struct using `mapstructure` tags
4. **Apply env overrides** - Runs `viper.AutomaticEnv()` and `env.Parse()` to apply environment variables
5. **Return typed config** - Returns strongly-typed configuration struct

## File-Based Configuration (Recommended)

### Step 1: Create Config Files

Create environment-specific config files in your project (typically in `config/` or `internal/config/`):

```json title="config/config.development.json"
{
	"fiberHttpOptions": {
		"port": ":8080",
		"host": "localhost",
		"basePath": "/api/v1",
		"name": "My API",
		"development": true,
		"timeout": 30
	},
	"gormOptions": {
		"type": 0,
		"host": "localhost",
		"port": 5432,
		"user": "postgres",
		"password": "postgres",
		"dbname": "myapp_dev",
		"sslmode": false
	},
	"logOptions": {
		"level": "debug",
		"logType": 0,
		"callerEnabled": true
	}
}
```

```json title="config/config.production.json"
{
	"fiberHttpOptions": {
		"port": ":8080",
		"host": "0.0.0.0",
		"basePath": "/api/v1",
		"name": "My API",
		"development": false,
		"timeout": 60
	},
	"gormOptions": {
		"type": 0,
		"host": "db.example.com",
		"port": 5432,
		"user": "prod_user",
		"password": "${DB_PASSWORD}",
		"dbname": "myapp_prod",
		"sslmode": true
	},
	"logOptions": {
		"level": "info",
		"logType": 0,
		"callerEnabled": false
	}
}
```

### Step 2: Define Configuration Structs

```go title="internal/config/app_options.go"
package config

import (
    "github.com/phatnt199/go-infra/pkg/application/config"
    "github.com/phatnt199/go-infra/pkg/application/environment"
)

type AppOptions struct {
    ServiceName  string `mapstructure:"serviceName" default:"myapp"`
    DeliveryType string `mapstructure:"deliveryType" default:"http"`
}

type FiberHttpOptions struct {
    Host        string `mapstructure:"host" default:"localhost"`
    Port        string `mapstructure:"port" default:":8080"`
    BasePath    string `mapstructure:"basePath" default:"/api/v1"`
    Name        string `mapstructure:"name" default:"My API"`
    Development bool   `mapstructure:"development" default:"true"`
    Timeout     int    `mapstructure:"timeout" default:"30"`
}

type GormOptions struct {
    Type     int    `mapstructure:"type" default:"0"`
    Host     string `mapstructure:"host" default:"localhost"`
    Port     int    `mapstructure:"port" default:"5432"`
    User     string `mapstructure:"user" env:"DB_USER" default:"postgres"`
    Password string `mapstructure:"password" env:"DB_PASSWORD"`
    DBName   string `mapstructure:"dbname" env:"DB_NAME" default:"myapp"`
    SSLMode  bool   `mapstructure:"sslmode" default:"false"`
}
```

:::tip Struct Tags Explained

- `mapstructure:"fieldName"` - Maps JSON field to struct field (used by Viper)
- `env:"ENV_VAR"` - Maps environment variable to field (used by `env.Parse`)
- `default:"value"` - Default value if not set (used by `go-defaults`)
  :::

### Step 3: Create Provider Functions

```go title="internal/config/config_fx.go"
package config

import (
    "github.com/phatnt199/go-infra/pkg/application/config"
    "github.com/phatnt199/go-infra/pkg/application/environment"
    "go.uber.org/fx"
)

// Module provides configuration to fx dependency injection
var Module = fx.Module(
    "config",
    fx.Provide(
        ProvideAppOptions,
        ProvideFiberOptions,
        ProvideGormOptions,
    ),
)

func ProvideAppOptions(env environment.Environment) (*AppOptions, error) {
    return config.BindConfigKey[*AppOptions]("appOptions", env)
}

func ProvideFiberOptions(env environment.Environment) (*FiberHttpOptions, error) {
    return config.BindConfigKey[*FiberHttpOptions]("fiberHttpOptions", env)
}

func ProvideGormOptions(env environment.Environment) (*GormOptions, error) {
    return config.BindConfigKey[*GormOptions]("gormOptions", env)
}
```

### Step 4: Use Configuration in Your App

```go title="cmd/api/main.go"
package main

import (
    "myapp/internal/config"
    "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
    customfiber "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
    postgresgorm "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
)

func main() {
    app := fxapp.NewApplicationBuilder().
        ProvideModule(config.Module).        // Load config first
        ProvideModule(customfiber.Module).   // Uses FiberHttpOptions
        ProvideModule(postgresgorm.Module).  // Uses GormOptions
        Build()

    app.Run()
}
```

## Environment Variable Overrides

Config files are loaded first, then environment variables override them:

```bash
# Override database password from environment
export DB_PASSWORD="production-secret"

# Override log level
export LOG_LEVEL="debug"

# Run application - config file values are used, but DB_PASSWORD is overridden
go run cmd/api/main.go
```

**Priority (highest to lowest):**

1. Environment variables with matching `env` tag
2. Values in `config.<env>.json` file
3. Default values from `default` tag

## Configuration Discovery

### How Config Files Are Found

go-infra automatically searches for config files:

1. **Check `CONFIG_PATH`** - If set, loads from this directory directly
2. **Search from `APP_ROOT_PATH`** - Recursively searches subdirectories
3. **Look for** `config.<env>.json`, `config.<env>.yaml`, or `config.<env>.yml`

**Recommended approach for production:**

```bash
# Set CONFIG_PATH to avoid filesystem search
export CONFIG_PATH=/app/config
export APP_ENV=production
./myapp
```

### How `.env` Files Are Found

The `environment.ConfigAppEnv()` function searches recursively:

```
myapp/
├── .env                          # Found! (project root)
├── cmd/
│   └── api/
│       └── main.go              # Starts here, searches upward
└── config/
    └── config.development.json
```

## Environment-Only Configuration (Optional)

If you prefer not to use config files, go-infra provides an environment-only loader:

```go
import "github.com/phatnt199/go-infra/pkg/application/config"

func main() {
    // Load configuration entirely from environment variables
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Use typed configuration
    log.Printf("Starting %s on %s", cfg.App.Name, cfg.Server.HTTP.Address())

    // Check environment
    if cfg.App.IsProduction() {
        // Production-specific logic
    }
}
```

**Supported environment variables:**

See the comprehensive list in [`pkg/application/config/README.md`](https://github.com/phatnt199/go-infra/blob/main/pkg/application/config/README.md) including:

- `APP_*` - Application settings
- `HTTP_*` / `GRPC_*` - Server configuration
- `DB_*` - Database connection
- `REDIS_*` - Redis configuration
- `LOG_*` - Logging settings
- `JWT_*` - Authentication
- And many more...

## Complete Example

Here's a real-world example from the `examples/users-api`:

```go title="cmd/api/main.go"
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/swagger"
    "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
    "github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
    customfiber "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
    postgresgorm "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
    "myapp/internal/modules"
)

func main() {
    // Create application builder (automatically calls environment.ConfigAppEnv())
    appBuilder := fxapp.NewApplicationBuilder()

    // Register modules (they load their own config via BindConfigKey)
    appBuilder.ProvideModule(customfiber.Module)   // Loads fiberHttpOptions
    appBuilder.ProvideModule(postgresgorm.Module)  // Loads gormOptions
    appBuilder.ProvideModule(modules.Module)       // Your business logic

    // Build and run
    app := appBuilder.Build()
    app.Run()
}
```

**Directory structure:**

```
myapp/
├── .env                                # APP_ENV=development
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   └── config/
│       ├── config.development.json    # Loaded when APP_ENV=development
│       └── config.production.json     # Loaded when APP_ENV=production
└── go.mod
```

## Best Practices

### ✅ Do

- **Use `config.<env>.json` files** for structured configuration
- **Set `CONFIG_PATH` in production** to avoid filesystem scanning
- **Use environment variables for secrets** (passwords, API keys)
- **Use `default` tags** for sensible defaults
- **Validate configuration** on startup (fail fast)
- **Create separate files** for each environment

### ❌ Don't

- **Don't commit secrets** to config files (use env vars or secret managers)
- **Don't hardcode values** that differ between environments
- **Don't skip validation** - catch config errors early
- **Don't rely on filesystem search in production** - set `CONFIG_PATH`

## Troubleshooting

### Config file not found

```
Error: viper.ReadInConfig: No directory with config file found
```

**Solution:** Set `CONFIG_PATH` to the directory containing your config file:

```bash
export CONFIG_PATH=$(pwd)/internal/config
go run cmd/api/main.go
```

### Environment variable not overriding config

**Problem:** Environment variable not taking effect

**Solution:** Add `env` tag to struct field:

```go
type GormOptions struct {
    Password string `mapstructure:"password" env:"DB_PASSWORD"`  // Add env tag
}
```

### Application can't find `.env` file

**Problem:** `.env` file exists but not loaded

**Solution:** Ensure `.env` is in project root or ancestor directory. The loader searches upward from the current working directory.

### Working directory issues

**Problem:** Config/migration paths not working

**Solution:** The framework calls `FixProjectRootWorkingDirectoryPath()` which changes the working directory to project root. This ensures consistent file access regardless of where you run your app from.

## Advanced Topics

### Custom Config Sections

Add your own config sections by creating new struct types and provider functions:

```go
type MyCustomConfig struct {
    Feature1 bool   `mapstructure:"feature1" default:"false"`
    Feature2 string `mapstructure:"feature2" default:"default-value"`
}

func ProvideMyCustomConfig(env environment.Environment) (*MyCustomConfig, error) {
    return config.BindConfigKey[*MyCustomConfig]("myCustomConfig", env)
}
```

Then add to your `config.development.json`:

```json
{
	"myCustomConfig": {
		"feature1": true,
		"feature2": "custom-value"
	}
}
```

### Testing with Mock Configuration

```go
func TestWithConfig(t *testing.T) {
    // Create test config
    testCfg := &config.FiberHttpOptions{
        Port: ":9999",
        Host: "localhost",
        Development: true,
    }

    // Use in test
    // ... test code
}
```

## Related Examples

- **[`examples/env-usage`](https://github.com/phatnt199/go-infra/tree/main/examples/env-usage)** - Comprehensive configuration examples
- **[`examples/users-api`](https://github.com/phatnt199/go-infra/tree/main/examples/users-api)** - Real-world API with config files
- **[`examples/authentication-service`](https://github.com/phatnt199/go-infra/tree/main/examples/authentication-service)** - Complex multi-module configuration

## Learn More

- **[pkg/application/config/README.md](https://github.com/phatnt199/go-infra/blob/main/pkg/application/config/README.md)** - Complete config API reference
- **[Viper Documentation](https://github.com/spf13/viper)** - Underlying config library
- **[go-defaults](https://github.com/mcuadros/go-defaults)** - Default value handling
