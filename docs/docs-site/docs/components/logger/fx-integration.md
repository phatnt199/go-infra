---
sidebar_position: 5
---

# Fx Integration

Use the logger with Uber's Fx dependency injection framework.

## Overview

The go-infra logger has first-class support for [Uber Fx](https://uber-go.github.io/fx/), making it easy to:

- Provide logger to all components via dependency injection
- Configure logger from environment
- Share a single logger instance across your application
- Log Fx lifecycle events

## Quick Start

### Basic Fx Integration

```go
package main

import (
    "go.uber.org/fx"
    "github.com/phatnt199/go-infra/pkg/logger"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

func main() {
    app := fx.New(
        // Provide Zap logger module
        zaplogger.Module,

        // Use logger in your components
        fx.Invoke(runApp),
    )

    app.Run()
}

func runApp(log logger.Logger) {
    log.Info("Application started with Fx")
}
```

## Zap Logger Module

The `zaplogger.Module` provides:

```go
var Module = fx.Module("zapfx",
    fx.Provide(
        config.ProvideLogConfig,    // Load config from environment
        NewZapLogger,                // Create Zap logger
        fx.Annotate(
            NewZapLogger,
            fx.As(new(logger.Logger)), // Provide as logger.Logger interface
        ),
    ),
)
```

This module:

1. Loads logger configuration from environment variables
2. Creates a ZapLogger instance
3. Provides it as both `zaplogger.ZapLogger` and `logger.Logger` interfaces

### Using the Module

```go
import (
    "go.uber.org/fx"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

app := fx.New(
    zaplogger.Module,

    // Logger is now available for injection
    fx.Invoke(func(log logger.Logger) {
        log.Info("Logger injected")
    }),
)
```

## Provide Your Own Logger

If you want to provide a pre-configured logger instance:

```go
package main

import (
    "go.uber.org/fx"
    "github.com/phatnt199/go-infra/pkg/application/constants"
    "github.com/phatnt199/go-infra/pkg/logger"
    "github.com/phatnt199/go-infra/pkg/logger/config"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

func main() {
    // Create custom logger
    cfg := &config.LogOptions{
        LogLevel:      "debug",
        CallerEnabled: true,
        EnableTracing: false,
    }
    customLog := zaplogger.NewZapLogger(cfg, constants.DEV_ENV)

    app := fx.New(
        // Use ModuleFunc to provide existing logger
        zaplogger.ModuleFunc(customLog),

        fx.Invoke(runApp),
    )

    app.Run()
}

func runApp(log logger.Logger) {
    log.Info("Using custom logger")
}
```

### ModuleFunc

The `ModuleFunc` creates an Fx module from an existing logger:

```go
var ModuleFunc = func(l logger.Logger) fx.Option {
    return fx.Module(
        "zapfx",
        fx.Provide(config.ProvideLogConfig),
        fx.Supply(fx.Annotate(l, fx.As(new(logger.Logger)))),
        fx.Supply(fx.Annotate(l, fx.As(new(ZapLogger)))),
    )
}
```

## Empty Logger Module

For testing or when you don't want logging:

```go
import (
    "go.uber.org/fx"
    emptylogger "github.com/phatnt199/go-infra/pkg/logger/empty"
)

app := fx.New(
    emptylogger.Module, // No-op logger

    fx.Invoke(func(log logger.Logger) {
        log.Info("This won't output anything")
    }),
)
```

## Fx Logger Adapter

Log Fx lifecycle events with your logger:

```go
package main

import (
    "go.uber.org/fx"
    "github.com/phatnt199/go-infra/pkg/logger"
    fxlog "github.com/phatnt199/go-infra/pkg/logger/external/fxlog"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

func main() {
    app := fx.New(
        zaplogger.Module,

        // Add Fx logger adapter
        fx.WithLogger(func(log logger.Logger) fxevent.Logger {
            return fxlog.NewCustomFxLogger(log)
        }),

        fx.Invoke(runApp),
    )

    app.Run()
}

func runApp(log logger.Logger) {
    log.Info("App running")
}
```

Or use the convenience option:

```go
app := fx.New(
    zaplogger.Module,
    fxlog.FxLogger, // Automatically sets up Fx logging
    fx.Invoke(runApp),
)
```

### Fx Event Logging

The Fx logger adapter logs lifecycle events:

- **OnStart/OnStop hooks** - When hooks execute
- **Provided dependencies** - What's registered
- **Invoked functions** - What's being run
- **Errors** - Dependency resolution failures

Example output:

```
INFO | provided | constructor=NewUserRepository | type=*repository.UserRepository
INFO | provided | constructor=NewUserService | type=*service.UserService
INFO | invoking | function=runApp
INFO | OnStart hook executing | caller=server.Start
INFO | OnStart hook executed | caller=server.Start | runtime=150ms
INFO | started
```

## Inject Logger into Components

### Repository Example

```go
type UserRepository struct {
    db     *gorm.DB
    logger logger.Logger
}

func NewUserRepository(db *gorm.DB, log logger.Logger) *UserRepository {
    return &UserRepository{
        db:     db,
        logger: log,
    }
}

func (r *UserRepository) Create(user *User) error {
    r.logger.Infow("Creating user", logger.Fields{
        "email": user.Email,
    })

    err := r.db.Create(user).Error
    if err != nil {
        r.logger.Err("Failed to create user", err)
        return err
    }

    r.logger.Infow("User created", logger.Fields{
        "user_id": user.ID,
    })
    return nil
}
```

### Service Example

```go
type UserService struct {
    repo   *UserRepository
    logger logger.Logger
}

func NewUserService(repo *UserRepository, log logger.Logger) *UserService {
    return &UserService{
        repo:   repo,
        logger: log,
    }
}

func (s *UserService) Register(email, password string) error {
    s.logger.Infow("User registration started", logger.Fields{
        "email": email,
    })

    // Business logic...

    return nil
}
```

### HTTP Handler Example

```go
type UserHandler struct {
    service *UserService
    logger  logger.Logger
}

func NewUserHandler(service *UserService, log logger.Logger) *UserHandler {
    return &UserHandler{
        service: service,
        logger:  log,
    }
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    h.logger.Infow("Create user request", logger.Fields{
        "path": c.Path(),
        "method": c.Method(),
    })

    // Handle request...

    return nil
}
```

## Complete Example

```go
package main

import (
    "go.uber.org/fx"
    "github.com/gofiber/fiber/v2"
    customfiber "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
    "github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
    "github.com/phatnt199/go-infra/pkg/logger"
    fxlog "github.com/phatnt199/go-infra/pkg/logger/external/fxlog"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

func main() {
    app := fx.New(
        // Logger module
        zaplogger.Module,
        fxlog.FxLogger,

        // HTTP server module
        customfiber.Module,

        // Application components
        fx.Provide(
            NewUserRepository,
            NewUserService,
            NewUserHandler,
        ),

        // Setup routes
        fx.Invoke(setupRoutes),
    )

    app.Run()
}

// Repository
type UserRepository struct {
    logger logger.Logger
}

func NewUserRepository(log logger.Logger) *UserRepository {
    log.Info("UserRepository initialized")
    return &UserRepository{logger: log}
}

// Service
type UserService struct {
    repo   *UserRepository
    logger logger.Logger
}

func NewUserService(repo *UserRepository, log logger.Logger) *UserService {
    log.Info("UserService initialized")
    return &UserService{
        repo:   repo,
        logger: log,
    }
}

// Handler
type UserHandler struct {
    service *UserService
    logger  logger.Logger
}

func NewUserHandler(service *UserService, log logger.Logger) *UserHandler {
    log.Info("UserHandler initialized")
    return &UserHandler{
        service: service,
        logger:  log,
    }
}

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
    h.logger.Info("Fetching users")
    return c.JSON(fiber.Map{"users": []string{}})
}

// Setup routes
func setupRoutes(
    server contracts.HttpServer,
    handler *UserHandler,
    log logger.Logger,
) {
    log.Info("Setting up routes")

    server.RouteBuilder().RegisterHandler(func(instance interface{}) {
        if app, ok := instance.(*fiber.App); ok {
            app.Get("/users", handler.GetUsers)
            log.Info("Routes registered")
        }
    })
}
```

## Module Composition

Combine logger with other modules:

```go
import (
    "go.uber.org/fx"
    fxapp "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
    customfiber "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
    postgresgorm "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

func main() {
    appBuilder := fxapp.NewApplicationBuilder()

    // Use the built-in logger from appBuilder
    log := appBuilder.Logger()
    log.Info("Building application")

    // Or provide modules that include logger
    appBuilder.ProvideModule(zaplogger.Module)
    appBuilder.ProvideModule(customfiber.Module)
    appBuilder.ProvideModule(postgresgorm.Module)

    // Build and run
    app := appBuilder.Build()
    app.Run()
}
```

## Configuration with Fx

Load logger configuration from environment:

```bash
# Set environment variables
export LOG_OPTIONS_LEVEL="info"
export LOG_OPTIONS_CALLER_ENABLED="true"
export LOG_OPTIONS_ENABLE_TRACING="true"
```

```go
import (
    "go.uber.org/fx"
    "github.com/phatnt199/go-infra/pkg/logger/config"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

app := fx.New(
    zaplogger.Module, // Automatically loads config

    fx.Invoke(func(cfg *config.LogOptions, log logger.Logger) {
        log.Infow("Logger configured", logger.Fields{
            "level": cfg.LogLevel,
            "caller_enabled": cfg.CallerEnabled,
        })
    }),
)
```

## Testing with Fx

Use empty logger for testing:

```go
func TestApp(t *testing.T) {
    app := fxtest.New(t,
        emptylogger.Module, // No log output in tests

        fx.Provide(NewUserService),

        fx.Invoke(func(service *UserService) {
            // Test service
        }),
    )

    app.RequireStart()
    app.RequireStop()
}
```

## Next Steps

- **[External Adapters](./adapters.md)** - GORM and Fx logger adapters
- **[Best Practices](./best-practices.md)** - Production patterns
- **[Dependency Injection](../../core-concepts/dependency-injection.md)** - Learn more about Fx
