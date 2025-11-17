---
sidebar_position: 3
---

# Configuration

Configure the logger for different environments and use cases.

## Configuration Options

The `LogOptions` struct provides configuration for the logger:

```go
type LogOptions struct {
    LogLevel      string         // Log level: "debug", "info", "warn", "error", "fatal"
    LogType       models.LogType // Logger implementation (Zap or Logrus)
    CallerEnabled bool           // Include file and line number
    EnableTracing bool           // Enable OpenTelemetry tracing (default: true)
}
```

## Creating Configuration

### Manual Configuration

```go
import (
    "github.com/phatnt199/go-infra/pkg/logger/config"
    "github.com/phatnt199/go-infra/pkg/logger/models"
)

cfg := &config.LogOptions{
    LogLevel:      "info",
    LogType:       models.Zap,
    CallerEnabled: true,
    EnableTracing: true,
}
```

### Configuration from Environment (with Fx)

Use `ProvideLogConfig` to load configuration from environment variables:

```go
import (
    "go.uber.org/fx"
    "github.com/phatnt199/go-infra/pkg/logger/config"
    "github.com/phatnt199/go-infra/pkg/application/environment"
)

app := fx.New(
    fx.Provide(config.ProvideLogConfig),
    fx.Invoke(func(cfg *config.LogOptions) {
        // Configuration loaded from environment
        fmt.Printf("Log Level: %s\n", cfg.LogLevel)
    }),
)
```

## Environment Variables

The logger configuration can be controlled via environment variables using the naming convention:

```bash
# Log level
export LOG_OPTIONS_LEVEL="info"

# Log type (0 = Zap, 1 = Logrus)
export LOG_OPTIONS_LOG_TYPE="0"

# Enable caller information
export LOG_OPTIONS_CALLER_ENABLED="true"

# Enable tracing
export LOG_OPTIONS_ENABLE_TRACING="true"
```

### Default Logger Environment Variables

The default logger specifically looks for:

```bash
# Set logger type for default logger
export LogConfig_LogType="Zap"
```

## Log Levels

Available log levels in order of severity:

| Level   | Description                     | Use Case                                      |
| ------- | ------------------------------- | --------------------------------------------- |
| `debug` | Detailed diagnostic information | Development and troubleshooting               |
| `info`  | General informational messages  | Normal application flow                       |
| `warn`  | Warning messages                | Potential issues that don't prevent operation |
| `error` | Error events                    | Errors that need attention                    |
| `fatal` | Severe errors causing exit      | Critical failures                             |

```go
cfg := &config.LogOptions{
    LogLevel: "info", // Set minimum log level
}
```

:::tip
In production, use `info` or higher to reduce log volume. In development, use `debug` for detailed information.
:::

## Environment-Based Configuration

The logger adapts its output format based on the environment:

### Development Environment

```go
import "github.com/phatnt199/go-infra/pkg/application/constants"

log := zaplogger.NewZapLogger(cfg, constants.DEV_ENV)
```

Features:

- **Console encoding** - Human-readable format
- **Colored output** - Visual level indicators
- **Full caller paths** - Complete file paths
- **Readable timestamps** - ISO8601 format
- **Console separator** - `|` for readability

Example output:

```
2024-01-01T10:00:00Z | INFO | main.go:42 | Application started | user_id=123
```

### Production Environment

```go
import "github.com/phatnt199/go-infra/pkg/application/constants"

log := zaplogger.NewZapLogger(cfg, constants.PROD_ENV)
```

Features:

- **JSON encoding** - Machine-parseable format
- **Structured fields** - Easy log aggregation
- **Short caller paths** - Relative paths only
- **Optimized performance** - Minimal overhead

Example output:

```json
{
	"level": "info",
	"time": "2024-01-01T10:00:00Z",
	"caller": "main.go:42",
	"message": "Application started",
	"user_id": "123"
}
```

## Caller Information

Enable/disable file and line number in logs:

```go
cfg := &config.LogOptions{
    CallerEnabled: true, // Show file:line
}
```

With `CallerEnabled: true`:

```
2024-01-01T10:00:00Z | INFO | main.go:42 | User logged in
```

With `CallerEnabled: false`:

```
2024-01-01T10:00:00Z | INFO | User logged in
```

:::caution
Enabling caller information adds a small performance overhead. Disable in high-throughput scenarios if not needed.
:::

## OpenTelemetry Tracing

Enable distributed tracing integration:

```go
cfg := &config.LogOptions{
    EnableTracing: true, // Add logs as events to traces
}
```

When enabled:

- Logs are automatically added as events to active OpenTelemetry spans
- Correlate logs with distributed traces
- Enhanced observability in microservices

See [OpenTelemetry Integration](../../observability/tracing.md) for more details.

## Complete Configuration Examples

### Development Setup

```go
package main

import (
    "github.com/phatnt199/go-infra/pkg/application/constants"
    "github.com/phatnt199/go-infra/pkg/logger/config"
    "github.com/phatnt199/go-infra/pkg/logger/models"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

func main() {
    cfg := &config.LogOptions{
        LogLevel:      "debug",      // Show all logs
        LogType:       models.Zap,
        CallerEnabled: true,         // Helpful for debugging
        EnableTracing: false,        // Usually not needed in dev
    }

    log := zaplogger.NewZapLogger(cfg, constants.DEV_ENV)
    log.Debug("Development logger initialized")
}
```

### Production Setup

```go
package main

import (
    "github.com/phatnt199/go-infra/pkg/application/constants"
    "github.com/phatnt199/go-infra/pkg/logger/config"
    "github.com/phatnt199/go-infra/pkg/logger/models"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

func main() {
    cfg := &config.LogOptions{
        LogLevel:      "info",       // Reduce log volume
        LogType:       models.Zap,
        CallerEnabled: false,        // Optimize performance
        EnableTracing: true,         // Enable for observability
    }

    log := zaplogger.NewZapLogger(cfg, constants.PROD_ENV)
    log.Info("Production logger initialized")
}
```

### Configuration with Fx

```go
package main

import (
    "go.uber.org/fx"
    "github.com/phatnt199/go-infra/pkg/logger"
    "github.com/phatnt199/go-infra/pkg/logger/config"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

func main() {
    app := fx.New(
        // Zap module includes config provider
        zaplogger.Module,

        fx.Invoke(runWithConfig),
    )
    app.Run()
}

func runWithConfig(log logger.Logger, cfg *config.LogOptions) {
    log.Infof("Logger configured with level: %s", cfg.LogLevel)
}
```

### Configuration File (JSON/YAML)

If using go-infra's config system, create a config file:

```json
{
	"logOptions": {
		"level": "info",
		"logType": 0,
		"callerEnabled": true,
		"enableTracing": true
	}
}
```

Load it with:

```go
import (
    "github.com/phatnt199/go-infra/pkg/application/environment"
    "github.com/phatnt199/go-infra/pkg/logger/config"
)

// Assuming you have environment setup
cfg, err := config.ProvideLogConfig(env)
if err != nil {
    panic(err)
}
```

## Logger Types

### Zap (models.Zap = 0)

High-performance structured logger:

```go
cfg.LogType = models.Zap // or just 0
```

### Logrus (models.Logrus = 1)

Alternative logger implementation:

```go
cfg.LogType = models.Logrus // or just 1
```

:::note
Currently, only Zap logger is fully implemented in the codebase. Logrus support is defined in the enum but not implemented.
:::

## Access Internal Logger

For advanced use cases, access the underlying Zap logger:

```go
zapLog := log.(zaplogger.ZapLogger)
internalZap := zapLog.InternalLogger()

// Use Zap-specific features
internalZap.With(zap.String("key", "value"))
```

## Next Steps

- **[Usage Guide](./usage.md)** - Learn all logging methods
- **[Fx Integration](./fx-integration.md)** - Use with dependency injection
- **[Best Practices](./best-practices.md)** - Production patterns
