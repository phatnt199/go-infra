---
sidebar_position: 6
---

# External Adapters

Adapters for integrating go-infra logger with external libraries.

## Overview

The logger package provides adapters to integrate with popular libraries:

- **GORM Logger** - Database query logging
- **Fx Logger** - Uber Fx lifecycle event logging

## GORM Logger Adapter

Integrate go-infra logger with GORM for database query logging.

### Setup

```go
import (
    "github.com/phatnt199/go-infra/pkg/logger"
    gormlog "github.com/phatnt199/go-infra/pkg/logger/external/gormlog"
    "gorm.io/gorm"
)

func setupDatabase(log logger.Logger) *gorm.DB {
    // Create GORM logger adapter
    gormLogger := gormlog.NewGormCustomLogger(log)

    // Use with GORM
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: gormLogger,
    })

    return db
}
```

### Complete Example

```go
package main

import (
    "go.uber.org/fx"
    defaultlogger "github.com/phatnt199/go-infra/pkg/logger/default_logger"
    gormlog "github.com/phatnt199/go-infra/pkg/logger/external/gormlog"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    log := defaultlogger.GetLogger()

    // Create GORM logger
    gormLogger := gormlog.NewGormCustomLogger(log)

    // Open database with logger
    dsn := "host=localhost user=postgres password=postgres dbname=mydb"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: gormLogger,
    })
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Use database - queries will be logged
    var users []User
    db.Find(&users)
}
```

### With Fx

```go
package main

import (
    "go.uber.org/fx"
    "github.com/phatnt199/go-infra/pkg/logger"
    gormlog "github.com/phatnt199/go-infra/pkg/logger/external/gormlog"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    app := fx.New(
        zaplogger.Module,

        fx.Provide(provideDatabase),

        fx.Invoke(runApp),
    )

    app.Run()
}

func provideDatabase(log logger.Logger) (*gorm.DB, error) {
    gormLogger := gormlog.NewGormCustomLogger(log)

    dsn := "host=localhost user=postgres password=postgres dbname=mydb"
    return gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: gormLogger,
    })
}

func runApp(db *gorm.DB, log logger.Logger) {
    log.Info("Database initialized with GORM logger")
}
```

### Log Levels

The GORM logger adapter supports GORM's log levels:

```go
import gormlogger "gorm.io/gorm/logger"

gormLogger := gormlog.NewGormCustomLogger(log)

// Set log level
gormLogger = gormLogger.LogMode(gormlogger.Info)  // Log all queries
gormLogger = gormLogger.LogMode(gormlogger.Warn)  // Log slow queries
gormLogger = gormLogger.LogMode(gormlogger.Error) // Log errors only
gormLogger = gormLogger.LogMode(gormlogger.Silent) // No logging
```

### GORM Logger Methods

The adapter implements the GORM logger interface:

```go
type GormCustomLogger struct {
    logger.Logger
    gormlogger.Config
}

// LogMode - Set log mode
func (l *GormCustomLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface

// Info - Log info messages
func (l GormCustomLogger) Info(ctx context.Context, str string, args ...interface{})

// Warn - Log warnings
func (l GormCustomLogger) Warn(ctx context.Context, str string, args ...interface{})

// Error - Log errors
func (l GormCustomLogger) Error(ctx context.Context, str string, args ...interface{})

// Trace - Log SQL queries
func (l GormCustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error)
```

### Query Logging Output

Example of logged SQL queries:

```
DEBUG | [150 ms, 1 rows] sql -> SELECT * FROM users WHERE id = $1
DEBUG | [50 ms, 10 rows] sql -> SELECT * FROM orders WHERE user_id = $1
WARN  | [2500 ms, 0 rows] sql -> SELECT * FROM products WHERE category = $1 -- Slow query
ERROR | [0 ms, 0 rows] sql -> INSERT INTO users VALUES (...) -- Duplicate key error
```

## Fx Logger Adapter

Integrate go-infra logger with Uber Fx lifecycle events.

### Setup

```go
import (
    "go.uber.org/fx"
    "github.com/phatnt199/go-infra/pkg/logger"
    fxlog "github.com/phatnt199/go-infra/pkg/logger/external/fxlog"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

func main() {
    app := fx.New(
        zaplogger.Module,

        // Add Fx logger
        fx.WithLogger(func(log logger.Logger) fxevent.Logger {
            return fxlog.NewCustomFxLogger(log)
        }),

        fx.Invoke(runApp),
    )

    app.Run()
}
```

### Using FxLogger Option

Convenience option for quick setup:

```go
import (
    "go.uber.org/fx"
    fxlog "github.com/phatnt199/go-infra/pkg/logger/external/fxlog"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

func main() {
    app := fx.New(
        zaplogger.Module,
        fxlog.FxLogger, // One-liner setup

        fx.Invoke(runApp),
    )

    app.Run()
}
```

### Logged Fx Events

The Fx logger adapter logs various lifecycle events:

#### Dependency Provision

```
DEBUG | provided | constructor=NewUserRepository | type=*repository.UserRepository | module=app
DEBUG | provided | constructor=NewUserService | type=*service.UserService | module=app
```

#### Function Invocation

```
DEBUG | invoking | function=setupRoutes | module=app
DEBUG | invoking | function=runApp | module=app
```

#### OnStart Hooks

```
DEBUG | OnStart hook executing | caller=server.Start | function=(*Server).Start
INFO  | OnStart hook executed | caller=server.Start | function=(*Server).Start | runtime=150ms
```

#### OnStop Hooks

```
DEBUG | OnStop hook executing | caller=server.Stop | function=(*Server).Stop
DEBUG | OnStop hook executed | caller=server.Stop | function=(*Server).Stop | runtime=50ms
```

#### Errors

```
ERROR | error encountered while applying options | module=app | error=missing type: *config.Config
ERROR | invoke failed | function=runApp | error=failed to build container | module=app
```

#### Application Lifecycle

```
DEBUG | started
DEBUG | received signal | signal=INTERRUPT
DEBUG | stop
```

### Custom Error Formatting

The Fx logger includes custom error formatting for better readability:

```go
// Before
ERROR | invoke failed | error=missing type: *config.Config
      caused by: missing provider for *config.Config
      in provider NewService from example.com/app
      in call stack: NewService -> NewApp -> main

// After
ERROR | invoke failed | error=missing type: *config.Config | caused by: missing provider | in provider NewService
```

### Complete Example with Fx Logging

```go
package main

import (
    "context"
    "go.uber.org/fx"
    "github.com/phatnt199/go-infra/pkg/logger"
    fxlog "github.com/phatnt199/go-infra/pkg/logger/external/fxlog"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
)

func main() {
    app := fx.New(
        // Logger with Fx event logging
        zaplogger.Module,
        fxlog.FxLogger,

        // Application components
        fx.Provide(
            NewServer,
            NewDatabase,
        ),

        fx.Invoke(runApp),
    )

    app.Run()
}

type Server struct {
    log logger.Logger
}

func NewServer(log logger.Logger, lc fx.Lifecycle) *Server {
    s := &Server{log: log}

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            s.log.Info("Starting server")
            // Start logic
            return nil
        },
        OnStop: func(ctx context.Context) error {
            s.log.Info("Stopping server")
            // Stop logic
            return nil
        },
    })

    return s
}

type Database struct {
    log logger.Logger
}

func NewDatabase(log logger.Logger, lc fx.Lifecycle) *Database {
    db := &Database{log: log}

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            db.log.Info("Connecting to database")
            return nil
        },
        OnStop: func(ctx context.Context) error {
            db.log.Info("Closing database connection")
            return nil
        },
    })

    return db
}

func runApp(server *Server, db *Database) {
    server.log.Info("Application running")
}
```

Output:

```
DEBUG | provided | constructor=NewServer | type=*main.Server
DEBUG | provided | constructor=NewDatabase | type=*main.Database
DEBUG | invoking | function=runApp
DEBUG | OnStart hook executing | caller=NewServer | function=OnStart
INFO  | Starting server
INFO  | OnStart hook executed | caller=NewServer | runtime=10ms
DEBUG | OnStart hook executing | caller=NewDatabase | function=OnStart
INFO  | Connecting to database
INFO  | OnStart hook executed | caller=NewDatabase | runtime=100ms
DEBUG | started
INFO  | Application running
```

## Using Both Adapters Together

Combine GORM and Fx logger adapters:

```go
package main

import (
    "go.uber.org/fx"
    "github.com/phatnt199/go-infra/pkg/logger"
    fxlog "github.com/phatnt199/go-infra/pkg/logger/external/fxlog"
    gormlog "github.com/phatnt199/go-infra/pkg/logger/external/gormlog"
    zaplogger "github.com/phatnt199/go-infra/pkg/logger/zap"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func main() {
    app := fx.New(
        // Logger with Fx logging
        zaplogger.Module,
        fxlog.FxLogger,

        // Database with GORM logging
        fx.Provide(provideDatabase),

        fx.Invoke(runApp),
    )

    app.Run()
}

func provideDatabase(log logger.Logger) (*gorm.DB, error) {
    gormLogger := gormlog.NewGormCustomLogger(log)

    dsn := "host=localhost user=postgres password=postgres dbname=mydb"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: gormLogger,
    })

    if err != nil {
        return nil, err
    }

    return db, nil
}

func runApp(db *gorm.DB, log logger.Logger) {
    log.Info("Application running with GORM and Fx logging")

    // Database queries will be logged via GORM adapter
    var users []User
    db.Find(&users)
}
```

## Next Steps

- **[Best Practices](./best-practices.md)** - Production patterns
- **[Usage Guide](./usage.md)** - Learn all logging methods
- **[Configuration](./configuration.md)** - Configure logger
