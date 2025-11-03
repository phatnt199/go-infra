# How to Use Logs in Go-Infra Project

This guide explains how to use the logging system in the go-infra project. The logging system is built on top of Uber's Zap library for high performance and structured logging.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Basic Usage](#basic-usage)
3. [Log Levels](#log-levels)
4. [Structured Logging](#structured-logging)
5. [Configuration](#configuration)
6. [Common Patterns](#common-patterns)
7. [Real-World Examples](#real-world-examples)
8. [Best Practices](#best-practices)

---

## Quick Start

### Import the Logger

```go
import "github.com/phatnt199/go-infra/pkg/logger"
```

### Using the Logger in Your Service/Handler

The logger is injected via dependency injection (using `go.uber.org/fx`):

```go
package myservice

import "github.com/phatnt199/go-infra/pkg/logger"

type MyService struct {
    logger logger.Logger
}

func NewMyService(logger logger.Logger) *MyService {
    return &MyService{
        logger: logger,
    }
}

func (s *MyService) DoSomething() {
    s.logger.Info("Starting operation")
    // ... your code ...
    s.logger.Info("Operation completed")
}
```

---

## Basic Usage

### Simple Info Logging

```go
// Simple message
s.logger.Info("User created successfully")

// With formatted message
s.logger.Infof("User created with ID: %d", userID)
```

### Debug Logging (Development Only)

```go
s.logger.Debug("Processing request payload")
s.logger.Debugf("Request data: %v", requestData)
```

### Warning Logging

```go
s.logger.Warn("Cache miss, falling back to database")
s.logger.Warnf("High memory usage: %d MB", memoryUsage)
```

### Error Logging

```go
s.logger.Error("Database connection failed")
s.logger.Errorf("Cannot parse JSON: %s", rawData)
```

### Fatal Logging (Stops Application)

```go
s.logger.Fatal("Cannot initialize critical service")
s.logger.Fatalf("Configuration error: %s", configPath)
```

---

## Log Levels

The project supports the following log levels (in order of severity):

| Level     | Usage                                   | Example                            |
| --------- | --------------------------------------- | ---------------------------------- |
| **DEBUG** | Detailed diagnostic information         | Request parsing, loop iterations   |
| **INFO**  | General informational messages          | Service started, request processed |
| **WARN**  | Warning messages for potential issues   | Deprecated API usage, cache miss   |
| **ERROR** | Error events that need investigation    | Failed request, database error     |
| **FATAL** | Severe errors that stop the application | Failed initialization              |

### Selecting Log Level

The log level is determined by the `LogConfig_LogLevel` environment variable or configuration:

- **Development**: Defaults to `debug` (see all logs)
- **Production**: Defaults to `info` (only info, warn, error, fatal)

```bash
# Set log level
export LogConfig_LogLevel=debug    # All logs
export LogConfig_LogLevel=info     # Info and above
export LogConfig_LogLevel=error    # Only errors and fatal
```

---

## Structured Logging

Structured logging adds context to your logs using fields. This is more powerful than string formatting.

### With Structured Fields

```go
// Instead of:
s.logger.Infof("User login from %s in %d ms", ipAddress, duration)

// Use:
s.logger.Infow("User login successful",
    logger.String("user_id", userID),
    logger.String("ip_address", ipAddress),
    logger.Int64("duration_ms", duration),
)
```

### Available Field Types

```go
import "github.com/phatnt199/go-infra/pkg/logger"

// String fields
logger.String("key", "value")
logger.String("user_id", "user-123")

// Numeric fields
logger.Int("count", 42)
logger.Int64("user_id", 123456)
logger.Float64("price", 99.99)

// Boolean fields
logger.Bool("is_active", true)

// Time fields
logger.Time("created_at", time.Now())
logger.Duration("elapsed", 150*time.Millisecond)

// Error fields
logger.Error("operation failed")  // Built-in error handling
logger.Err(err)                   // Log an error

// Any type
logger.Any("config", configObj)
```

### Error Logging with Context

```go
// Log error with context fields
if err != nil {
    s.logger.Errorw("Failed to fetch user",
        logger.Int64("user_id", userID),
        logger.String("database", "postgres"),
        logger.Err(err),
    )
    return err
}

// Or use Errorf for quick formatting
if err != nil {
    s.logger.Errorf("Operation failed for user %d: %v", userID, err)
    return err
}
```

---

## Configuration

### Configuration Structure

The logger is configured through the `LogOptions` struct:

```go
type LogOptions struct {
    LogLevel       string          // "debug", "info", "warn", "error", "fatal"
    LogType        models.LogType  // Zap (0) or Logrus (1)
    CallerEnabled  bool            // Show file:line in logs
    EnableTracing  bool            // Add logs as events to tracing
}
```

### Configuration File

Configuration is typically done via YAML or environment variables:

```yaml
# config.yaml
logConfig:
  level: debug # or info, warn, error
  logType: 0 # 0 = Zap, 1 = Logrus
  callerEnabled: true # Show caller info
  enableTracing: true # Integration with tracing
```

### Environment Variables

```bash
# Set log level
export LogConfig_LogLevel=debug

# Set log type (0=Zap, 1=Logrus)
export LogConfig_LogType=0

# Enable caller information
export LogConfig_CallerEnabled=true

# Enable tracing integration
export LogConfig_EnableTracing=true

# Application environment
export APP_ENV=development         # development, production, staging
```

### Automatic Configuration Based on Environment

The logger automatically adjusts based on `APP_ENV`:

**Development Environment:**

- Format: Console (human-readable)
- Level: debug
- Colors: Enabled
- Caller: Full file path

**Production Environment:**

- Format: JSON (machine-readable)
- Level: info
- Colors: Disabled
- Caller: Short file path

```bash
# Development (console format, debug level)
export APP_ENV=development

# Production (JSON format, info level)
export APP_ENV=production
```

---

## Common Patterns

### Pattern 1: Dependency Injection in Handlers

```go
package handlers

import (
    "github.com/phatnt199/go-infra/pkg/logger"
)

type UserHandler struct {
    logger logger.Logger
    // ... other dependencies
}

func NewUserHandler(logger logger.Logger) *UserHandler {
    return &UserHandler{
        logger: logger,
    }
}

func (h *UserHandler) GetUser(ctx context.Context, id string) (*User, error) {
    h.logger.Debugf("Fetching user with ID: %s", id)

    // ... fetch user ...

    if err != nil {
        h.logger.Errorw("Failed to fetch user",
            logger.String("user_id", id),
            logger.Err(err),
        )
        return nil, err
    }

    h.logger.Infof("User %s fetched successfully", id)
    return user, nil
}
```

### Pattern 2: Transaction Logging (Real Example from Project)

From `pkg/infra/postgres/gorm/pipelines/mediator_transaction_pipeline.go`:

```go
type mediatorTransactionPipeline struct {
    logger logger.Logger
    db     *gorm.DB
}

func (m *mediatorTransactionPipeline) Handle(ctx context.Context, request interface{}, next mediatr.RequestHandlerFunc) (interface{}, error) {
    requestName := typeMapper.GetSnakeTypeName(request)

    // Log transaction start
    m.logger.Infof("beginning database transaction for request `%s`", requestName)

    tx := m.db.WithContext(ctx).Begin()

    defer func() {
        if r := recover(); r != nil {
            // Log panic and rollback
            m.logger.Errorf(
                "panic in the transaction, rolling back with message: %+v", r)
            tx.WithContext(ctx).Rollback()
        }
    }()

    middlewareResponse, err := next(ctx)

    if err != nil {
        // Log error and rollback
        m.logger.Errorf(
            "rolling back transaction for request `%s`: %v",
            requestName, err)
        tx.WithContext(ctx).Rollback()
        return nil, err
    }

    // Log commit
    m.logger.Infof("committing transaction for request `%s`", requestName)

    if err = tx.WithContext(ctx).Commit().Error; err != nil {
        m.logger.Errorf("transaction commit error: %v", err)
    }

    return middlewareResponse, nil
}
```

### Pattern 3: Error Handling in Application

From `pkg/adapter/fxapp/error_handler.go`:

```go
type FxErrorHandler struct {
    logger logger.Logger
}

func NewFxErrorHandler(logger logger.Logger) *FxErrorHandler {
    return &FxErrorHandler{logger: logger}
}

func (h *FxErrorHandler) HandleError(e error) {
    h.logger.Error(e)
}
```

### Pattern 4: Service Lifecycle Logging

From `pkg/adapter/fxapp/application.go`:

```go
func (a *application) Stop() error {
    if a.fxapp == nil {
        a.logger.Fatal("Failed to stop because application not started.")
    }

    ctx, cancel := context.WithTimeout(
        context.Background(),
        60*time.Second,
    )
    defer cancel()

    return a.fxapp.Stop(ctx)
}
```

### Pattern 5: Database Connection Logging

From `pkg/infra/postgres/pgx/postgres_pgx_fx.go`:

```go
// Log successful connection closure
logger.Info("Pgx postgres connection closed gracefully")

// Log initialization or critical events
logger.Fatalf("Cannot initialize postgres connection: %v", err)
```

---

## Real-World Examples

### Example 1: User Service with Structured Logging

```go
package services

import (
    "context"
    "github.com/phatnt199/go-infra/pkg/logger"
)

type UserService struct {
    logger logger.Logger
    repo   UserRepository
}

func NewUserService(logger logger.Logger, repo UserRepository) *UserService {
    return &UserService{
        logger: logger,
        repo:   repo,
    }
}

func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Log operation start with request ID
    s.logger.Infow("Creating new user",
        logger.String("email", req.Email),
        logger.String("operation", "create_user"),
    )

    // Validate input
    if err := req.Validate(); err != nil {
        s.logger.Errorw("Invalid user creation request",
            logger.String("email", req.Email),
            logger.Err(err),
        )
        return nil, err
    }

    // Create user
    user := &User{
        Email: req.Email,
        Name:  req.Name,
    }

    if err := s.repo.Save(ctx, user); err != nil {
        // Log error with context
        s.logger.Errorw("Failed to save user",
            logger.String("email", req.Email),
            logger.String("database", "postgres"),
            logger.Err(err),
        )
        return nil, err
    }

    // Log success
    s.logger.Infow("User created successfully",
        logger.String("user_id", user.ID),
        logger.String("email", user.Email),
    )

    return user, nil
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    s.logger.Debugf("Fetching user: %s", userID)

    user, err := s.repo.FindByID(ctx, userID)

    if err != nil {
        s.logger.Errorw("User not found",
            logger.String("user_id", userID),
            logger.Err(err),
        )
        return nil, err
    }

    s.logger.Debugf("User found: %s", userID)
    return user, nil
}
```

### Example 2: HTTP Handler with Request Logging

```go
package handlers

import (
    "net/http"
    "time"
    "github.com/phatnt199/go-infra/pkg/logger"
)

type OrderHandler struct {
    logger       logger.Logger
    orderService OrderService
}

func NewOrderHandler(logger logger.Logger, orderService OrderService) *OrderHandler {
    return &OrderHandler{
        logger:       logger,
        orderService: orderService,
    }
}

func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
    start := time.Now()

    // Extract order data
    var order Order
    if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
        h.logger.Errorw("Failed to decode order request",
            logger.String("path", r.URL.Path),
            logger.Err(err),
        )
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Log operation start
    h.logger.Infow("Processing order",
        logger.String("order_id", order.ID),
        logger.Float64("total", order.Total),
        logger.Int("items", len(order.Items)),
    )

    // Place order
    result, err := h.orderService.PlaceOrder(r.Context(), &order)
    if err != nil {
        h.logger.Errorw("Failed to place order",
            logger.String("order_id", order.ID),
            logger.String("error_type", fmt.Sprintf("%T", err)),
            logger.Err(err),
        )
        http.Error(w, "Order processing failed", http.StatusInternalServerError)
        return
    }

    // Log success with duration
    h.logger.Infow("Order placed successfully",
        logger.String("order_id", result.ID),
        logger.String("status", result.Status),
        logger.Duration("processing_time", time.Since(start)),
    )

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

### Example 3: Database Migration Logging

```go
package migrations

import (
    "github.com/phatnt199/go-infra/pkg/logger"
    "gorm.io/gorm"
)

type Migration struct {
    logger logger.Logger
    db     *gorm.DB
}

func NewMigration(logger logger.Logger, db *gorm.DB) *Migration {
    return &Migration{
        logger: logger,
        db:     db,
    }
}

func (m *Migration) RunMigrations() error {
    m.logger.Info("Starting database migrations")

    // Create users table
    if err := m.db.AutoMigrate(&User{}); err != nil {
        m.logger.Errorw("Failed to migrate users table",
            logger.Err(err),
        )
        return err
    }
    m.logger.Info("Users table migrated successfully")

    // Create orders table
    if err := m.db.AutoMigrate(&Order{}); err != nil {
        m.logger.Errorw("Failed to migrate orders table",
            logger.Err(err),
        )
        return err
    }
    m.logger.Info("Orders table migrated successfully")

    m.logger.Info("All database migrations completed successfully")
    return nil
}
```

---

## Best Practices

### ✅ DO

1. **Use structured logging with fields**

   ```go
   // Good
   s.logger.Infow("User login successful",
       logger.String("user_id", userID),
       logger.String("ip_address", ipAddress),
       logger.Duration("response_time", duration),
   )
   ```

2. **Include request/correlation IDs**

   ```go
   s.logger.Infow("Processing request",
       logger.String("request_id", requestID),
       logger.String("operation", "create_user"),
   )
   ```

3. **Log at the right level**

   ```go
   s.logger.Debug("Parsed request body")           // Development
   s.logger.Info("User created successfully")      // Important events
   s.logger.Warn("Database connection slow")       // Potential issues
   s.logger.Error("Database connection failed")    // Errors
   ```

4. **Include error context**

   ```go
   if err != nil {
       s.logger.Errorw("Operation failed",
           logger.String("operation", "fetch_user"),
           logger.String("user_id", userID),
           logger.Err(err),
       )
   }
   ```

5. **Use semantic field names**

   ```go
   // Good - clear what the value represents
   logger.String("user_id", userID)
   logger.String("database", "postgres")
   logger.Duration("query_time", duration)

   // Avoid - unclear meaning
   logger.String("val1", userID)
   logger.String("db", "postgres")
   ```

### ❌ DON'T

1. **Don't use string formatting for complex data**

   ```go
   // Bad
   s.logger.Infof("User: %+v, Order: %+v", user, order)

   // Good
   s.logger.Infow("Processing transaction",
       logger.Any("user", user),
       logger.Any("order", order),
   )
   ```

2. **Don't log sensitive information**

   ```go
   // Bad
   s.logger.Infof("User password reset: %s", password)

   // Good
   s.logger.Info("User password reset successfully")
   ```

3. **Don't use Fatal/Panic for recoverable errors**

   ```go
   // Bad - kills the entire application
   s.logger.Fatal("Failed to fetch optional cache")

   // Good - recoverable
   s.logger.Warn("Failed to fetch cache, using fallback")
   ```

4. **Don't log the same error multiple times in a chain**

   ```go
   // Bad - error logged multiple times
   if err != nil {
       logger.Error("Database error", err)   // First log
       return fmt.Errorf("db error: %w", err)  // Returns
   }
   // Handler logs again
   s.logger.Error("Operation failed", err)  // Second log

   // Good - log once at appropriate level
   if err != nil {
       return fmt.Errorf("db error: %w", err)
   }
   // Log at handler level
   s.logger.Errorw("Operation failed", logger.Err(err))
   ```

5. **Don't use Debug logs in production code paths**

   ```go
   // Bad - creates noise in development
   s.logger.Debugf("Processing item %d", i)  // Inside a loop

   // Good - log periodically or on important events
   if i % 1000 == 0 {
       s.logger.Debugf("Processed %d items", i)
   }
   ```

---

## Output Examples

### Development Console Output

```
[INFO]  | 2025-01-15T10:30:45.123+0700 | main.go:42  | Application started | [SERVICE]=MyApp
[DEBUG] | 2025-01-15T10:30:46.456+0700 | handler.go:15 | Processing request with ID: req-123 | [SERVICE]=MyApp
[INFO]  | 2025-01-15T10:30:47.789+0700 | service.go:28 | User created successfully | user_id=usr-456 email=user@example.com | [SERVICE]=MyApp
[WARN]  | 2025-01-15T10:30:48.101+0700 | cache.go:55  | Cache miss for key: user-456 | [SERVICE]=MyApp
[ERROR] | 2025-01-15T10:30:49.234+0700 | db.go:72    | Failed to fetch user | user_id=usr-789 error=connection timeout | [SERVICE]=MyApp
```

### Production JSON Output

```json
{
  "level": "info",
  "timestamp": "2025-01-15T10:30:47.789+0700",
  "caller": "service.go:28",
  "message": "User created successfully",
  "service": "MyApp",
  "user_id": "usr-456",
  "email": "user@example.com"
}

{
  "level": "error",
  "timestamp": "2025-01-15T10:30:49.234+0700",
  "caller": "db.go:72",
  "message": "Failed to fetch user",
  "service": "MyApp",
  "user_id": "usr-789",
  "error": "connection timeout"
}
```

---

## Common Issues & Solutions

### Issue: Logger Not Showing Logs

**Problem:** Logs are not appearing in the console.

**Solution:** Check the log level configuration.

```bash
# Set log level to debug to see all logs
export LogConfig_LogLevel=debug
```

### Issue: No Request Context Available

**Problem:** Can't pass logger through the application.

**Solution:** Use fx dependency injection to provide logger at startup.

```go
// In your fx.Options
fx.Provide(func(logger logger.Logger) *MyService {
    return NewMyService(logger)
})
```

### Issue: Performance Concerns

**Problem:** Logging is using too much CPU/memory.

**Solution:**

- Use `Debug` logs for detailed diagnostics (they're disabled in production)
- Avoid logging in tight loops
- Use structured fields instead of reflection-based `Any()`

```go
// Bad - in a loop with reflection
for _, item := range items {
    logger.Info("Processing", logger.Any("item", item))  // Slow
}

// Good - only log periodically
for i, item := range items {
    if i % 1000 == 0 {
        logger.Infof("Processed %d items", i)
    }
}
```

---

## See Also

- [Logger Package Documentation](../pkg/logger/README.md)
- [Go-Infra Project](https://github.com/phatnt199/go-infra)
- [Zap Logger Documentation](https://pkg.go.dev/go.uber.org/zap)
- [Structured Logging Best Practices](https://www.kartar.net/2015/12/structured-logging/)

---

## Quick Reference

```go
// Import
import "github.com/phatnt199/go-infra/pkg/logger"

// Simple logs
logger.Debug(msg)
logger.Info(msg)
logger.Warn(msg)
logger.Error(msg)
logger.Fatal(msg)

// Formatted logs
logger.Debugf("format", args...)
logger.Infof("format", args...)
logger.Warnf("format", args...)
logger.Errorf("format", args...)
logger.Fatalf("format", args...)

// Structured logs with fields
logger.Infow(msg,
    logger.String("key", "value"),
    logger.Int("count", 42),
    logger.Err(error),
)

// Field types
logger.String(key, value)
logger.Int(key, value)
logger.Int64(key, value)
logger.Float64(key, value)
logger.Bool(key, value)
logger.Time(key, value)
logger.Duration(key, value)
logger.Err(error)
logger.Any(key, value)

// Error specific
logger.Error(msg, err)              // Simple error
logger.Errorf("format: %v", err)   // Formatted error
logger.Errorw(msg, logger.Err(err)) // Structured error
```

---

**Happy Logging! 📝**
