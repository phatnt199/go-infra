---
sidebar_position: 4
---

# Usage Guide

Comprehensive guide to using the go-infra logger.

## Getting a Logger Instance

### Default Logger

```go
import defaultlogger "github.com/phatnt199/go-infra/pkg/logger/default_logger"

log := defaultlogger.GetLogger()
```

### Via Dependency Injection

```go
func MyService(log logger.Logger) {
    log.Info("Service initialized")
}
```

### Create Logger Explicitly

```go
import (
    "github.com/phatnt199/go-infra/pkg/application/constants"
    "github.com/phatnt199/go-infra/pkg/logger/config"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

cfg := &config.LogOptions{
    LogLevel:      "info",
    CallerEnabled: true,
}
log := zaplogger.NewZapLogger(cfg, constants.DEV_ENV)
```

## Basic Logging Methods

### Simple Messages

Log simple messages without formatting:

```go
log.Debug("This is a debug message")
log.Info("This is an info message")
log.Warn("This is a warning message")
log.Error("This is an error message")
```

### Formatted Messages

Use Printf-style formatting with `*f` methods:

```go
userName := "Alice"
userID := 123

log.Debugf("User %s (ID: %d) accessed debug endpoint", userName, userID)
log.Infof("User %s logged in successfully", userName)
log.Warnf("User %s exceeded rate limit", userName)
log.Errorf("Failed to process request for user %s: %v", userName, err)
```

### Structured Logging

Use `*w` methods with field maps for structured logging:

```go
import "github.com/phatnt199/go-infra/pkg/logger"

log.Debugw("Debug information", logger.Fields{
    "component": "auth",
    "operation": "verify",
})

log.Infow("User logged in", logger.Fields{
    "user_id": "user-123",
    "ip": "192.168.1.1",
    "timestamp": time.Now(),
})

log.Errorw("Database query failed", logger.Fields{
    "query": sqlQuery,
    "error": err.Error(),
    "duration_ms": elapsed.Milliseconds(),
})
```

## Error Logging

### Err Method

Log an error with a message:

```go
err := database.Connect()
if err != nil {
    log.Err("Failed to connect to database", err)
}
```

### WarnMsg Method

Log a warning with an error:

```go
err := cache.Set(key, value)
if err != nil {
    log.WarnMsg("Cache write failed, continuing without cache", err)
}
```

### Error with Fields

Use `Errorw` for structured error logging:

```go
log.Errorw("Payment processing failed", logger.Fields{
    "order_id": orderID,
    "amount": amount,
    "error": err.Error(),
    "payment_gateway": "stripe",
})
```

## Field Types

The logger automatically converts Go types to appropriate Zap fields:

```go
log.Infow("Type examples", logger.Fields{
    // Strings
    "name": "Alice",

    // Integers
    "age": 30,
    "count": int64(1000000),

    // Floats
    "price": 99.99,

    // Booleans
    "active": true,

    // Time
    "created_at": time.Now(),
    "duration": 150 * time.Millisecond,

    // Errors
    "error": err,

    // Complex types (uses reflection)
    "user": userStruct,
})
```

Supported types with optimized field conversion:

- `string`, `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `bool`, `float32`, `float64`
- `time.Duration`, `time.Time`
- `error`
- Any other type (uses reflection - slower)

## Special Logging Methods

### Printf

Compatible with standard library `log` package:

```go
log.Printf("Server listening on port %d", port)
```

### WithName

Add a service name prefix to all logs:

```go
log.WithName("auth-service")
log.Info("Service started") // Logs: [auth-service] Service started
```

### Configure

Access internal logger for advanced configuration:

```go
log.Configure(func(internalLog interface{}) {
    if zapLog, ok := internalLog.(*zap.Logger); ok {
        // Add global fields, configure options, etc.
        zapLog = zapLog.With(zap.String("service", "my-service"))
    }
})
```

## Log Levels

### Debug

Detailed information for diagnosing problems:

```go
log.Debug("Entering function processOrder")
log.Debugf("Processing order %s with %d items", orderID, itemCount)
log.Debugw("Request details", logger.Fields{
    "headers": headers,
    "body": body,
})
```

### Info

General informational messages about application flow:

```go
log.Info("Application started successfully")
log.Infof("Server listening on port %d", 8080)
log.Infow("User registered", logger.Fields{
    "user_id": userID,
    "email": email,
})
```

### Warn

Warnings about potentially harmful situations:

```go
log.Warn("High memory usage detected")
log.Warnf("Retry attempt %d of %d", attempt, maxRetries)
log.WarnMsg("Cache unavailable", err)
```

### Error

Error events that need attention:

```go
log.Error("Failed to save user data")
log.Errorf("Database connection failed after %d attempts", attempts)
log.Err("Payment processing error", err)
log.Errorw("API request failed", logger.Fields{
    "endpoint": endpoint,
    "status_code": statusCode,
    "error": err.Error(),
})
```

### Fatal

Severe errors that cause application exit:

```go
if err := initDatabase(); err != nil {
    log.Fatal("Cannot start without database connection")
    // Application exits here
}

log.Fatalf("Critical configuration error: %v", err)
```

:::danger
`Fatal` methods call `os.Exit(1)` after logging. Use only for unrecoverable errors during startup.
:::

## gRPC Logging

Special methods for gRPC middleware logging:

### Server Middleware Logger

```go
log.GrpcMiddlewareAccessLogger(
    method,      // gRPC method name
    duration,    // Request duration
    metadata,    // Request metadata
    err,         // Error if any
)
```

Example:

```go
log.GrpcMiddlewareAccessLogger(
    "/user.UserService/GetUser",
    150*time.Millisecond,
    map[string][]string{
        "authorization": {"Bearer token..."},
    },
    nil,
)
```

### Client Interceptor Logger

```go
log.GrpcClientInterceptorLogger(
    method,      // gRPC method name
    request,     // Request object
    response,    // Response object
    duration,    // Call duration
    metadata,    // Call metadata
    err,         // Error if any
)
```

Example:

```go
log.GrpcClientInterceptorLogger(
    "/user.UserService/CreateUser",
    &CreateUserRequest{Name: "Alice"},
    &CreateUserResponse{ID: "123"},
    100*time.Millisecond,
    metadata,
    nil,
)
```

## Practical Examples

### HTTP Request Logging

```go
func LoggingMiddleware(log logger.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            log.Infow("Request started", logger.Fields{
                "method": r.Method,
                "path": r.URL.Path,
                "remote_addr": r.RemoteAddr,
            })

            next.ServeHTTP(w, r)

            log.Infow("Request completed", logger.Fields{
                "method": r.Method,
                "path": r.URL.Path,
                "duration_ms": time.Since(start).Milliseconds(),
            })
        })
    }
}
```

### Database Operation Logging

```go
func (r *UserRepository) GetUser(id string) (*User, error) {
    r.log.Infow("Fetching user", logger.Fields{
        "user_id": id,
    })

    var user User
    err := r.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        r.log.Errorw("Failed to fetch user", logger.Fields{
            "user_id": id,
            "error": err.Error(),
        })
        return nil, err
    }

    r.log.Infow("User fetched successfully", logger.Fields{
        "user_id": id,
        "email": user.Email,
    })

    return &user, nil
}
```

### Background Job Logging

```go
func ProcessJob(log logger.Logger, job *Job) error {
    log.Infow("Job processing started", logger.Fields{
        "job_id": job.ID,
        "job_type": job.Type,
    })

    start := time.Now()
    err := job.Execute()
    duration := time.Since(start)

    if err != nil {
        log.Errorw("Job processing failed", logger.Fields{
            "job_id": job.ID,
            "duration_ms": duration.Milliseconds(),
            "error": err.Error(),
        })
        return err
    }

    log.Infow("Job processing completed", logger.Fields{
        "job_id": job.ID,
        "duration_ms": duration.Milliseconds(),
    })

    return nil
}
```

### Startup Logging

```go
func main() {
    log := defaultlogger.GetLogger()

    log.WithName("my-service")

    log.Info("Application starting...")

    // Load configuration
    cfg, err := loadConfig()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    log.Infow("Configuration loaded", logger.Fields{
        "env": cfg.Environment,
        "port": cfg.Port,
    })

    // Initialize database
    db, err := initDatabase(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    log.Info("Database connection established")

    // Start server
    log.Infof("Server listening on port %d", cfg.Port)
    if err := startServer(cfg.Port); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

## Testing with Empty Logger

For tests where you don't want log output:

```go
import (
    "testing"
    "github.com/phatnt199/go-infra/pkg/logger/empty"
)

func TestUserService(t *testing.T) {
    log := empty.EmptyLogger
    service := NewUserService(log)

    // Test without log noise
    result := service.DoSomething()
    // assertions...
}
```

## Next Steps

- **[Fx Integration](./fx-integration.md)** - Use logger with Uber Fx
- **[External Adapters](./adapters.md)** - GORM and Fx adapters
- **[Best Practices](./best-practices.md)** - Production patterns
