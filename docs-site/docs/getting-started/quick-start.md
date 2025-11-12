---
sidebar_position: 2
---

# Quick Start

Build your first REST API with go-infra in **5 minutes**.

## Step 1: Create Project

```bash
mkdir myapp && cd myapp
go mod init myapp
```

## Step 2: Install go-infra

```bash
go get github.com/phatnt199/go-infra
```

## Step 3: Setup Database

Start PostgreSQL with Docker:

```bash
docker run --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=myapp \
  -p 5432:5432 \
  -d postgres:15-alpine
```

## Step 4: Create Configuration

Create `.env`:

```bash title=".env"
APP_ENV=development
HTTP_PORT=8080
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=myapp
```

## Step 5: Define Your Model

```go title="internal/domain/user.go"
package domain

import "github.com/phatnt199/go-infra/pkg/domain/entity"

type User struct {
    entity.BaseModel
    Name  string `json:"name" gorm:"not null"`
    Email string `json:"email" gorm:"uniqueIndex;not null"`
}
```

## Step 6: Create Main Application

```go title="cmd/api/main.go"
package main

import (
    "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
    "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
    "myapp/internal/domain"
    "go.uber.org/fx"
)

func main() {
    app := fxapp.NewApplicationBuilder().
        ProvideModule(fiber_adapter.Module).
        ProvideModule(gorm.Module).
        Provide(fx.Invoke(autoMigrate)).
        Build()

    app.Run()
}

func autoMigrate(db *gorm.DB) {
    db.AutoMigrate(&domain.User{})
}
```

## Step 7: Run Your App

```bash
go run cmd/api/main.go
```

You should see:

```
[INFO] Server listening on :8080
[INFO] Database connected
```

## Step 8: Test It

Check health endpoint:

```bash
curl http://localhost:8080/health
```

## What You Just Built

In 5 minutes, you created an application with:

✅ **HTTP Server** - Fiber web framework  
✅ **Database** - PostgreSQL with GORM  
✅ **Configuration** - Environment-based config  
✅ **Logging** - Structured logging with context  
✅ **Health Checks** - Built-in health endpoint  
✅ **Graceful Shutdown** - Proper cleanup on exit  
✅ **Dependency Injection** - Clean architecture with Fx

## Next Steps

### Add REST Endpoints

Want to build REST endpoints? See [Building APIs](../http-server/building-apis).

### Add Authentication

Need user authentication? See [Authentication](../authentication/getting-started).

### Add More Features

- [Custom Routes](../http-server/routing)
- [Middleware](../http-server/middleware)
- [Database Migrations](../database/migrations)
- [Advanced Examples](../examples/users-api)

## Common Next Actions

### Add a Custom Handler

```go title="internal/handler/user_handler.go"
package handler

import (
    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
    "myapp/internal/domain"
)

type UserHandler struct {
    db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
    return &UserHandler{db: db}
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
    var users []domain.User
    h.db.Find(&users)
    return c.JSON(users)
}
```

### Register Routes

```go title="cmd/api/main.go"
func setupRoutes(
    server contracts.HttpServer,
    handler *handler.UserHandler,
) {
    server.RouteBuilder().RegisterHandler(func(router interface{}) {
        if app, ok := router.(*fiber.App); ok {
            api := app.Group("/api/v1")
            api.Get("/users", handler.GetUsers)
        }
    })
}
```

Update `main.go`:

```go
app := fxapp.NewApplicationBuilder().
    ProvideModule(fiber_adapter.Module).
    ProvideModule(gorm.Module).
    Provide(handler.NewUserHandler).
    Provide(fx.Invoke(setupRoutes)).
    Build()
```

## Troubleshooting

**Database connection fails?**

```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check connection
psql -U postgres -h localhost
```

**Port already in use?**

```bash
# Change port in .env
HTTP_PORT=8081
```

**Module errors?**

```bash
go mod tidy
go mod download
```

## What's Next?

You now have a working application! Here's what to explore:

1. **[Architecture Overview](../core-concepts/architecture)** - Understand the framework
2. **[HTTP Server Guide](../http-server/building-apis)** - Build REST APIs
3. **[Database Guide](../database/setup)** - Work with databases
4. **[Real Examples](../examples/users-api)** - See complete applications

Happy coding! 🚀
