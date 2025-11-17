---
sidebar_position: 1
---

# Users API Example

A complete example demonstrating a REST API with CRUD operations, database integration, and Swagger documentation.

## Overview

This example shows how to build a production-ready REST API using go-infra. It includes:

- ✅ **REST API** with Fiber
- ✅ **PostgreSQL** database with GORM
- ✅ **CRUD operations** for users
- ✅ **Database migrations**
- ✅ **Swagger documentation**
- ✅ **Health checks**
- ✅ **Dependency injection** with Fx

## Project Structure

```
users-api/
├── cmd/
│   ├── api/
│   │   └── main.go          # HTTP server entry point
│   └── migration/
│       └── main.go          # Migration runner
├── internal/
│   ├── domain/
│   │   └── user.go          # User model
│   ├── repository/
│   │   └── user_repository.go
│   ├── handler/
│   │   └── user_handler.go  # HTTP handlers
│   └── modules/
│       └── module.go        # Module wiring
├── migrations/              # SQL migration files
├── docs/                    # Generated Swagger docs
├── .env                     # Configuration
├── Makefile                 # Build commands
└── go.mod
```

## Quick Start

### 1. Setup Database

```bash
docker run --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=usersdb \
  -p 5432:5432 \
  -d postgres:15-alpine
```

### 2. Configure Environment

Create `.env`:

```bash
HTTP_PORT=8080
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=usersdb
```

### 3. Run Migrations

```bash
cd examples/users-api
make migrate
```

### 4. Start Server

```bash
make run
```

Server starts at `http://localhost:8080`

## API Endpoints

### Get All Users

```bash
curl http://localhost:8080/api/v1/users
```

Response:

```json
{
	"data": [
		{
			"id": "123e4567-e89b-12d3-a456-426614174000",
			"name": "John Doe",
			"email": "john@example.com",
			"created_at": "2024-01-01T00:00:00Z"
		}
	],
	"total": 1
}
```

### Get User by ID

```bash
curl http://localhost:8080/api/v1/users/{id}
```

### Create User

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com"
  }'
```

### Update User

```bash
curl -X PUT http://localhost:8080/api/v1/users/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com"
  }'
```

### Delete User

```bash
curl -X DELETE http://localhost:8080/api/v1/users/{id}
```

### Health Check

```bash
curl http://localhost:8080/api/v1/health
```

## Code Walkthrough

### 1. Define Model

```go title="internal/domain/user.go"
package domain

import "github.com/phatnt199/go-infra/pkg/domain/entity"

type User struct {
    entity.BaseModel
    Name  string `json:"name" gorm:"not null"`
    Email string `json:"email" gorm:"uniqueIndex;not null"`
}
```

### 2. Create Repository

```go title="internal/repository/user_repository.go"
package repository

import (
    "context"
    "gorm.io/gorm"
    "myapp/internal/domain"
)

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
    var users []domain.User
    err := r.db.WithContext(ctx).Find(&users).Error
    return users, err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    var user domain.User
    err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
    return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
    return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
    return r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error
}
```

### 3. Create Handler

```go title="internal/handler/user_handler.go"
package handler

import (
    "github.com/gofiber/fiber/v2"
    "myapp/internal/repository"
    "myapp/internal/domain"
)

type UserHandler struct {
    repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
    return &UserHandler{repo: repo}
}

// @Summary Get all users
// @Tags users
// @Produce json
// @Success 200 {array} domain.User
// @Router /users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
    users, err := h.repo.FindAll(c.Context())
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(users)
}

// @Summary Get user by ID
// @Tags users
// @Param id path string true "User ID"
// @Produce json
// @Success 200 {object} domain.User
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
    id := c.Params("id")
    user, err := h.repo.FindByID(c.Context(), id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "User not found"})
    }
    return c.JSON(user)
}

// @Summary Create user
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.User true "User"
// @Success 201 {object} domain.User
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var user domain.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    if err := h.repo.Create(c.Context(), &user); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(201).JSON(user)
}
```

### 4. Wire Everything

```go title="cmd/api/main.go"
package main

import (
    "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
    "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
    "myapp/internal/modules"
)

// @title Users API
// @version 1.0
// @description REST API for user management
// @host localhost:8080
// @BasePath /api/v1
func main() {
    appBuilder := fxapp.NewApplicationBuilder()
    appBuilder.ProvideModule(fiber_adapter.Module)
    appBuilder.ProvideModule(gorm.Module)
    appBuilder.ProvideModule(modules.Module)
    app := appBuilder.Build()

    app.Run()
}
```

## Swagger Documentation

Access Swagger UI at: `http://localhost:8080/swagger/index.html`

### Generate Swagger Docs

```bash
make swagger
```

Or manually:

```bash
swag init -g cmd/api/main.go -o docs
```

## Makefile Commands

```makefile
# Install dependencies
make deps

# Run migrations
make migrate

# Generate swagger docs
make swagger

# Run server
make run

# Start PostgreSQL
make docker-db
```

## Testing

### Manual Testing

```bash
# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test User", "email": "test@example.com"}'

# Get all users
curl http://localhost:8080/api/v1/users

# Get specific user
curl http://localhost:8080/api/v1/users/{id}

# Update user
curl -X PUT http://localhost:8080/api/v1/users/{id} \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Name", "email": "updated@example.com"}'

# Delete user
curl -X DELETE http://localhost:8080/api/v1/users/{id}
```

### Using the Test Script

```bash
./test-apis.sh
```

## Key Learnings

### 1. Module System

go-infra uses modules for dependency injection:

```go
appBuilder := fxapp.NewApplicationBuilder()
appBuilder.ProvideModule(fiber_adapter.Module)    // HTTP server
appBuilder.ProvideModule(gorm.Module)             // Database
appBuilder.ProvideModule(modules.Module)          // Your app module
app := appBuilder.Build()
```

### 2. Repository Pattern

Separates data access from business logic:

```go
type UserRepository struct {
    db *gorm.DB
}

func (r *UserRepository) FindAll(ctx context.Context) ([]domain.User, error) {
    var users []domain.User
    err := r.db.WithContext(ctx).Find(&users).Error
    return users, err
}
```

### 3. Dependency Injection

Dependencies are injected via constructors:

```go
func NewUserHandler(repo *repository.UserRepository) *UserHandler {
    return &UserHandler{repo: repo}
}
```

### 4. Context Propagation

Context flows through all layers:

```go
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
    users, err := h.repo.FindAll(c.Context())  // Pass context
    // ...
}
```

## Next Steps

- Add **authentication** - See [Authentication Example](./authentication-service)
- Add **validation** - Validate request data
- Add **pagination** - Handle large datasets
- Add **filtering** - Filter and search users
- Add **tests** - Write unit and integration tests

## Full Source Code

The complete source code is available at:

```
examples/users-api/
```

Run the example:

```bash
cd examples/users-api
make deps
make migrate
make run
```
