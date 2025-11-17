---
sidebar_position: 1
---

# Logger

High-performance structured logging for Go applications built on Uber's Zap logger.

## Overview

The go-infra logger package provides a flexible logging abstraction with multiple implementations:

- **Zap Logger** - Production-ready, high-performance structured logger
- **Empty Logger** - No-op logger for testing or when logging is not needed
- **Default Logger** - Quick-start logger that auto-initializes with sensible defaults

The logger uses an interface-based design that integrates seamlessly with Uber's Fx dependency injection framework.

## Key Features

- ✅ **Interface-based design** - Easy to mock and test
- ✅ **High performance** - Built on Uber's Zap
- ✅ **Structured logging** - Type-safe fields with auto-complete
- ✅ **Fx integration** - First-class support for dependency injection
- ✅ **Configurable** - Environment-based configuration
- ✅ **OpenTelemetry support** - Optional distributed tracing integration
- ✅ **External adapters** - GORM and Fx logger adapters included

## Quick Start

### Using Default Logger

For simple applications or quick prototyping:

```go
package main

import (
    defaultlogger "github.com/phatnt199/go-infra/pkg/logger/default_logger"
)

func main() {
    log := defaultlogger.GetLogger()

    log.Info("Application started")
    log.Infow("User logged in", map[string]interface{}{
        "user_id": "12345",
        "ip": "192.168.1.1",
    })
}
```

### Using with Fx Dependency Injection (Recommended)

For production applications:

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

        // Use logger in your components
        fx.Invoke(run),
    )

    app.Run()
}

func run(log logger.Logger) {
    log.Info("Application started with dependency injection")
}
```

## Documentation Structure

- **[Quick Start](./quickstart.md)** - Get started quickly
- **[Configuration](./configuration.md)** - Configure logger for different environments
- **[Usage Guide](./usage.md)** - Learn how to use the logger effectively
- **[Fx Integration](./fx-integration.md)** - Use logger with Uber Fx
- **[External Adapters](./adapters.md)** - GORM and Fx logger adapters
- **[Best Practices](./best-practices.md)** - Production-ready patterns

## Logger Interface

All logger implementations conform to the `logger.Logger` interface:

```go
type Logger interface {
    Configure(cfg func(internalLog interface{}))

    // Debug level
    Debug(args ...interface{})
    Debugf(template string, args ...interface{})
    Debugw(msg string, fields Fields)

    // Info level
    Info(args ...interface{})
    Infof(template string, args ...interface{})
    Infow(msg string, fields Fields)

    // Warn level
    Warn(args ...interface{})
    Warnf(template string, args ...interface{})
    WarnMsg(msg string, err error)

    // Error level
    Error(args ...interface{})
    Errorw(msg string, fields Fields)
    Errorf(template string, args ...interface{})
    Err(msg string, err error)

    // Fatal level (exits application)
    Fatal(args ...interface{})
    Fatalf(template string, args ...interface{})

    // Utility
    Printf(template string, args ...interface{})
    WithName(name string)
    LogType() models.LogType

    // gRPC middleware logging
    GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error)
    GrpcClientInterceptorLogger(method string, req, reply interface{}, time time.Duration, metaData map[string][]string, err error)
}
```

## Logger Types

### Zap Logger

Production-ready logger with:

- High performance (4+ million logs/second)
- Structured logging with type-safe fields
- JSON encoding for production
- Console encoding for development
- OpenTelemetry tracing support
- Configurable caller information and stack traces

### Empty Logger

No-op logger for:

- Testing without log noise
- Benchmarking without I/O overhead
- Temporary logging disablement

### Default Logger

Auto-initializing logger for:

- Quick prototyping
- Simple applications
- Migration utilities
- Command-line tools

## Next Steps

Choose your path:

- **New to go-infra?** Start with [Quick Start](./quickstart.md)
- **Production deployment?** See [Configuration](./configuration.md) and [Best Practices](./best-practices.md)
- **Using Fx?** Read [Fx Integration](./fx-integration.md)
- **Need examples?** Check [Usage Guide](./usage.md)
