---
sidebar_position: 1
---

# Building REST APIs

Learn how to build REST APIs with go-infra's HTTP server using Fiber.

## Basic Setup

### 1. Add HTTP Module

```go
import (
    "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
)

func main() {
    app := fxapp.NewApplicationBuilder().
        ProvideModule(fiber_adapter.Module).
        Build()

    app.Run()
}
```

### 2. Configuration

Set HTTP server options in `.env`:

```bash
HTTP_PORT=8080
HTTP_HOST=0.0.0.0
HTTP_READ_TIMEOUT=30s
HTTP_WRITE_TIMEOUT=30s
HTTP_IDLE_TIMEOUT=120s
```

Or in code:

```go
config := &fiber_adapter.Config{
    Port:         8080,
    Host:         "0.0.0.0",
    ReadTimeout:  30 * time.Second,
    WriteTimeout: 30 * time.Second,
}
```

## Creating Handlers

### Basic Handler

```go
package handler

import (
    "github.com/gofiber/fiber/v2"
    "github.com/phatnt199/go-infra/pkg/logger"
)

type UserHandler struct {
    logger logger.Logger
}

func NewUserHandler(logger logger.Logger) *UserHandler {
    return &UserHandler{logger: logger}
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
    h.logger.Info("Fetching users")

    users := []map[string]string{
        {"id": "1", "name": "John"},
        {"id": "2", "name": "Jane"},
    }

    return c.JSON(users)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    type CreateUserRequest struct {
        Name  string `json:"name" validate:"required"`
        Email string `json:"email" validate:"required,email"`
    }

    var req CreateUserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    h.logger.Info("Creating user", logger.String("email", req.Email))

    return c.Status(201).JSON(fiber.Map{
        "message": "User created",
        "user":    req,
    })
}
```

## Registering Routes

### Method 1: Using RegisterHandler

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
    "myapp/internal/handler"
    "go.uber.org/fx"
)

func main() {
    app := fxapp.NewApplicationBuilder().
        ProvideModule(fiber_adapter.Module).
        Provide(handler.NewUserHandler).
        Provide(fx.Invoke(setupRoutes)).
        Build()

    app.Run()
}

func setupRoutes(
    server contracts.HttpServer,
    userHandler *handler.UserHandler,
) {
    server.RouteBuilder().RegisterHandler(func(router interface{}) {
        app := router.(*fiber.App)

        api := app.Group("/api/v1")
        api.Get("/users", userHandler.GetUsers)
        api.Post("/users", userHandler.CreateUser)
    })
}
```

### Method 2: Route Groups

```go
func setupRoutes(
    server contracts.HttpServer,
    userHandler *handler.UserHandler,
    productHandler *handler.ProductHandler,
) {
    server.RouteBuilder().RegisterHandler(func(router interface{}) {
        app := router.(*fiber.App)

        // API v1
        v1 := app.Group("/api/v1")

        // User routes
        users := v1.Group("/users")
        users.Get("/", userHandler.GetUsers)
        users.Post("/", userHandler.CreateUser)
        users.Get("/:id", userHandler.GetUser)
        users.Put("/:id", userHandler.UpdateUser)
        users.Delete("/:id", userHandler.DeleteUser)

        // Product routes
        products := v1.Group("/products")
        products.Get("/", productHandler.GetProducts)
        products.Post("/", productHandler.CreateProduct)
    })
}
```

## Request Handling

### Path Parameters

```go
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
    id := c.Params("id")

    user, err := h.service.GetUserByID(id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "error": "User not found",
        })
    }

    return c.JSON(user)
}
```

### Query Parameters

```go
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
    // /users?page=1&limit=10
    page := c.QueryInt("page", 1)
    limit := c.QueryInt("limit", 10)
    search := c.Query("search", "")

    users, err := h.service.GetUsers(page, limit, search)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to fetch users",
        })
    }

    return c.JSON(fiber.Map{
        "data":  users,
        "page":  page,
        "limit": limit,
    })
}
```

### Request Body

```go
type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=2"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"min=0,max=150"`
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req CreateUserRequest

    // Parse JSON body
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // Validate request
    if err := validator.Validate(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error":   "Validation failed",
            "details": err.Error(),
        })
    }

    // Process request
    user, err := h.service.CreateUser(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(201).JSON(user)
}
```

### Headers

```go
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
    // Read headers
    authToken := c.Get("Authorization")
    userAgent := c.Get("User-Agent")

    // Set response headers
    c.Set("X-Custom-Header", "value")

    return c.JSON(user)
}
```

## Response Handling

### JSON Response

```go
// Simple JSON
return c.JSON(fiber.Map{
    "message": "Success",
    "data":    data,
})

// Struct response
return c.JSON(user)

// With status code
return c.Status(201).JSON(user)
```

### Error Responses

```go
// Standard error
return c.Status(400).JSON(fiber.Map{
    "error": "Invalid input",
})

// Detailed error
return c.Status(422).JSON(fiber.Map{
    "error": "Validation failed",
    "details": []string{
        "Email is required",
        "Name must be at least 2 characters",
    },
})
```

### Success Responses

```go
// 200 OK
return c.JSON(data)

// 201 Created
return c.Status(201).JSON(data)

// 204 No Content
return c.SendStatus(204)
```

## Middleware

### Global Middleware

```go
func setupRoutes(server contracts.HttpServer) {
    server.RouteBuilder().RegisterHandler(func(router interface{}) {
        app := router.(*fiber.App)

        // Apply to all routes
        app.Use(LoggingMiddleware())
        app.Use(RecoveryMiddleware())
        app.Use(CORSMiddleware())

        // Routes
        app.Get("/users", handler.GetUsers)
    })
}
```

### Route-Specific Middleware

```go
// Apply to specific routes
api := app.Group("/api/v1")
api.Use(AuthMiddleware())

// Public routes (no middleware)
app.Get("/health", handler.Health)

// Protected routes (with middleware)
api.Get("/users", handler.GetUsers)
```

### Custom Middleware

```go
func LoggingMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()

        // Before request
        logger.Info("Request started",
            logger.String("method", c.Method()),
            logger.String("path", c.Path()),
        )

        // Process request
        err := c.Next()

        // After request
        duration := time.Since(start)
        logger.Info("Request completed",
            logger.Duration("duration", duration),
            logger.Int("status", c.Response().StatusCode()),
        )

        return err
    }
}
```

## Validation

### Request Validation

```go
import "github.com/phatnt199/go-infra/pkg/validator"

type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"min=18,max=120"`
    Password string `json:"password" validate:"required,min=8"`
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req CreateUserRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request",
        })
    }

    // Validate
    if err := validator.Validate(&req); err != nil {
        return c.Status(422).JSON(fiber.Map{
            "error":   "Validation failed",
            "details": err,
        })
    }

    // Process...
    return c.JSON(result)
}
```

### Custom Validation

```go
func ValidateCreateUser(req CreateUserRequest) error {
    if req.Age < 18 {
        return errors.New("must be 18 or older")
    }

    if !strings.Contains(req.Email, "@") {
        return errors.New("invalid email format")
    }

    return nil
}
```

## Error Handling

### Centralized Error Handler

```go
func ErrorHandler() fiber.ErrorHandler {
    return func(c *fiber.Ctx, err error) error {
        code := fiber.StatusInternalServerError
        message := "Internal server error"

        // Check error type
        if e, ok := err.(*fiber.Error); ok {
            code = e.Code
            message = e.Message
        }

        // Log error
        logger.Error("Request error",
            logger.Err(err),
            logger.String("path", c.Path()),
            logger.Int("status", code),
        )

        // Return error response
        return c.Status(code).JSON(fiber.Map{
            "error": message,
        })
    }
}

// Register error handler
app := fiber.New(fiber.Config{
    ErrorHandler: ErrorHandler(),
})
```

### Error Utilities

```go
// Not found error
if user == nil {
    return fiber.NewError(404, "User not found")
}

// Bad request error
if invalid {
    return fiber.NewError(400, "Invalid input")
}

// Unauthorized error
if !authenticated {
    return fiber.NewError(401, "Unauthorized")
}
```

## File Upload

```go
func (h *Handler) UploadFile(c *fiber.Ctx) error {
    // Get file from request
    file, err := c.FormFile("file")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "No file uploaded",
        })
    }

    // Validate file
    if file.Size > 10*1024*1024 { // 10MB
        return c.Status(400).JSON(fiber.Map{
            "error": "File too large",
        })
    }

    // Save file
    filename := fmt.Sprintf("./uploads/%s", file.Filename)
    if err := c.SaveFile(file, filename); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to save file",
        })
    }

    return c.JSON(fiber.Map{
        "message":  "File uploaded",
        "filename": file.Filename,
        "size":     file.Size,
    })
}
```

## Static Files

```go
func setupRoutes(server contracts.HttpServer) {
    server.RouteBuilder().RegisterHandler(func(router interface{}) {
        app := router.(*fiber.App)

        // Serve static files
        app.Static("/", "./public")
        app.Static("/uploads", "./uploads")

        // With custom config
        app.Static("/assets", "./assets", fiber.Static{
            Compress:  true,
            ByteRange: true,
            MaxAge:    3600,
        })
    })
}
```

## CORS Configuration

```go
import "github.com/gofiber/fiber/v2/middleware/cors"

func setupRoutes(server contracts.HttpServer) {
    server.RouteBuilder().RegisterHandler(func(router interface{}) {
        app := router.(*fiber.App)

        app.Use(cors.New(cors.Config{
            AllowOrigins: "https://example.com,https://app.example.com",
            AllowMethods: "GET,POST,PUT,DELETE",
            AllowHeaders: "Origin,Content-Type,Authorization",
            AllowCredentials: true,
            MaxAge: 3600,
        }))
    })
}
```

## Rate Limiting

```go
import "github.com/gofiber/fiber/v2/middleware/limiter"

func setupRoutes(server contracts.HttpServer) {
    server.RouteBuilder().RegisterHandler(func(router interface{}) {
        app := router.(*fiber.App)

        app.Use(limiter.New(limiter.Config{
            Max:        100,
            Expiration: 1 * time.Minute,
            LimitReached: func(c *fiber.Ctx) error {
                return c.Status(429).JSON(fiber.Map{
                    "error": "Too many requests",
                })
            },
        }))
    })
}
```

## Complete Example

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
    "github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
    "github.com/phatnt199/go-infra/pkg/logger"
    "go.uber.org/fx"
    "myapp/internal/handler"
)

func main() {
    app := fxapp.NewApplicationBuilder().
        ProvideModule(fiber_adapter.Module).
        Provide(handler.NewUserHandler).
        Provide(fx.Invoke(setupRoutes)).
        Build()

    app.Run()
}

func setupRoutes(
    server contracts.HttpServer,
    userHandler *handler.UserHandler,
    log logger.Logger,
) {
    server.RouteBuilder().RegisterHandler(func(router interface{}) {
        app := router.(*fiber.App)

        // Global middleware
        app.Use(LoggingMiddleware(log))
        app.Use(RecoveryMiddleware())

        // Public routes
        app.Get("/health", func(c *fiber.Ctx) error {
            return c.JSON(fiber.Map{"status": "ok"})
        })

        // API routes
        api := app.Group("/api/v1")

        // User routes
        users := api.Group("/users")
        users.Get("/", userHandler.GetUsers)
        users.Post("/", userHandler.CreateUser)
        users.Get("/:id", userHandler.GetUser)
        users.Put("/:id", userHandler.UpdateUser)
        users.Delete("/:id", userHandler.DeleteUser)

        log.Info("Routes registered")
    })
}

func LoggingMiddleware(log logger.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()

        err := c.Next()

        log.Info("Request",
            logger.String("method", c.Method()),
            logger.String("path", c.Path()),
            logger.Int("status", c.Response().StatusCode()),
            logger.Duration("duration", time.Since(start)),
        )

        return err
    }
}

func RecoveryMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        defer func() {
            if r := recover(); r != nil {
                logger.Error("Panic recovered",
                    logger.Any("panic", r),
                )
                c.Status(500).JSON(fiber.Map{
                    "error": "Internal server error",
                })
            }
        }()
        return c.Next()
    }
}
```

## Next Steps

- **[CRUD Operations](./crud-operations)** - Auto-generate CRUD endpoints
- **[Middleware](./middleware)** - Deep dive into middleware
- **[Authentication](../authentication/getting-started)** - Add auth to your API
- **[Swagger](./swagger)** - Generate API documentation
