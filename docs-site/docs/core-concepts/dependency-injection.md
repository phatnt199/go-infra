---
sidebar_position: 3
---

# Dependency Injection

Learn how to use dependency injection patterns in go-infra applications.

## Why Dependency Injection?

Dependency Injection (DI) is a design pattern that helps you:

- **Decouple components** - Reduce tight coupling between modules
- **Improve testability** - Easy to mock dependencies in tests
- **Enhance flexibility** - Swap implementations without changing code
- **Better organization** - Clear dependency hierarchy

## Constructor Injection

The most common pattern in Go - pass dependencies via constructor:

```go
// Define dependencies
type UserRepository interface {
    Create(user *User) error
    FindByID(id string) (*User, error)
}

type EmailService interface {
    Send(to, subject, body string) error
}

// Service with injected dependencies
type UserService struct {
    repo  UserRepository
    email EmailService
    logger logger.Logger
}

// Constructor injection
func NewUserService(
    repo UserRepository,
    email EmailService,
    logger logger.Logger,
) *UserService {
    return &UserService{
        repo:   repo,
        email:  email,
        logger: logger,
    }
}

// Use the service
func (s *UserService) CreateUser(user *User) error {
    // Validate
    if err := user.Validate(); err != nil {
        return err
    }

    // Save to database
    if err := s.repo.Create(user); err != nil {
        s.logger.Error("Failed to create user", logger.Field("error", err))
        return err
    }

    // Send welcome email
    s.email.Send(user.Email, "Welcome!", "Welcome to our app!")

    s.logger.Info("User created", logger.Field("userId", user.ID))
    return nil
}
```

## Wire Package Integration

go-infra supports Google Wire for automatic dependency injection:

### Step 1: Define Providers

```go
// internal/wire/wire.go
//go:build wireinject
// +build wireinject

package wire

import (
    "github.com/google/wire"
    "myapp/internal/repository"
    "myapp/internal/service"
)

// Provider sets
var RepositorySet = wire.NewSet(
    repository.NewUserRepository,
    repository.NewProductRepository,
)

var ServiceSet = wire.NewSet(
    service.NewUserService,
    service.NewProductService,
)

// Wire function
func InitializeApp() (*App, error) {
    wire.Build(
        // Database
        ProvideDatabaseConnection,

        // Repositories
        RepositorySet,

        // Services
        ServiceSet,

        // HTTP
        ProvideHTTPServer,

        // Application
        NewApp,
    )
    return nil, nil
}
```

### Step 2: Define Providers

```go
// internal/wire/providers.go
package wire

import (
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func ProvideDatabaseConnection() (*gorm.DB, error) {
    dsn := "host=localhost user=postgres password=postgres dbname=myapp"
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func ProvideHTTPServer() (fiber.App, fiber.Router) {
    return fiber.New()
}
```

### Step 3: Generate Wire Code

```bash
# Install wire
go install github.com/google/wire/cmd/wire@latest

# Generate wire_gen.go
cd internal/wire
wire
```

### Step 4: Use in Main

```go
// main.go
package main

import (
    "log"
    "myapp/internal/wire"
)

func main() {
    app, err := wire.InitializeApp()
    if err != nil {
        log.Fatal(err)
    }

    app.Run()
}
```

## Manual DI Container

If you prefer not to use Wire, create a manual container:

```go
// internal/container/container.go
package container

import (
    "gorm.io/gorm"
    "myapp/internal/repository"
    "myapp/internal/service"
    "github.com/phatnt199/go-infra/pkg/logger"
)

type Container struct {
    // Infrastructure
    DB     *gorm.DB
    Logger logger.Logger

    // Repositories
    UserRepo    repository.UserRepository
    ProductRepo repository.ProductRepository

    // Services
    UserService    *service.UserService
    ProductService *service.ProductService
}

func New(db *gorm.DB) *Container {
    c := &Container{
        DB:     db,
        Logger: logger.Default(),
    }

    // Initialize repositories
    c.UserRepo = repository.NewUserRepository(db)
    c.ProductRepo = repository.NewProductRepository(db)

    // Initialize services
    c.UserService = service.NewUserService(c.UserRepo, c.Logger)
    c.ProductService = service.NewProductService(c.ProductRepo, c.Logger)

    return c
}
```

Usage:

```go
// main.go
func main() {
    db := setupDatabase()
    container := container.New(db)

    app, router := fiber.New()

    // Register handlers with injected services
    handler.NewUserHandler(container.UserService).Register(router)
    handler.NewProductHandler(container.ProductService).Register(router)

    app.Listen(":3000")
}
```

## Interface-Based DI

Define interfaces for flexibility:

```go
// domain/repository/user.go
package repository

type UserRepository interface {
    Create(user *User) error
    FindByID(id string) (*User, error)
    FindByEmail(email string) (*User, error)
    Update(user *User) error
    Delete(id string) error
}

// infra/repository/user_gorm.go
package repository

type GormUserRepository struct {
    db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
    return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(user *User) error {
    return r.db.Create(user).Error
}

// Can easily swap with different implementation
// infra/repository/user_mongo.go
type MongoUserRepository struct {
    client *mongo.Client
}
```

## Scoped Dependencies

Some dependencies should be scoped to a request:

```go
// Middleware to inject request-scoped dependencies
func RequestScopeMiddleware(container *Container) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Create request-scoped logger with request ID
        requestID := c.Get("X-Request-ID")
        scopedLogger := container.Logger.With(
            logger.Field("requestId", requestID),
        )

        // Store in context
        c.Locals("logger", scopedLogger)

        return c.Next()
    }
}

// Handler uses scoped logger
func (h *UserHandler) Create(c *fiber.Ctx) error {
    logger := c.Locals("logger").(logger.Logger)
    logger.Info("Creating user")

    // ...
}
```

## Testing with DI

DI makes testing much easier:

```go
// Mock repository
type MockUserRepository struct {
    CreateFunc func(user *User) error
}

func (m *MockUserRepository) Create(user *User) error {
    return m.CreateFunc(user)
}

// Test service with mock
func TestUserService_CreateUser(t *testing.T) {
    // Setup mock
    mockRepo := &MockUserRepository{
        CreateFunc: func(user *User) error {
            user.ID = "test-id"
            return nil
        },
    }

    mockEmail := &MockEmailService{}
    mockLogger := logger.NewNoop()

    // Create service with mocks
    service := NewUserService(mockRepo, mockEmail, mockLogger)

    // Test
    user := &User{Name: "John", Email: "john@example.com"}
    err := service.CreateUser(user)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "test-id", user.ID)
}
```

## Circular Dependencies

Avoid circular dependencies:

```go
// ❌ Bad - circular dependency
package user

import "myapp/order"

type UserService struct {
    orderService *order.OrderService
}

// order/service.go
package order

import "myapp/user"

type OrderService struct {
    userService *user.UserService
}

// ✅ Good - use interfaces and events
package user

type OrderService interface {
    GetOrdersByUser(userID string) ([]*Order, error)
}

type UserService struct {
    orderService OrderService
}

// Or use event-driven approach
func (s *UserService) CreateUser(user *User) error {
    // Create user

    // Publish event instead of calling OrderService
    s.eventBus.Publish("user.created", UserCreatedEvent{UserID: user.ID})

    return nil
}
```

## Configuration Injection

Inject configuration objects:

```go
// config/config.go
type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
}

type AppConfig struct {
    Database DatabaseConfig
    JWT      JWTConfig
    Email    EmailConfig
}

// Load configuration
func Load() (*AppConfig, error) {
    var cfg AppConfig
    // Load from file, env vars, etc.
    return &cfg, nil
}

// Inject into providers
func ProvideDatabaseConnection(cfg *AppConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.Username,
        cfg.Database.Password,
        cfg.Database.Database,
    )
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
```

## Best Practices

### 1. Depend on Abstractions

```go
// ✅ Good - depend on interface
type UserService struct {
    repo UserRepository  // interface
}

// ❌ Bad - depend on concrete type
type UserService struct {
    repo *GormUserRepository
}
```

### 2. Keep Constructors Simple

```go
// ✅ Good
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

// ❌ Bad - doing too much in constructor
func NewUserService(dsn string) *UserService {
    db, _ := gorm.Open(postgres.Open(dsn))  // Bad!
    repo := NewUserRepository(db)
    return &UserService{repo: repo}
}
```

### 3. Use Option Pattern for Optional Dependencies

```go
type UserServiceOptions struct {
    Cache      CacheService
    RateLimit  RateLimiter
}

func NewUserService(
    repo UserRepository,
    opts *UserServiceOptions,
) *UserService {
    s := &UserService{repo: repo}

    if opts != nil {
        s.cache = opts.Cache
        s.rateLimit = opts.RateLimit
    }

    return s
}
```

## Real-World Example

Complete application with DI:

```go
// main.go
package main

import (
    "log"
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber"
    "github.com/phatnt199/go-infra/pkg/logger"
    "myapp/internal/config"
    "myapp/internal/repository"
    "myapp/internal/service"
    "myapp/internal/handler"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal(err)
    }

    // Initialize logger
    logger.Init()

    // Setup database
    db, err := setupDatabase(cfg.Database)
    if err != nil {
        log.Fatal(err)
    }

    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    productRepo := repository.NewProductRepository(db)

    // Initialize services
    userService := service.NewUserService(userRepo, logger.Default())
    productService := service.NewProductService(productRepo, logger.Default())

    // Setup HTTP server
    app, router := fiber.New()

    // Register handlers
    userHandler := handler.NewUserHandler(userService)
    userHandler.Register(router)

    productHandler := handler.NewProductHandler(productService)
    productHandler.Register(router)

    // Start server
    log.Fatal(app.Listen(":3000"))
}
```

## Next Steps

- Learn about [Module System](./modules)
- Explore [Configuration Management](./configuration)
- See [Testing Guide](../testing/unit-testing)
