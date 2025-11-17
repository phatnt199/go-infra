---
sidebar_position: 2
---

# Quick Start

Build your first REST API with go-infra in **10 minutes**. This guide creates a complete, production-ready API with database, configuration, and health checks.

## Prerequisites

- Go 1.21+ installed
- PostgreSQL running (Docker or local)
- Basic Go knowledge

## Step 1: Create Project Structure

```bash
mkdir myapp && cd myapp
go mod init myapp

# Create directory structure
mkdir -p cmd/api internal/{config,domain,handler,modules,repository} config
```

Your project structure:

```
myapp/
├── cmd/
│   └── api/
│       └── main.go
├── config/
│   └── config.development.json
├── internal/
│   ├── config/
│   ├── domain/
│   ├── handler/
│   ├── modules/
│   └── repository/
├── .env
└── go.mod
```

## Step 2: Install go-infra

```bash
go get github.com/phatnt199/go-infra
```

## Step 3: Start PostgreSQL

Using Docker (recommended):

```bash
docker run --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=myapp \
  -p 5432:5432 \
  -d postgres:15-alpine
```

Or use your local PostgreSQL installation.

## Step 4: Create Configuration

### Create `.env` file

```bash title=".env"
APP_ENV=development
APP_NAME=myapp
```

:::important
The `.env` file must be in your project root. go-infra automatically loads it recursively from the current directory upward.
:::

### Create `config.development.json`

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

:::tip
go-infra uses Viper to load `config.<env>.json` files based on `APP_ENV`. The framework automatically searches for this file in your project subdirectories.
:::

## Step 5: Define Your Domain Model

```go title="internal/domain/user.go"
package domain

import "github.com/phatnt199/go-infra/pkg/domain/entity"

type User struct {
    entity.BaseModel
    Name  string `json:"name" gorm:"not null"`
    Email string `json:"email" gorm:"uniqueIndex;not null"`
}

func (User) TableName() string {
    return "users"
}
```

## Step 6: Create Repository

```go title="internal/repository/user_repository.go"
package repository

import (
    "myapp/internal/domain"
    "gorm.io/gorm"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindAll() ([]domain.User, error) {
    var users []domain.User
    err := r.db.Find(&users).Error
    return users, err
}

func (r *UserRepository) FindByID(id string) (*domain.User, error) {
    var user domain.User
    err := r.db.Where("id = ?", id).First(&user).Error
    return &user, err
}
```

## Step 7: Create HTTP Handler

```go title="internal/handler/user_handler.go"
package handler

import (
    "myapp/internal/domain"
    "myapp/internal/repository"
    "github.com/gofiber/fiber/v2"
    "github.com/phatnt199/go-infra/pkg/logger"
)

type UserHandler struct {
    repo   *repository.UserRepository
    logger logger.Logger
}

func NewUserHandler(repo *repository.UserRepository, logger logger.Logger) *UserHandler {
    return &UserHandler{
        repo:   repo,
        logger: logger,
    }
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var user domain.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    if err := h.repo.Create(&user); err != nil {
        h.logger.Errorf("Failed to create user: %v", err)
        return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
    }

    return c.Status(201).JSON(user)
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
    users, err := h.repo.FindAll()
    if err != nil {
        h.logger.Errorf("Failed to get users: %v", err)
        return c.Status(500).JSON(fiber.Map{"error": "Failed to get users"})
    }

    return c.JSON(users)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
    id := c.Params("id")
    user, err := h.repo.FindByID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "User not found"})
    }

    return c.JSON(user)
}
```

## Step 8: Setup Routes Module

```go title="internal/modules/module.go"
package modules

import (
    "myapp/internal/handler"
    "myapp/internal/repository"
    "github.com/gofiber/fiber/v2"
    "github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
    "github.com/phatnt199/go-infra/pkg/logger"
    "go.uber.org/fx"
)

// Module provides all dependencies using fx
var Module = fx.Module(
    "app_module",
    fx.Provide(repository.NewUserRepository),
    fx.Provide(handler.NewUserHandler),
    fx.Invoke(setupRoutes),
)

func setupRoutes(
    server contracts.HttpServer,
    userHandler *handler.UserHandler,
    logger logger.Logger,
) {
    server.RouteBuilder().RegisterHandler(func(instance interface{}) {
        if app, ok := instance.(*fiber.App); ok {
            api := app.Group("/api/v1")

            // User routes
            users := api.Group("/users")
            users.Get("/", userHandler.GetUsers)
            users.Get("/:id", userHandler.GetUser)
            users.Post("/", userHandler.CreateUser)

            log.Info("Routes registered successfully")
        }
    })
}
```

## Step 9: Create Main Application

```go title="cmd/api/main.go"
package main

import (
    "myapp/internal/domain"
    "myapp/internal/modules"

    "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
    customfiber "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
    postgresgorm "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
    "github.com/phatnt199/go-infra/pkg/logger"

    "go.uber.org/fx"
    "gorm.io/gorm"
)

func main() {
    // Create application builder
    appBuilder := fxapp.NewApplicationBuilder()

    // Register go-infra modules
    appBuilder.ProvideModule(customfiber.Module)   // HTTP server
    appBuilder.ProvideModule(postgresgorm.Module)  // PostgreSQL

    // Register your app module
    appBuilder.ProvideModule(modules.Module)

    // Auto-migrate database
    appBuilder.Provide(fx.Invoke(runMigrations))

    // Build and run
    app := appBuilder.Build()
    app.Run()
}

func runMigrations(db *gorm.DB, log logger.Logger) {
    log.Info("Running database migrations...")

    if err := db.AutoMigrate(&domain.User{}); err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }

    log.Info("Database migrations completed")
}
```

:::tip Module Import Name
Notice we import the fiber adapter as `customfiber`. The package name is `customfiber`, not `fiber_adapter`. This is go-infra's naming convention to avoid conflicts with the standard `fiber` package.
:::

## Step 10: Run Your Application

From the project root:

```bash
go mod tidy
go run cmd/api/main.go
```

You should see output like:

```
[INFO] Running database migrations...
[INFO] Database migrations completed
[INFO] My API is listening on Host:localhost Http PORT: :8080
[INFO] Routes registered successfully
```

## Step 11: Test Your API

### Health check (built-in):

```bash
curl http://localhost:8080/health
```

Response:

```json
{
	"status": "ok",
	"timestamp": "2024-01-15T10:30:00Z"
}
```

### Create a user:

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

Response:

```json
{
	"id": "uuid-here",
	"name": "John Doe",
	"email": "john@example.com",
	"created_at": "2024-01-15T10:30:00Z",
	"updated_at": "2024-01-15T10:30:00Z"
}
```

### Get all users:

```bash
curl http://localhost:8080/api/v1/users
```

### Get specific user:

```bash
curl http://localhost:8080/api/v1/users/{user-id}
```

## What You Just Built

Congratulations! In 10 minutes, you created a production-ready application with:

✅ **HTTP Server** - Fiber web framework with routing  
✅ **Database** - PostgreSQL with GORM ORM  
✅ **Configuration** - Viper-powered config with `.env` support  
✅ **Logging** - Structured logging with context  
✅ **Health Checks** - Built-in `/health` endpoint  
✅ **Auto-Migration** - Automatic database schema updates  
✅ **Graceful Shutdown** - Proper cleanup on exit  
✅ **Dependency Injection** - Clean architecture with Uber Fx

## Project Structure Explained

```
myapp/
├── .env                          # Environment variables (APP_ENV, APP_NAME)
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── config/
│   └── config.development.json  # Environment-specific config
├── internal/
│   ├── domain/                  # Domain models
│   ├── handler/                 # HTTP handlers (controllers)
│   ├── modules/                 # Fx modules with DI wiring
│   └── repository/              # Data access layer
└── go.mod
```

## Common Next Steps

### Add More Routes

Expand your `setupRoutes` function:

```go
users.Put("/:id", userHandler.UpdateUser)
users.Delete("/:id", userHandler.DeleteUser)
```

### Add Validation

Install validator:

```bash
go get github.com/go-playground/validator/v10
```

Add to your domain:

```go
type User struct {
    entity.BaseModel
    Name  string `json:"name" gorm:"not null" validate:"required,min=2"`
    Email string `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
}
```

### Add Authentication

See the [Authentication Guide](../authentication/getting-started) for JWT authentication setup.

### Use Database Migrations

Instead of AutoMigrate, use Goose for version-controlled migrations:

```go
import "github.com/phatnt199/go-infra/pkg/migration/goose"

// In main.go
appBuilder.ProvideModule(goose.Module)
```

See [Database Migrations](../database/migrations) for details.

## Troubleshooting

### "Config file not found"

**Error:** `viper.ReadInConfig: No directory with config file found`

**Solution:** Ensure `config/config.development.json` exists and `APP_ENV=development` is set in `.env`.

You can also set `CONFIG_PATH`:

```bash
export CONFIG_PATH=$(pwd)/config
go run cmd/api/main.go
```

### "Cannot connect to database"

**Solution:** Verify PostgreSQL is running:

```bash
docker ps | grep postgres
# or
psql -U postgres -h localhost -c "SELECT 1"
```

Check your `gormOptions` in `config.development.json` matches your database credentials.

### Port already in use

**Solution:** Change the port in `config.development.json`:

```json
{
  "fiberHttpOptions": {
    "port": ":8081",
    ...
  }
}
```

### Module not found

```bash
go mod tidy
go get github.com/phatnt199/go-infra
```

## Understanding the Configuration

go-infra uses a two-level configuration system:

1. **`.env` file** - Sets `APP_ENV` (development/staging/production)
2. **`config.<env>.json`** - Environment-specific settings loaded via Viper

The framework:

- Automatically loads `.env` from project root
- Searches for `config.development.json` (or staging/production)
- Allows environment variable overrides

See [Configuration Management](../core-concepts/configuration) for complete details.

## Next Steps

Now that you have a working API, explore more features:

### Build More Features

- **[Building APIs](../http-server/building-apis)** - Advanced routing and middleware
- **[Database Guide](../database/setup)** - Relationships, transactions, advanced queries
- **[Authentication](../authentication/getting-started)** - Add JWT authentication
- **[Validation](../validation/setup)** - Request validation

### Explore Examples

- **[Users API Example](../examples/users-api)** - Complete CRUD API
- **[Authentication Service](../examples/authentication-service)** - Full auth implementation
- **[Microservices](../examples/microservices)** - Multi-service architecture

### Learn Architecture

- **[Architecture Overview](../core-concepts/architecture)** - Understand the framework design
- **[Dependency Injection](../core-concepts/dependency-injection)** - Learn Fx patterns
- **[CQRS Pattern](../core-concepts/cqrs)** - Command/Query separation

Ready to build something amazing? 🚀
