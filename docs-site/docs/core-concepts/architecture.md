---
sidebar_position: 1
---

# Architecture Overview

go-infra follows a **modular, layered architecture** designed for scalability and maintainability.

## Architecture Layers

```
┌─────────────────────────────────────────────┐
│           Application Layer                 │
│  (Your Business Logic & Domain Models)      │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│            Adapter Layer                    │
│  (HTTP, gRPC, WebSocket, CLI)               │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│           Component Layer                   │
│  (Authentication, Authorization, etc)       │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│          Infrastructure Layer               │
│  (Database, Cache, Queue, Storage)          │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│            Core Layer                       │
│  (Logger, Crypto, Utils, Config)            │
└─────────────────────────────────────────────┘
```

## Core Principles

### 1. Dependency Injection

go-infra uses **Uber Fx** for dependency injection, enabling:

- **Loose coupling** between components
- **Testability** through interface injection
- **Lifecycle management** for resources
- **Clear dependencies** in constructor functions

```go
type UserService struct {
    db     *gorm.DB
    logger logger.Logger
}

// Constructor with dependencies injected
func NewUserService(
    db *gorm.DB,
    logger logger.Logger,
) *UserService {
    return &UserService{
        db:     db,
        logger: logger,
    }
}
```

### 2. Module System

Components are organized as **modules** that can be composed:

```go
app := fxapp.NewApplicationBuilder().
    ProvideModule(fiber_adapter.Module).    // HTTP server
    ProvideModule(gorm.Module).             // Database
    ProvideModule(authentication.Module).   // Authentication
    Build()
```

Each module:

- Provides its dependencies
- Declares what it needs
- Can be enabled/disabled independently

### 3. Configuration Management

Environment-based configuration with validation:

```go
// Configuration is loaded automatically
type Config struct {
    HTTP     HTTPConfig     `mapstructure:"http"`
    Database DatabaseConfig `mapstructure:"database"`
    JWT      JWTConfig      `mapstructure:"jwt"`
}
```

Sources (in priority order):

1. Environment variables
2. `.env` file
3. Configuration files (`config.json`, `config.yaml`)
4. Default values

### 4. Clean Architecture

Follows Domain-Driven Design principles:

**Domain Layer** - Business logic and entities

```go
type User struct {
    entity.BaseModel
    Email    string
    Password string
}
```

**Repository Layer** - Data access abstraction

```go
type UserRepository interface {
    FindByEmail(email string) (*User, error)
    Create(user *User) error
}
```

**Service Layer** - Application logic

```go
type UserService struct {
    repo UserRepository
}

func (s *UserService) RegisterUser(email, password string) error {
    // Business logic here
}
```

**Handler Layer** - HTTP/gRPC handlers

```go
func (h *UserHandler) Register(c *fiber.Ctx) error {
    // Parse request, call service, return response
}
```

## Key Components

### Adapters

**Purpose**: Interface between external world and your application

| Adapter         | Purpose                       |
| --------------- | ----------------------------- |
| `fxapp`         | Application bootstrap with DI |
| `fiber_adapter` | HTTP server with Fiber        |
| `grpc`          | gRPC server                   |
| `websocket`     | WebSocket server              |

### Infrastructure

**Purpose**: External system integrations

| Package         | Purpose                  |
| --------------- | ------------------------ |
| `postgres/gorm` | PostgreSQL with GORM     |
| `redis`         | Redis cache              |
| `storage`       | File storage (S3, local) |
| `queue`         | Message queues           |

### Components

**Purpose**: Reusable business components

| Component        | Purpose              |
| ---------------- | -------------------- |
| `authentication` | Complete auth system |
| `authz`          | Authorization & RBAC |
| `metrics`        | Application metrics  |
| `migration`      | Database migrations  |

### Core

**Purpose**: Foundation utilities

| Package     | Purpose            |
| ----------- | ------------------ |
| `logger`    | Structured logging |
| `crypto`    | Security utilities |
| `utils`     | Common utilities   |
| `validator` | Input validation   |
| `mapper`    | Object mapping     |

## Request Flow

Here's how a typical HTTP request flows through the system:

```
1. HTTP Request
   ↓
2. Fiber Router → Finds matching route
   ↓
3. Middleware → Logging, Auth, Validation
   ↓
4. Handler → Parses request
   ↓
5. Service → Business logic
   ↓
6. Repository → Data access
   ↓
7. Database → Query execution
   ↓
8. Response → JSON serialization
   ↓
9. Client ← HTTP Response
```

Example with code:

```go
// 1-2: Route registered
app.Post("/users", handler.CreateUser)

// 3: Middleware applied
app.Use(authMiddleware.Handle())

// 4: Handler
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req CreateUserRequest
    c.BodyParser(&req)

    // 5: Call service
    user, err := h.service.CreateUser(req)
    if err != nil {
        return err
    }

    return c.JSON(user)
}

// 5: Service layer
func (s *UserService) CreateUser(req CreateUserRequest) (*User, error) {
    // Business logic
    user := &User{Email: req.Email}

    // 6: Call repository
    return s.repo.Create(user)
}

// 6-7: Repository & Database
func (r *UserRepository) Create(user *User) (*User, error) {
    err := r.db.Create(user).Error
    return user, err
}
```

## Dependency Graph

Components depend on each other in a clear hierarchy:

```
Application (Your Code)
    ↓ depends on
Adapters (HTTP, gRPC)
    ↓ depends on
Components (Auth, Metrics)
    ↓ depends on
Infrastructure (DB, Cache)
    ↓ depends on
Core (Logger, Crypto)
```

This ensures:

- **No circular dependencies**
- **Core has no external dependencies**
- **Infrastructure doesn't know about adapters**
- **Clear upgrade path**

## Error Handling

Consistent error handling throughout:

```go
import "emperror.dev/errors"

// Create errors with context
err := errors.New("user not found")

// Wrap errors
err = errors.Wrap(err, "failed to get user")

// Errors with codes
err = errors.NewWithDetails(
    "invalid input",
    "code", "INVALID_INPUT",
    "field", "email",
)
```

## Context Propagation

Request context flows through all layers:

```go
// Handler receives context
func (h *Handler) GetUser(c *fiber.Ctx) error {
    ctx := c.Context()

    // Context passed to service
    user, err := h.service.GetUser(ctx, id)
    return c.JSON(user)
}

// Service uses context
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    // Context passed to repository
    return s.repo.FindByID(ctx, id)
}

// Repository uses context for queries
func (r *Repository) FindByID(ctx context.Context, id string) (*User, error) {
    return r.db.WithContext(ctx).First(&User{}, id)
}
```

Benefits:

- Request tracing
- Cancellation
- Timeouts
- Logging context

## Lifecycle Management

go-infra handles application lifecycle:

```go
app := fxapp.NewApplicationBuilder().
    ProvideModule(fiber_adapter.Module).
    RegisterHook(onStart).   // Called on startup
    RegisterHook(onStop).    // Called on shutdown
    Build()

func onStart(lifecycle fx.Lifecycle, logger logger.Logger) {
    lifecycle.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            logger.Info("Application starting")
            return nil
        },
        OnStop: func(ctx context.Context) error {
            logger.Info("Application stopping")
            return nil
        },
    })
}
```

Lifecycle phases:

1. **Construction** - Create dependencies
2. **Initialization** - Run OnStart hooks
3. **Running** - Serve requests
4. **Shutdown** - Run OnStop hooks (graceful)

## Testing Architecture

The architecture makes testing easy:

```go
// Mock repository
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(user *User) error {
    args := m.Called(user)
    return args.Error(0)
}

// Test service with mock
func TestUserService_CreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo, logger.NewEmpty())

    mockRepo.On("Create", mock.Anything).Return(nil)

    err := service.CreateUser("test@example.com")
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

## Best Practices

### 1. Use Interfaces

Define interfaces for dependencies:

```go
type UserRepository interface {
    Create(user *User) error
    FindByID(id string) (*User, error)
}
```

### 2. Constructor Injection

Use constructors for dependencies:

```go
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

### 3. Layer Separation

Keep layers separate:

- Handlers don't access database directly
- Services don't know about HTTP
- Repositories only handle data access

### 4. Context Everywhere

Always pass context:

```go
func (s *Service) DoSomething(ctx context.Context) error
```

### 5. Error Wrapping

Add context to errors:

```go
if err != nil {
    return errors.Wrap(err, "failed to create user")
}
```

## Performance Considerations

### Database Connection Pooling

```go
// Automatically configured
db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
sqlDB, _ := db.DB()
sqlDB.SetMaxOpenConns(25)
sqlDB.SetMaxIdleConns(5)
```

### Logging Performance

- Zero-allocation logging in hot paths
- Structured logging more efficient than string formatting
- Log level filtering before serialization

### Dependency Injection

- One-time cost at startup
- No runtime overhead
- Compiled, not reflection-based

## Next Steps

- **[Module System](./modules)** - Deep dive into modules
- **[Dependency Injection](./dependency-injection)** - Master DI with Fx
- **[Configuration](./configuration)** - Advanced config patterns
