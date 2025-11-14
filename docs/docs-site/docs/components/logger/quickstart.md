---
sidebar_position: 2
---

# Quick Start

Get started with go-infra logger in minutes.

## Installation

```bash
go get github.com/phatnt199/go-infra
```

## Basic Usage

### Option 1: Default Logger (Simplest)

For quick prototyping or simple applications:

```go
package main

import (
    defaultlogger "github.com/phatnt199/go-infra/pkg/logger/default_logger"
)

func main() {
    // Get the default logger (auto-initializes on first call)
    log := defaultlogger.GetLogger()

    // Basic logging
    log.Info("Application started")
    log.Debug("Debug information")
    log.Warn("Warning message")
    log.Error("Error occurred")
}
```

The default logger automatically initializes with Zap logger when you first call `GetLogger()`.

### Option 2: Zap Logger (Explicit Configuration)

For more control over configuration:

```go
package main

import (
    "github.com/phatnt199/go-infra/pkg/application/constants"
    "github.com/phatnt199/go-infra/pkg/logger"
    "github.com/phatnt199/go-infra/pkg/logger/config"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

func main() {
    // Create logger configuration
    cfg := &config.LogOptions{
        LogLevel:      "info",
        LogType:       0, // Zap logger
        CallerEnabled: true,
        EnableTracing: false,
    }

    // Create logger instance
    log := zaplogger.NewZapLogger(cfg, constants.DEV_ENV)

    // Use the logger
    log.Info("Application started with custom configuration")
}
```

### Option 3: Fx Dependency Injection (Recommended for Production)

For production applications using Uber Fx:

```go
package main

import (
    "go.uber.org/fx"
    "github.com/phatnt199/go-infra/pkg/logger"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

func main() {
    app := fx.New(
        // Provide logger module
        zaplogger.Module,

        // Use logger in your application
        fx.Invoke(runApp),
    )

    app.Run()
}

func runApp(log logger.Logger) {
    log.Info("Application started with dependency injection")
    log.Infof("Running in %s mode", "production")
}
```

## Structured Logging

Use the `*w` methods (e.g., `Infow`, `Errorw`) for structured logging with fields:

```go
log.Infow("User logged in", logger.Fields{
    "user_id": "user-123",
    "ip": "192.168.1.1",
    "timestamp": time.Now(),
})
```

## Formatted Logging

Use the `*f` methods for formatted strings:

```go
log.Infof("User %s logged in from %s", userID, ipAddress)
```

## Error Logging

Special methods for logging errors:

```go
// Log error with message
err := processData()
if err != nil {
    log.Err("Failed to process data", err)
}

// Log error with structured fields
log.Errorw("Database query failed", logger.Fields{
    "query": sqlQuery,
    "error": err.Error(),
})
```

## Environment Configuration

The default logger respects the `LogConfig_LogType` environment variable:

```bash
# Use Zap logger (default)
export LogConfig_LogType="Zap"

# Run your application
go run main.go
```

## Complete Example

```go
package main

import (
    "fmt"
    defaultlogger "github.com/phatnt199/go-infra/pkg/logger/default_logger"
    "github.com/phatnt199/go-infra/pkg/logger"
)

func main() {
    log := defaultlogger.GetLogger()

    // Simple message
    log.Info("Application starting...")

    // Formatted message
    port := 8080
    log.Infof("Server listening on port %d", port)

    // Structured logging
    log.Infow("Request received", logger.Fields{
        "method": "GET",
        "path": "/api/users",
        "status": 200,
    })

    // Error logging
    err := someOperation()
    if err != nil {
        log.Err("Operation failed", err)
    }

    log.Info("Application shutting down")
}

func someOperation() error {
    return fmt.Errorf("simulated error")
}
```

## Log Levels

The logger supports standard log levels:

- **Debug** - Detailed information for diagnosing problems
- **Info** - General informational messages
- **Warn** - Warning messages for potential issues
- **Error** - Error events
- **Fatal** - Severe errors that cause application exit

```go
log.Debug("Detailed debug information")
log.Info("Normal flow of application")
log.Warn("Something unexpected happened")
log.Error("An error occurred")
log.Fatal("Critical error - application will exit")
```

## Next Steps

- **[Configuration](./configuration.md)** - Learn how to configure the logger
- **[Usage Guide](./usage.md)** - Explore all logging methods
- **[Fx Integration](./fx-integration.md)** - Use with Uber Fx
- **[Best Practices](./best-practices.md)** - Production-ready patterns
