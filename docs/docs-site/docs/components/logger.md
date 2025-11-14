---
sidebar_position: 1
---

# Logger

High-performance structured logging for Go applications.

## Features

- ✅ **Zero external imports** - Only import go-infra
- ✅ **High performance** - Built on Uber's Zap
- ✅ **Type-safe fields** - Auto-complete support
- ✅ **Context propagation** - Request tracing
- ✅ **Scoped logging** - Organize by component
- ✅ **Environment-based** - Auto-configure

## Quick Start

### Basic Usage

```go
import "github.com/phatnt199/go-infra/pkg/logger"

func main() {
    // Use default logger (auto-initialized)
    logger.Info("Application started")
    logger.Info("User logged in", logger.String("user_id", "12345"))

    defer logger.Sync()  // Flush logs before exit
}
```

### Custom Configuration

```go
func main() {
    config := &logger.Config{
        Environment:      "production",
        Level:            "info",
        OutputPaths:      []string{"stdout", "/var/log/app.log"},
        ErrorOutputPaths: []string{"stderr", "/var/log/error.log"},
        EnableCaller:     true,
        EnableStacktrace: true,
        Encoding:         "json",
        ServiceName:      "my-service",
    }

    if err := logger.Init(config); err != nil {
        panic(err)
    }
    defer logger.Sync()

    logger.Info("Application started with custom config")
}
```

## Log Levels

```go
// Debug - Detailed information for diagnosing problems
logger.Debug("Processing item", logger.Int("item_id", 123))

// Info - General informational messages
logger.Info("Server started", logger.Int("port", 8080))

// Warn - Warning messages
logger.Warn("High memory usage", logger.Int64("bytes", memUsage))

// Error - Error events
logger.Error("Database connection failed", logger.Err(err))

// Fatal - Severe errors (exits application)
logger.Fatal("Cannot start server", logger.Err(err))

// Panic - Severe errors (causes panic)
logger.Panic("Critical failure", logger.Err(err))
```

## Structured Fields

All field types are type-safe with auto-complete:

```go
// String fields
logger.String("key", "value")

// Numeric fields
logger.Int("count", 42)
logger.Int64("id", 123456789)
logger.Float64("price", 99.99)

// Boolean fields
logger.Bool("is_active", true)

// Time-related fields
logger.Time("created_at", time.Now())
logger.Duration("elapsed", 150*time.Millisecond)

// Error fields
logger.Err(err)

// Any type (uses reflection)
logger.Any("data", complexStruct)
```

## Structured Logging

```go
logger.Info("User action",
    logger.String("user_id", "user-123"),
    logger.String("action", "login"),
    logger.String("ip", "192.168.1.1"),
    logger.Duration("duration", 150*time.Millisecond),
)

// Output (JSON in production):
// {
//   "level": "info",
//   "timestamp": "2024-01-01T10:00:00Z",
//   "message": "User action",
//   "user_id": "user-123",
//   "action": "login",
//   "ip": "192.168.1.1",
//   "duration": 0.15
// }
```

## Logger with Persistent Fields

```go
// Create a logger with fields that persist across all logs
userLogger := logger.WithFields(
    logger.String("user_id", "user-456"),
    logger.String("session_id", "sess-789"),
)

// All logs from this logger include these fields
userLogger.Info("Viewing profile")
userLogger.Info("Updating settings")
userLogger.Info("Logging out")
```

## Scoped Logging

```go
// Create scoped loggers for different components
authLogger := logger.WithScope("auth")
dbLogger := logger.WithScope("database")
apiLogger := logger.WithScope("api")

authLogger.Info("User authentication started")
dbLogger.Info("Database connection established")
apiLogger.Info("API server listening", logger.Int("port", 8080))

// Chain scopes and fields
userAuthLogger := authLogger.WithFields(
    logger.String("user_id", "user-999"),
)
userAuthLogger.Info("Password verified")
```

## Context-Based Logging

```go
import "context"

func HandleRequest(w http.ResponseWriter, r *http.Request) {
    // Create a logger for this request
    requestLogger := logger.GetDefault().WithFields(
        logger.String("request_id", generateRequestID()),
        logger.String("method", r.Method),
        logger.String("path", r.URL.Path),
    )

    // Add logger to context
    ctx := requestLogger.ToContext(r.Context())

    // Pass context through your application
    ProcessOrder(ctx)
}

func ProcessOrder(ctx context.Context) {
    // Retrieve logger from context
    log := logger.FromContext(ctx)

    log.Info("Order processing started")
    // All logs will include request_id, method, path automatically
}
```

## Configuration

### Environment Variables

- `APP_ENV` - Sets environment (development, production)
- `SERVICE_NAME` - Sets service name in logs

### Config Options

```go
type Config struct {
    // Environment (development, production, staging)
    Environment string

    // Log level: debug, info, warn, error, fatal, panic
    Level string

    // Output paths for logs
    OutputPaths []string  // ["stdout", "/var/log/app.log"]

    // Output paths for errors
    ErrorOutputPaths []string  // ["stderr", "/var/log/error.log"]

    // Enable caller (file:line)
    EnableCaller bool

    // Enable stack trace on errors
    EnableStacktrace bool

    // Encoding: json or console
    Encoding string

    // Service name
    ServiceName string
}
```

### Default Configuration

```go
config := logger.DefaultConfig()
// Automatically configures based on APP_ENV:
// - Level: "debug" for dev, "info" for production
// - Encoding: "console" for dev, "json" for production
// - OutputPaths: ["stdout"]
// - ErrorOutputPaths: ["stderr"]
```

## File Output

```go
config := &logger.Config{
    Environment: "production",
    Level:       "info",
    OutputPaths: []string{
        "stdout",                    // Console
        "/var/log/app/app.log",      // Application logs
    },
    ErrorOutputPaths: []string{
        "stderr",                    // Console errors
        "/var/log/app/error.log",    // Error logs
    },
    Encoding: "json",
}
```

## HTTP Middleware Integration

```go
func LoggingMiddleware(log logger.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()

        // Create request logger
        requestLogger := log.WithFields(
            logger.String("request_id", generateRequestID()),
            logger.String("method", c.Method()),
            logger.String("path", c.Path()),
            logger.String("remote_addr", c.IP()),
        )

        // Add to context
        ctx := requestLogger.ToContext(c.Context())
        c.SetUserContext(ctx)

        requestLogger.Info("Request started")

        // Call next handler
        err := c.Next()

        // Log completion
        duration := time.Since(start)
        requestLogger.Info("Request completed",
            logger.Duration("duration", duration),
            logger.Int("status", c.Response().StatusCode()),
        )

        return err
    }
}
```

## Error Logging

```go
func FetchUser(id int) error {
    err := database.Find(id)
    if err != nil {
        // Log error with context and stack trace
        logger.Error("Failed to fetch user",
            logger.Err(err),
            logger.Int("user_id", id),
            logger.String("database", "postgres"),
        )
        return err
    }
    return nil
}
```

## Best Practices

### 1. Initialize at Startup

```go
func main() {
    // Initialize logger first
    config := logger.DefaultConfig()
    if err := logger.Init(config); err != nil {
        panic(err)
    }
    defer logger.Sync()

    // Rest of application
}
```

### 2. Use Context for Tracing

```go
// Add logger to context at entry points
ctx = logger.GetDefault().
    WithFields(logger.String("request_id", reqID)).
    ToContext(ctx)

// Retrieve logger from context in handlers
log := logger.FromContext(ctx)
```

### 3. Scope Loggers by Domain

```go
var (
    authLogger = logger.WithScope("auth")
    dbLogger   = logger.WithScope("database")
    apiLogger  = logger.WithScope("api")
)
```

### 4. Log Errors with Context

```go
if err != nil {
    logger.Error("Operation failed",
        logger.Err(err),
        logger.String("operation", "user_creation"),
        logger.String("user_id", userID),
    )
    return err
}
```

### 5. Use Appropriate Levels

- **Debug**: Development and troubleshooting
- **Info**: Normal application flow
- **Warn**: Potential issues
- **Error**: Errors requiring investigation
- **Fatal/Panic**: Critical errors (use sparingly)

## Performance

- **Zero-allocation** logging in hot paths
- **10x faster** than standard library log
- **4+ million logs/second**

Built on Uber's Zap with minimal abstraction overhead.

## Complete Example

```go
package main

import (
    "net/http"
    "time"
    "github.com/phatnt199/go-infra/pkg/logger"
)

func main() {
    // Initialize logger
    config := &logger.Config{
        Environment:      "production",
        Level:            "info",
        OutputPaths:      []string{"stdout", "/var/log/app.log"},
        ErrorOutputPaths: []string{"stderr", "/var/log/error.log"},
        EnableCaller:     true,
        Encoding:         "json",
        ServiceName:      "user-service",
    }

    if err := logger.Init(config); err != nil {
        panic(err)
    }
    defer logger.Sync()

    logger.Info("Application starting",
        logger.String("version", "1.0.0"),
        logger.String("environment", config.Environment),
    )

    // Setup HTTP server
    mux := http.NewServeMux()
    mux.HandleFunc("/users", GetUsersHandler)

    server := &http.Server{
        Addr:    ":8080",
        Handler: LoggingMiddleware(mux),
    }

    logger.Info("Server listening", logger.Int("port", 8080))

    if err := server.ListenAndServe(); err != nil {
        logger.Fatal("Server failed", logger.Err(err))
    }
}

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        requestLogger := logger.GetDefault().WithFields(
            logger.String("method", r.Method),
            logger.String("path", r.URL.Path),
        )

        ctx := requestLogger.ToContext(r.Context())
        r = r.WithContext(ctx)

        requestLogger.Info("Request started")
        next.ServeHTTP(w, r)

        requestLogger.Info("Request completed",
            logger.Duration("duration", time.Since(start)),
        )
    })
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
    log := logger.FromContext(r.Context())
    log.Info("Fetching users")
    // Handler logic
}
```

## Next Steps

- **[Crypto Package](./crypto)** - Security utilities
- **[Error Handling](../core-concepts/error-handling)** - Handle errors properly
- **[Observability](../observability/tracing)** - Add tracing
