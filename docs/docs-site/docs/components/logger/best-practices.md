---
sidebar_position: 7
---

# Best Practices

Production-ready patterns for using the go-infra logger.

## Initialization

### Use Dependency Injection

**✅ Recommended:**

```go
// Inject logger via Fx
func NewUserService(log logger.Logger, repo *UserRepository) *UserService {
    return &UserService{
        log:  log,
        repo: repo,
    }
}
```

**❌ Avoid:**

```go
// Creating logger instances everywhere
func NewUserService(repo *UserRepository) *UserService {
    log := defaultlogger.GetLogger() // Don't do this
    return &UserService{
        log:  log,
        repo: repo,
    }
}
```

### Initialize Once

**✅ Recommended:**

```go
func main() {
    app := fx.New(
        zaplogger.Module, // Single logger instance for entire app
        fx.Invoke(runApp),
    )
    app.Run()
}
```

**❌ Avoid:**

```go
// Multiple logger instances
func someFunction() {
    log := defaultlogger.GetLogger() // Creates new instance
}

func anotherFunction() {
    log := defaultlogger.GetLogger() // Another instance
}
```

## Log Levels

### Use Appropriate Levels

**✅ Recommended:**

```go
// Debug - Development and troubleshooting
log.Debugw("Parsed request", logger.Fields{
    "headers": headers,
    "body": body,
})

// Info - Normal application flow
log.Infow("User registered", logger.Fields{
    "user_id": userID,
})

// Warn - Potential issues
log.Warnw("Retry attempt", logger.Fields{
    "attempt": 3,
    "max": 5,
})

// Error - Errors requiring attention
log.Errorw("Payment failed", logger.Fields{
    "order_id": orderID,
    "error": err.Error(),
})

// Fatal - Only for unrecoverable errors
if criticalResourceUnavailable {
    log.Fatal("Cannot start without database")
}
```

**❌ Avoid:**

```go
// Wrong level usage
log.Error("User logged in") // Not an error
log.Info("Database connection failed") // Should be Error
log.Debug("Server started") // Should be Info
log.Fatal("User not found") // Don't exit for business logic errors
```

## Structured Logging

### Use Structured Fields

**✅ Recommended:**

```go
// Structured logging with fields
log.Infow("Order processed", logger.Fields{
    "order_id": orderID,
    "amount": amount,
    "currency": "USD",
    "duration_ms": elapsed.Milliseconds(),
})
```

**❌ Avoid:**

```go
// String concatenation
log.Infof("Order %s processed with amount %f %s in %dms",
    orderID, amount, "USD", elapsed.Milliseconds())
```

### Use Consistent Field Names

**✅ Recommended:**

```go
// Consistent naming across application
log.Infow("User action", logger.Fields{
    "user_id": userID,      // Always user_id
    "action": "login",
    "timestamp": time.Now(),
})

log.Infow("User query", logger.Fields{
    "user_id": userID,      // Same field name
    "query_type": "profile",
})
```

**❌ Avoid:**

```go
// Inconsistent naming
log.Infow("User action", logger.Fields{
    "userId": userID,   // camelCase
})

log.Infow("User query", logger.Fields{
    "user_id": userID,  // snake_case
})

log.Infow("User update", logger.Fields{
    "id": userID,       // Too generic
})
```

## Error Logging

### Include Context

**✅ Recommended:**

```go
// Rich error context
err := db.Create(&user).Error
if err != nil {
    log.Errorw("Failed to create user", logger.Fields{
        "error": err.Error(),
        "user_email": user.Email,
        "operation": "user_registration",
        "database": "postgres",
    })
    return err
}
```

**❌ Avoid:**

```go
// Minimal context
if err != nil {
    log.Error("Error") // What error? Where?
    return err
}
```

### Don't Log and Return

**✅ Recommended:**

```go
// Log at the point where you handle the error
func (s *Service) ProcessOrder(orderID string) error {
    err := s.repo.Save(order)
    if err != nil {
        // Only log here if you're handling it
        // Otherwise, return and let caller decide
        return fmt.Errorf("save order: %w", err)
    }
    return nil
}

func (h *Handler) CreateOrder(c *fiber.Ctx) error {
    err := h.service.ProcessOrder(orderID)
    if err != nil {
        // Log once at the boundary
        h.log.Errorw("Order processing failed", logger.Fields{
            "order_id": orderID,
            "error": err.Error(),
        })
        return c.Status(500).JSON(errorResponse)
    }
    return nil
}
```

**❌ Avoid:**

```go
// Logging at every level
func (r *Repository) Save(order *Order) error {
    err := r.db.Create(order).Error
    if err != nil {
        r.log.Error("DB error", err) // Logged
        return err
    }
    return nil
}

func (s *Service) ProcessOrder(orderID string) error {
    err := s.repo.Save(order)
    if err != nil {
        s.log.Error("Save failed", err) // Logged again
        return err
    }
    return nil
}

func (h *Handler) CreateOrder(c *fiber.Ctx) error {
    err := h.service.ProcessOrder(orderID)
    if err != nil {
        h.log.Error("Processing failed", err) // Logged again!
        return c.Status(500).JSON(errorResponse)
    }
    return nil
}
```

## Performance

### Use Appropriate Methods

**✅ Recommended:**

```go
// For structured logging
log.Infow("User action", logger.Fields{
    "user_id": userID,
    "action": "login",
})

// For simple messages
log.Info("Application started")

// For formatted messages
log.Infof("Server listening on port %d", port)
```

**❌ Avoid:**

```go
// Expensive string concatenation
log.Info(fmt.Sprintf("User %s performed %s", userID, action))

// Using Any() for primitive types
log.Infow("User action", logger.Fields{
    "user_id": userID, // OK - uses optimized string field
    "data": userData,  // OK for complex types, uses reflection
})
```

### Avoid Debug Logs in Hot Paths

**✅ Recommended:**

```go
// Check log level before expensive operations
if log.LogLevel() == "debug" {
    // Only compute if debug is enabled
    expensiveDebugData := computeExpensiveData()
    log.Debugw("Debug data", logger.Fields{
        "data": expensiveDebugData,
    })
}
```

Note: The current logger interface doesn't expose `LogLevel()`. In practice, just avoid expensive operations in debug logs in production (set level to `info`).

## Configuration

### Environment-Based Configuration

**✅ Recommended:**

```go
// Production
cfg := &config.LogOptions{
    LogLevel:      "info",     // Less verbose
    CallerEnabled: false,      // Better performance
    EnableTracing: true,       // Observability
}
log := zaplogger.NewZapLogger(cfg, constants.PROD_ENV)

// Development
cfg := &config.LogOptions{
    LogLevel:      "debug",    // More verbose
    CallerEnabled: true,       // Helpful for debugging
    EnableTracing: false,      // Not needed locally
}
log := zaplogger.NewZapLogger(cfg, constants.DEV_ENV)
```

### Use Config from Environment

**✅ Recommended:**

```go
// Load from environment variables
app := fx.New(
    zaplogger.Module, // Automatically loads from env
    fx.Invoke(runApp),
)
```

## Context and Tracing

### Add Request Context

**✅ Recommended:**

```go
func (h *Handler) HandleRequest(c *fiber.Ctx) error {
    // Create request-scoped logger
    requestLog := h.log // Add request-specific fields if needed

    requestLog.Infow("Request started", logger.Fields{
        "method": c.Method(),
        "path": c.Path(),
        "request_id": c.Get("X-Request-ID"),
    })

    // Pass logger to service
    err := h.service.Process(requestLog, data)

    requestLog.Infow("Request completed", logger.Fields{
        "status": c.Response().StatusCode(),
    })

    return err
}
```

## Testing

### Use Empty Logger

**✅ Recommended:**

```go
import (
    "testing"
    "github.com/phatnt199/go-infra/pkg/logger/empty"
)

func TestUserService(t *testing.T) {
    log := empty.EmptyLogger
    service := NewUserService(log, repo)

    // Test without log noise
    result := service.Process(data)
    assert.Equal(t, expected, result)
}
```

### Or Use Real Logger for Integration Tests

```go
func TestIntegration(t *testing.T) {
    log := defaultlogger.GetLogger()

    // See logs for debugging
    service := NewUserService(log, repo)
    result := service.Process(data)
}
```

## Service Naming

### Use WithName

**✅ Recommended:**

```go
func NewAuthService(log logger.Logger) *AuthService {
    log.WithName("auth-service")
    return &AuthService{log: log}
}

// All logs will include service name
log.Info("User authenticated") // [auth-service] User authenticated
```

## Sensitive Data

### Never Log Sensitive Information

**✅ Recommended:**

```go
log.Infow("User registered", logger.Fields{
    "user_id": user.ID,
    "email": user.Email, // OK if not considered sensitive
})
```

**❌ Avoid:**

```go
log.Infow("User registered", logger.Fields{
    "password": user.Password,     // NEVER
    "credit_card": user.CreditCard, // NEVER
    "ssn": user.SSN,               // NEVER
    "api_key": user.APIKey,        // NEVER
})
```

### Redact Sensitive Data

```go
func sanitizeEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return "***"
    }
    return parts[0][:2] + "***@" + parts[1]
}

log.Infow("Password reset", logger.Fields{
    "email": sanitizeEmail(user.Email), // al***@example.com
})
```

## Production Checklist

- ✅ Use `info` or higher log level
- ✅ Enable JSON encoding for log aggregation
- ✅ Enable tracing for observability
- ✅ Disable caller information for performance (unless debugging)
- ✅ Use structured logging with consistent field names
- ✅ Log errors at application boundaries, not at every level
- ✅ Never log sensitive data (passwords, tokens, PII)
- ✅ Add request IDs for request tracing
- ✅ Use appropriate log levels
- ✅ Forward logs to centralized logging system

## Common Patterns

### Repository Pattern

```go
type UserRepository struct {
    db  *gorm.DB
    log logger.Logger
}

func (r *UserRepository) GetByID(id string) (*User, error) {
    r.log.Debugw("Fetching user", logger.Fields{"user_id": id})

    var user User
    err := r.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        r.log.Errorw("Failed to fetch user", logger.Fields{
            "user_id": id,
            "error": err.Error(),
        })
        return nil, err
    }

    r.log.Debugw("User fetched", logger.Fields{"user_id": id})
    return &user, nil
}
```

### Middleware Pattern

```go
func LoggingMiddleware(log logger.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()

        log.Infow("Request started", logger.Fields{
            "method": c.Method(),
            "path": c.Path(),
        })

        err := c.Next()

        log.Infow("Request completed", logger.Fields{
            "method": c.Method(),
            "path": c.Path(),
            "status": c.Response().StatusCode(),
            "duration_ms": time.Since(start).Milliseconds(),
        })

        return err
    }
}
```

## Next Steps

- **[Configuration](./configuration.md)** - Learn configuration options
- **[Usage Guide](./usage.md)** - Explore all logging methods
- **[Fx Integration](./fx-integration.md)** - Use with Uber Fx
