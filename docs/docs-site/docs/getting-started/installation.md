---
sidebar_position: 1
---

# Installation

## Prerequisites

Before you begin, ensure you have:

- **Go 1.21+** installed
- **PostgreSQL** (if using database features)
- **Git** for version control

Check your Go version:

```bash
go version
```

## Project Structure

Create a new project with this recommended structure:

```
myapp/
├── cmd/
│   ├── api/
│   │   └── main.go          # HTTP server entry point
│   └── migration/
│       └── main.go          # Migration runner (optional)
├── internal/
│   ├── domain/              # Domain models
│   ├── repository/          # Data access layer
│   ├── handler/             # HTTP handlers
│   ├── modules/             # Fx modules for DI
│   └── config/              # Config structs (optional)
├── config/                  # Config files
│   ├── config.development.json
│   ├── config.staging.json
│   └── config.production.json
├── migrations/              # Database migrations (optional)
├── .env.example             # Environment template
├── .env                     # Environment variables (git-ignored)
├── go.mod
└── go.sum
```

Create this structure:

```bash
mkdir -p myapp/{cmd/{api,migration},internal/{domain,repository,handler,modules,config},migrations,config}
cd myapp
go mod init myapp
```

## Install go-infra

### Option 1: Add to Existing Project

Add go-infra to your `go.mod`:

```bash
go get github.com/phatnt199/go-infra
```

### Option 2: Local Development

For local development or customization:

```bash
# Clone the repository
git clone https://github.com/phatnt199/go-infra.git

# In your project's go.mod, add:
replace github.com/phatnt199/go-infra => /path/to/go-infra
```

## Verify Installation

Create a simple test file:

```go title="main.go"
package main

import (
    defaultlogger "github.com/phatnt199/go-infra/pkg/logger/default_logger"
)

func main() {
    log := defaultlogger.GetLogger()
    log.Info("go-infra installed successfully!")
}
```

Run it:

```bash
go run main.go
```

You should see a log message indicating successful installation.

## Database Setup (Optional)

If you plan to use database features:

### PostgreSQL with Docker

```bash
docker run --name postgres-dev \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=myapp \
  -p 5432:5432 \
  -d postgres:15-alpine
```

### PostgreSQL Installation

**macOS:**

```bash
brew install postgresql@15
brew services start postgresql@15
```

**Ubuntu/Debian:**

```bash
sudo apt-get update
sudo apt-get install postgresql postgresql-contrib
sudo systemctl start postgresql
```

**Windows:**

Download and install from [PostgreSQL Official Site](https://www.postgresql.org/download/windows/)

## Configuration Setup

go-infra uses **Viper** to load environment-aware configuration files.

### Step 1: Create `.env` File

```bash title=".env"
# Application
APP_ENV=development
APP_NAME=myapp

# Note: HTTP/DB config is in config.development.json
```

Copy to `.env.example` for version control:

```bash
cp .env .env.example
```

:::tip
Add `.env` to your `.gitignore` to avoid committing secrets:

```bash
echo ".env" >> .gitignore
```

:::

### Step 2: Create Config File

go-infra loads `config.<env>.json` based on your `APP_ENV`:

```json title="config/config.development.json"
{
	"fiberHttpOptions": {
		"port": ":8080",
		"host": "localhost",
		"basePath": "/api/v1",
		"name": "My App",
		"development": true,
		"timeout": 30
	},
	"gormOptions": {
		"type": 0,
		"host": "localhost",
		"port": 5432,
		"user": "postgres",
		"password": "postgres",
		"dbname": "myapp",
		"sslmode": false
	},
	"logOptions": {
		"level": "debug",
		"logType": 0,
		"callerEnabled": true
	}
}
```

:::important Required Configuration
The `config.<env>.json` file is **required** for most go-infra modules to work:

- `customfiber.Module` requires `fiberHttpOptions`
- `postgresgorm.Module` requires `gormOptions`
- Logger uses `logOptions`

Without this file, the application will fail to start.
:::

### Create Production Config

```json title="config/config.production.json"
{
	"fiberHttpOptions": {
		"port": ":8080",
		"host": "0.0.0.0",
		"basePath": "/api/v1",
		"name": "My App",
		"development": false,
		"timeout": 60
	},
	"gormOptions": {
		"type": 0,
		"host": "db.production.com",
		"port": 5432,
		"user": "prod_user",
		"password": "secret-password",
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

### Environment Variable Overrides

go-infra uses struct tags to map environment variables. You can override config file values by setting environment variables that match the struct field tags.

**Common environment variables:**

```bash
# HTTP Server
export TcpPort=:8080
export Host=0.0.0.0
export BasePath=/api/v1

# Database (GORM)
export DB_HOST=db.production.com
export DB_PORT=5432
export DB_USER=prod_user
export DB_PASSWORD=secret-password
export DB_NAME=myapp_prod
export SslMode=true

# Logger
export LOG_LEVEL=info
```

:::important
Environment variable names are **case-sensitive** and must match the struct tag names exactly. See your module's config struct for the exact tag names (e.g., `env:"TcpPort"` in `FiberHttpOptions`).
:::

See [Configuration Management](../core-concepts/configuration) for complete details.

## Database Migrations (Optional)

go-infra has **built-in Goose support** for database migrations.

### Setup Goose Migration Module

```go title="cmd/api/main.go"
import (
    "github.com/phatnt199/go-infra/pkg/migration/goose"
)

func main() {
    appBuilder := fxapp.NewApplicationBuilder()

    appBuilder.ProvideModule(customfiber.Module).
    appBuilder.ProvideModule(postgresgorm.Module).
    appBuilder.ProvideModule(goose.Module).  // Built-in Goose support

		app := appBuilder.Build()

    app.Run()
}
```

### Create Migration Files

```bash
# Create migrations directory
mkdir -p migrations

# Create migration file
cat > migrations/00001_create_users_table.sql << 'EOF'
-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users;
EOF
```

### Configure Migration Options

Add to your `config.development.json`:

```json
{
	"migrationOptions": {
		"host": "localhost",
		"port": 5432,
		"user": "postgres",
		"password": "postgres",
		"dbName": "myapp",
		"sslMode": false,
		"migrationsDir": "migrations",
		"skipMigration": false
	}
}
```

Migrations will run automatically on application start when `goose.Module` is included.

:::tip
Set `"skipMigration": true` to disable automatic migrations in production and run them manually.
:::

## Additional Tools

### Air (Hot Reload for Development)

```bash
go install github.com/cosmtrek/air@latest
```

Create `.air.toml`:

```toml title=".air.toml"
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./cmd/api"
bin = "tmp/main"
include_ext = ["go", "json"]
exclude_dir = ["tmp", "vendor"]
```

Run with hot reload:

```bash
air
```

### Swagger (API Documentation)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

See [API Documentation](../http-server/swagger) for setup.

## Verify Complete Setup

Create a complete test application:

```go title="cmd/api/main.go"
package main

import (
    "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
    customfiber "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
    postgresgorm "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
)

func main() {
    // Create application builder (includes logger)
    appBuilder := fxapp.NewApplicationBuilder()
    log := appBuilder.Logger()

    log.Info("Starting application...")

    // Build app with modules
    appBuilder.ProvideModule(customfiber.Module)
    appBuilder.ProvideModule(postgresgorm.Module)
    app := appBuilder.Build()

    app.Run()
}
```

Run it:

```bash
go mod tidy
go run cmd/api/main.go
```

You should see:

```
[INFO] Starting application...
[INFO] My App is listening on Host:localhost Http PORT: :8080
[INFO] Database connected
```

Visit `http://localhost:8080/health` to verify it's running.

## Troubleshooting

### "Config file not found"

**Error:** `viper.ReadInConfig: No directory with config file found`

**Solution:**

1. Ensure `config/config.development.json` exists
2. Set `APP_ENV=development` in `.env`
3. Or set `CONFIG_PATH` explicitly:

```bash
export CONFIG_PATH=$(pwd)/config
go run cmd/api/main.go
```

### "Cannot find package"

```bash
go mod download
go mod tidy
```

### Database Connection Fails

Verify PostgreSQL is running:

```bash
# macOS/Linux
psql -U postgres -h localhost

# Docker
docker ps | grep postgres
```

Check your `gormOptions` in `config.development.json` matches your database credentials.

### Port Already in Use

Change the port in `config.development.json`:

```json
{
  "fiberHttpOptions": {
    "port": ":8081",
    ...
  }
}
```

### Wrong Module Import Name

**Error:** `undefined: fiber_adapter`

**Solution:** Use `customfiber` instead:

```go
import customfiber "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"

// Then use:
appBuilder.ProvideModule(customfiber.Module)
```

The package is named `customfiber` to avoid conflicts with the Fiber web framework.

## Understanding go-infra Modules

go-infra uses **Uber Fx** for dependency injection. Modules are self-contained units:

```go
// HTTP Server Module
customfiber.Module       // Provides Fiber HTTP server

// Database Module
postgresgorm.Module      // Provides PostgreSQL with GORM

// Migration Module
goose.Module            // Provides Goose migrations (built-in)

// Health Module
health.Module           // Provides /health endpoint

// Your Custom Module
modules.Module          // Your app's business logic
```

Each module:

- Loads its own configuration via `config.BindConfigKey`
- Provides dependencies to the DI container
- Registers lifecycle hooks (start/stop)

See [Dependency Injection](../core-concepts/dependency-injection) for details.

## Configuration Discovery

go-infra automatically discovers configuration:

1. **Loads `.env`** - Recursively searches from current directory upward
2. **Reads `APP_ENV`** - Determines which config file to load
3. **Finds `config.<env>.json`** - Searches project subdirectories from `APP_ROOT_PATH`
4. **Loads with Viper** - Unmarshals into typed structs
5. **Applies overrides** - Environment variables override file values

**Recommended for production:**

```bash
# Set explicit paths to avoid filesystem search
export CONFIG_PATH=/app/config
export APP_ENV=production
./myapp
```

See [Configuration Management](../core-concepts/configuration) for complete details.

## Next Steps

Now that you have go-infra installed and configured:

- **[Quick Start](./quick-start)** - Build your first API in 10 minutes
- **[Configuration Guide](../core-concepts/configuration)** - Deep dive into config system
- **[Architecture](../core-concepts/architecture)** - Understand the framework design
- **[Examples](../examples/users-api)** - See real-world implementations

Happy coding! 🚀
