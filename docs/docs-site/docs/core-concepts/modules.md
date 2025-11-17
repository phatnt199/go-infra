---
sidebar_position: 2
---

# Module System

Learn how go-infra uses modules to organize and manage application components.

## Overview

go-infra uses a modular architecture where each feature is encapsulated in a module. This approach promotes:

- **Separation of concerns** - Each module handles a specific domain
- **Reusability** - Modules can be shared across projects
- **Testability** - Isolated modules are easier to test
- **Maintainability** - Changes are localized to specific modules

## Module Structure

A typical module in go-infra follows this structure:

```
pkg/component/authentication/
├── handler.go          # HTTP handlers
├── service.go          # Business logic
├── repository.go       # Data access
├── models.go           # Data structures
├── validation.go       # Input validation
└── middleware.go       # HTTP middleware
```

## Core Modules

### HTTP Adapter Module

The HTTP adapter module provides web server functionality:

```go
import (
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber"
    "github.com/phatnt199/go-infra/pkg/adapter/http/crud"
)

// Create HTTP server
app, router := fiber.New()

// Register CRUD routes
crud.RegisterCRUD[*User](router, db, &crud.CRUDOptions[*User, string]{
    BasePath: "/api/users",
})
```

### Authentication Module

Handles user authentication and authorization:

```go
import (
    "github.com/phatnt199/go-infra/pkg/component/authentication"
    "github.com/phatnt199/go-infra/pkg/crypto"
)

// Hash password
hashedPassword, _ := crypto.HashPassword("secret123")

// Generate JWT token
token, _ := crypto.GenerateJWT("user-id-123", 24*time.Hour)

// Verify token
claims, _ := crypto.VerifyJWT(token)
```

### Database Module

Provides database connectivity and ORM support:

```go
import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

// Initialize database
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// Run migrations
db.AutoMigrate(&User{}, &Product{})
```

### Logger Module

Structured logging with dependency injection:

```go
import (
    "github.com/phatnt199/go-infra/pkg/logger"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

// Logger is provided via Fx module
// In your component:
func MyService(log logger.Logger) {
    // Log with structured fields
    log.Infow("User created", logger.Fields{
        "user_id": user.ID,
        "email": user.Email,
    })
}
```

## Creating Custom Modules

### Step 1: Define Module Structure

```go
// pkg/modules/notification/types.go
package notification

type NotificationService interface {
    Send(recipient string, message string) error
}

type EmailNotification struct {
    smtpHost string
    smtpPort int
}
```

### Step 2: Implement Business Logic

```go
// pkg/modules/notification/service.go
package notification

func NewEmailNotification(host string, port int) *EmailNotification {
    return &EmailNotification{
        smtpHost: host,
        smtpPort: port,
    }
}

func (e *EmailNotification) Send(recipient, message string) error {
    // Implementation
    return nil
}
```

### Step 3: Register with HTTP Handler

```go
// pkg/modules/notification/handler.go
package notification

import "github.com/gofiber/fiber/v2"

type Handler struct {
    service NotificationService
}

func NewHandler(service NotificationService) *Handler {
    return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router fiber.Router) {
    router.Post("/notifications/send", h.Send)
}

func (h *Handler) Send(c *fiber.Ctx) error {
    var req SendRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    if err := h.service.Send(req.Recipient, req.Message); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"success": true})
}
```

### Step 4: Wire It Up

```go
// main.go
package main

import (
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber"
    "myapp/pkg/modules/notification"
)

func main() {
    app, router := fiber.New()

    // Create notification service
    notifService := notification.NewEmailNotification("smtp.gmail.com", 587)

    // Register handler
    handler := notification.NewHandler(notifService)
    handler.RegisterRoutes(router)

    app.Listen(":3000")
}
```

## Module Communication

### Direct Dependencies

Modules can depend on each other directly:

```go
type OrderService struct {
    userService   *UserService
    notifService  NotificationService
}

func (s *OrderService) CreateOrder(userID string, items []Item) error {
    user, err := s.userService.GetByID(userID)
    if err != nil {
        return err
    }

    // Create order...

    // Send notification
    s.notifService.Send(user.Email, "Order created successfully")
    return nil
}
```

### Event-Based Communication

For loosely coupled modules, use events:

```go
// Define event
type OrderCreatedEvent struct {
    OrderID   string
    UserID    string
    Total     float64
}

// Publish event
publisher.Publish("order.created", OrderCreatedEvent{
    OrderID: order.ID,
    UserID:  order.UserID,
    Total:   order.Total,
})

// Subscribe to event
subscriber.Subscribe("order.created", func(event OrderCreatedEvent) {
    // Handle event
    notifService.Send(event.UserID, "Your order is confirmed")
})
```

## Module Configuration

Each module can have its own configuration:

```go
// config/config.go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Email    EmailConfig
}

type EmailConfig struct {
    Host     string
    Port     int
    Username string
    Password string
}

// Load configuration
cfg, err := config.Load()
emailService := notification.NewEmailNotification(
    cfg.Email.Host,
    cfg.Email.Port,
)
```

## Best Practices

### 1. Single Responsibility

Each module should have one clear responsibility:

```go
// ✅ Good - focused responsibility
type UserAuthService struct {}  // Handles user authentication
type UserProfileService struct {} // Handles user profiles

// ❌ Bad - too many responsibilities
type UserService struct {}  // Handles auth, profiles, notifications, etc.
```

### 2. Dependency Injection

Inject dependencies rather than creating them:

```go
// ✅ Good
type OrderService struct {
    db     *gorm.DB
    logger logger.Logger
}

func NewOrderService(db *gorm.DB, logger logger.Logger) *OrderService {
    return &OrderService{db: db, logger: logger}
}

// ❌ Bad
type OrderService struct {}

func NewOrderService() *OrderService {
    db := gorm.Open(...)  // Creates its own dependencies
    return &OrderService{}
}
```

### 3. Interface-Based Design

Program to interfaces for flexibility:

```go
// Define interface
type PaymentProcessor interface {
    ProcessPayment(amount float64) error
}

// Multiple implementations
type StripeProcessor struct {}
type PayPalProcessor struct {}

// Service depends on interface
type OrderService struct {
    payment PaymentProcessor  // Can use any implementation
}
```

### 4. Error Handling

Handle errors consistently across modules:

```go
import "emperror.dev/errors"

func (s *UserService) CreateUser(user *User) error {
    if err := s.validate(user); err != nil {
        return errors.Wrap(err, "validation failed")
    }

    if err := s.db.Create(user).Error; err != nil {
        return errors.Wrap(err, "failed to create user")
    }

    return nil
}
```

## Testing Modules

### Unit Testing

Test modules in isolation:

```go
func TestUserService_CreateUser(t *testing.T) {
    // Setup
    db := setupTestDB()
    service := NewUserService(db)

    // Test
    user := &User{Name: "John", Email: "john@example.com"}
    err := service.CreateUser(user)

    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, user.ID)
}
```

### Integration Testing

Test modules together:

```go
func TestOrderFlow(t *testing.T) {
    // Setup all required services
    db := setupTestDB()
    userService := NewUserService(db)
    orderService := NewOrderService(db, userService)

    // Create user
    user := &User{Name: "John"}
    userService.CreateUser(user)

    // Create order
    order := &Order{UserID: user.ID}
    err := orderService.CreateOrder(order)

    assert.NoError(t, err)
}
```

## Next Steps

- Learn about [Dependency Injection](./dependency-injection)
- Explore [Configuration Management](./configuration)
- See [Complete Examples](../examples/users-api)
