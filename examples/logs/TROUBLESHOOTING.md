# Logger - Troubleshooting & FAQ

Common issues, solutions, and frequently asked questions about using the logger in go-infra.

## Table of Contents

- [Common Issues](#common-issues)
- [Frequently Asked Questions](#frequently-asked-questions)
- [Performance Tips](#performance-tips)
- [Integration Guides](#integration-guides)

---

## Common Issues

### Issue 1: Logger Not Showing Any Output

**Problem:** You've added logger calls but no logs appear in the console.

**Causes:**

1. Logger not initialized properly
2. Log level too high (e.g., set to `fatal` when logging `info`)
3. Logger not injected via dependency injection
4. Application exits before logs are flushed

**Solutions:**

```go
// ✅ CORRECT - Always initialize and sync
func main() {
    logOptions := &config.LogOptions{
        LogLevel:      "debug",  // Lower level to see more logs
        LogType:       models.Zap,
        CallerEnabled: true,
    }

    env := environment.NewEnvironment("development")
    appLogger := zap.NewZapLogger(logOptions, env)
    defer appLogger.Sync()  // IMPORTANT: Flush logs on exit

    appLogger.Info("Application started")
}

// ❌ WRONG - Logs not flushed
func main() {
    appLogger := zap.NewZapLogger(logOptions, env)
    // Missing: defer appLogger.Sync()
    appLogger.Info("This might not appear!")
}
```

**Check log level:**

```bash
# View current log level (check environment or config)
echo $LogConfig_LogLevel

# Set to debug to see all logs
export LogConfig_LogLevel=debug
```

---

### Issue 2: "Logger is Nil" Error

**Problem:** Panic with `nil pointer dereference` when calling logger methods.

**Causes:**

1. Logger not injected via fx.Provide
2. Constructor parameter name mismatch
3. Missing dependency in fx options

**Solutions:**

```go
// ✅ CORRECT - Logger injected via fx
type MyService struct {
    logger logger.Logger  // Correct type
}

func NewMyService(logger logger.Logger) *MyService {  // Parameter name doesn't matter
    return &MyService{logger: logger}
}

// In fx.Options:
fx.Provide(NewMyService)

// ❌ WRONG - Logger not injected
type MyService struct {
    logger logger.Logger
}

func NewMyService() *MyService {  // Logger parameter missing!
    return &MyService{logger: nil}  // panic!
}
```

**Verify fx setup:**

```go
// Make sure logger is provided in fx.Options
app := fx.New(
    fx.Provide(func() logger.Logger {
        logOptions := &config.LogOptions{LogLevel: "info"}
        return zap.NewZapLogger(logOptions, env)
    }),
    fx.Provide(NewMyService),  // Depends on logger
    fx.Invoke(startApp),
)
```

---

### Issue 3: Logs Not Showing Request Details

**Problem:** Logs don't include important context like user ID or request ID.

**Causes:**

1. Using simple logging instead of structured fields
2. Not passing context through application
3. Creating separate logger instances instead of using fields

**Solutions:**

```go
// ❌ WRONG - No context, hard to parse
logger.Infof("User login from %s in %d ms", ipAddress, duration)
// Output: "User login from 192.168.1.1 in 250 ms"

// ✅ CORRECT - Structured fields, easy to parse
logger.Infow("User login successful",
    logger.String("user_id", userID),
    logger.String("ip_address", ipAddress),
    logger.Int64("duration_ms", duration),
)
// Output: {"level":"info","user_id":"usr-123","ip_address":"192.168.1.1","duration_ms":250}
```

---

### Issue 4: Performance Degradation with Logging

**Problem:** Application becomes slow after adding logging.

**Causes:**

1. Logging in tight loops
2. Using `Any()` field type with complex objects
3. Logging at too low a level (debug in production)
4. Not using defer for cleanup

**Solutions:**

```go
// ❌ WRONG - Logs in tight loop with reflection
for i := 0; i < 1000000; i++ {
    logger.Infow("Processing item", logger.Any("item", complexObject))
}

// ✅ CORRECT - Log periodically, use specific types
for i := 0; i < 1000000; i++ {
    if i % 10000 == 0 {
        logger.Infof("Processed %d items", i)
    }
}

// ❌ WRONG - Using Any() for everything
logger.Infow("User data", logger.Any("user", user))  // Slow

// ✅ CORRECT - Use specific field types
logger.Infow("User login",
    logger.String("user_id", user.ID),
    logger.String("email", user.Email),
    logger.Bool("verified", user.Verified),
)
```

---

### Issue 5: Duplicate Log Entries

**Problem:** Each log message appears multiple times.

**Causes:**

1. Logger registered multiple times in fx.Provide
2. Multiple Sync() calls
3. Logger error handlers creating recursive logs

**Solutions:**

```go
// ❌ WRONG - Logger provided twice
fx.New(
    fx.Provide(CreateLogger),  // First provider
    fx.Provide(CreateLogger),  // Duplicate!
)

// ✅ CORRECT - Provide once
fx.New(
    fx.Provide(CreateLogger),
    fx.Provide(NewMyService),
)

// Check for duplicate Sync() calls
defer logger.Sync()  // Call only once in main()
// Don't call in each function
```

---

### Issue 6: Lost Logs on Application Panic

**Problem:** When application panics, recent logs are lost.

**Causes:**

1. `defer logger.Sync()` not called before panic
2. Buffered logs not flushed before OS.Exit

**Solutions:**

```go
// ✅ CORRECT - Ensure Sync is always called
func main() {
    logger := initializeLogger()
    defer logger.Sync()  // Will run even if panic occurs

    defer func() {
        if r := recover(); r != nil {
            logger.Panicf("Application panicked: %v", r)
            logger.Sync()  // Extra flush before panic propagates
            panic(r)
        }
    }()

    startApplication()
}
```

---

### Issue 7: JSON Logs Not Readable

**Problem:** Production logs are JSON format but hard to debug.

**Causes:**

1. APP_ENV set to production
2. Encoding set to JSON in config

**Solutions:**

```go
// For development debugging:
export APP_ENV=development
export LogConfig_LogLevel=debug

// If you need to parse JSON logs:
# Using jq to format JSON logs
tail -f app.log | jq '.'

# Extract specific fields
jq '.user_id, .duration' app.log

# Filter by level
jq 'select(.level=="error")' app.log
```

---

### Issue 8: Caller Information Shows Wrong File

**Problem:** Log shows wrong file:line number, or no caller info at all.

**Causes:**

1. CallerEnabled set to false
2. Wrong skip count for wrapped loggers
3. Logging through multiple abstraction layers

**Solutions:**

```go
// ✅ Enable caller information
logOptions := &config.LogOptions{
    CallerEnabled: true,  // Enable
}

// For wrapped loggers, adjust skip count:
// zap.AddCallerSkip(N) where N is number of wrapper layers
logOptions := &config.LogOptions{
    LogLevel:      "info",
    CallerEnabled: true,
    // CallerSkip adjusted in zap_logger.go if needed
}
```

---

## Frequently Asked Questions

### Q1: Can I use the logger from a package/library?

**A:** Yes! Services and repositories should accept `logger.Logger` as a constructor parameter.

```go
package mylib

import "github.com/phatnt199/go-infra/pkg/logger"

type MyLibrary struct {
    logger logger.Logger
}

func NewMyLibrary(logger logger.Logger) *MyLibrary {
    return &MyLibrary{logger: logger}
}

func (m *MyLibrary) DoWork() {
    m.logger.Info("Doing work")
}
```

---

### Q2: Should I create a logger instance in each file?

**A:** No! Always use dependency injection. Create logger once and pass to services.

```go
// ❌ WRONG - Don't do this
package userservice

import "github.com/phatnt199/go-infra/pkg/logger"

var logger = createLogger()  // Global logger

// ✅ CORRECT - Inject via constructor
package userservice

import "github.com/phatnt199/go-infra/pkg/logger"

type UserService struct {
    logger logger.Logger
}

func NewUserService(logger logger.Logger) *UserService {
    return &UserService{logger: logger}
}
```

---

### Q3: How do I add request/trace ID to all logs in a request?

**A:** Add fields when creating a scoped logger for the request.

```go
// In middleware or handler
requestID := r.Header.Get("X-Request-ID")
if requestID == "" {
    requestID = generateRequestID()
}

// Create a logger with request context
requestLogger := logger.WithFields(
    logger.String("request_id", requestID),
    logger.String("method", r.Method),
    logger.String("path", r.URL.Path),
)

// All logs through requestLogger will include these fields
requestLogger.Info("Processing request")  // Includes request_id
```

---

### Q4: Can I use different log levels for different packages?

**A:** Currently, log level is global. For per-package control, use Debug/Info strategically:

```go
// Development
export LogConfig_LogLevel=debug   # See everything

// Production
export LogConfig_LogLevel=info    # Only important events

// Inside your code
logger.Debug("Detailed diagnostic")    // Only in debug level
logger.Info("Important event")         // Always shown
```

---

### Q5: How do I handle sensitive data in logs?

**A:** Never log passwords, tokens, or API keys. Log only safe identifiers.

```go
// ❌ WRONG - Exposes password
logger.Infof("User login with password: %s", password)

// ✅ CORRECT - Log only safe info
logger.Infow("User authentication",
    logger.String("user_id", userID),
    logger.String("auth_method", "password"),
)

// ✅ CORRECT - Mask sensitive data if needed
maskedToken := token[:10] + "***"
logger.Infof("Token provided: %s", maskedToken)
```

---

### Q6: Should I catch and log errors or let them bubble up?

**A:** Log once at the appropriate boundary (handler, middleware, service level), then return the error.

```go
// ❌ WRONG - Logs error multiple times
func (s *Service) GetUser(id string) (*User, error) {
    user, err := s.repo.Find(id)
    if err != nil {
        logger.Error("User not found", err)  // First log
        return nil, err
    }
}

// Handler logs again
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.service.GetUser(id)
    if err != nil {
        logger.Error("Failed to get user", err)  // Second log - WRONG
    }
}

// ✅ CORRECT - Log once at handler level
func (s *Service) GetUser(id string) (*User, error) {
    user, err := s.repo.Find(id)
    if err != nil {
        return nil, err  // Don't log here
    }
    return user, nil
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.service.GetUser(id)
    if err != nil {
        logger.Errorw("Failed to get user",  // Log once
            logger.String("user_id", id),
            logger.Err(err),
        )
        http.Error(w, "User not found", 404)
    }
}
```

---

### Q7: How do I test code that uses logging?

**A:** Create a test logger with debug level enabled.

```go
func TestMyService(t *testing.T) {
    // Create test logger
    logOptions := &config.LogOptions{
        LogLevel: "debug",
        LogType:  models.Zap,
    }
    testLogger := zap.NewZapLogger(logOptions, env)
    defer testLogger.Sync()

    // Use in test
    service := NewMyService(testLogger)
    result, err := service.DoSomething()

    // Logger will output debug info during test
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
}
```

---

### Q8: Can I send logs to external services?

**A:** Configure OutputPaths in the logger configuration to send to files or services.

```go
logOptions := &config.LogOptions{
    OutputPaths: []string{
        "stdout",                    // Console
        "/var/log/app/app.log",     // File
        "syslog://localhost:514",   // Syslog
    },
    ErrorOutputPaths: []string{
        "stderr",
        "/var/log/app/error.log",
    },
}
```

---

### Q9: What's the performance overhead of structured logging?

**A:** Minimal with proper usage. Zap is designed for high performance:

- Structured logging: ~1-2% overhead vs printf
- Field allocation: O(1) for small number of fields
- Sync is blocking: ~100μs per call
- In hot paths: Use field types instead of `Any()`

---

### Q10: Why use go-infra logger instead of Zap directly?

**A:** Several advantages:

1. **Clean API**: Single import, no zap imports in user code
2. **Consistency**: All services follow same pattern
3. **Flexibility**: Can switch underlying implementation without breaking user code
4. **Maintainability**: Logger initialization in one place
5. **Type Safety**: All field types are provided by go-infra

```go
// ✅ WITH go-infra
import "github.com/phatnt199/go-infra/pkg/logger"
logger.Info("message", logger.String("key", "value"))

// ❌ WITHOUT go-infra (Zap directly)
import "go.uber.org/zap"
import "go.uber.org/zap/zapcore"
zapLogger.Info("message", zap.String("key", "value"))  // More verbose
```

---

## Performance Tips

### Tip 1: Use Debug Level for Development Only

```bash
# Development
export APP_ENV=development
export LogConfig_LogLevel=debug

# Production
export APP_ENV=production
export LogConfig_LogLevel=info
```

### Tip 2: Avoid Reflection in Hot Paths

```go
// ❌ SLOW - Uses reflection
for i := 0; i < 1000000; i++ {
    logger.Debugf("Item: %v", complexObject)  // Reflection
}

// ✅ FAST - Pre-allocated fields
for i := 0; i < 1000000; i++ {
    if i % 10000 == 0 {
        logger.Debugf("Processed %d items", i)  // No reflection
    }
}
```

### Tip 3: Disable Caller Info in Production if Needed

```go
// Caller info has small overhead (~10-50ns per call)
logOptions := &config.LogOptions{
    CallerEnabled: true,  // Set false in high-throughput scenarios
}
```

### Tip 4: Batch Write Large Amounts of Data

```go
// ❌ WRONG - 1000 log calls
for _, item := range items {
    logger.Debugf("Item: %v", item)
}

// ✅ CORRECT - Batch log
logger.Debugf("Processing %d items", len(items))
// ... process ...
if hasErrors {
    logger.Warnf("Completed with %d errors", errorCount)
}
```

---

## Integration Guides

### Integration with HTTP Server

```go
// main.go
func main() {
    // Initialize logger
    logOptions := &config.LogOptions{LogLevel: "info"}
    appLogger := zap.NewZapLogger(logOptions, env)
    defer appLogger.Sync()

    // Create router with logging middleware
    mux := http.NewServeMux()

    // Add logging middleware
    handler := LoggingMiddleware(appLogger)(mux)

    server := &http.Server{
        Addr:    ":8080",
        Handler: handler,
    }

    appLogger.Info("Server starting on :8080")
    if err := server.ListenAndServe(); err != nil {
        appLogger.Fatal("Server error", logger.Err(err))
    }
}
```

### Integration with Database Queries

```go
// repository.go
type UserRepository struct {
    logger logger.Logger
    db     *gorm.DB
}

func (r *UserRepository) Create(user *User) error {
    r.logger.Debugf("Creating user: %s", user.Email)

    if err := r.db.Create(user).Error; err != nil {
        r.logger.Errorw("Create user failed",
            logger.String("email", user.Email),
            logger.Err(err),
        )
        return err
    }

    r.logger.Infof("User created: %s", user.ID)
    return nil
}
```

### Integration with Background Jobs

```go
// worker.go
type BackgroundWorker struct {
    logger logger.Logger
}

func (w *BackgroundWorker) Process(job Job) {
    start := time.Now()

    w.logger.Infow("Starting job",
        logger.String("job_id", job.ID),
        logger.String("job_type", job.Type),
    )

    if err := job.Execute(); err != nil {
        w.logger.Errorw("Job failed",
            logger.String("job_id", job.ID),
            logger.Duration("duration", time.Since(start)),
            logger.Err(err),
        )
        return
    }

    w.logger.Infow("Job completed",
        logger.String("job_id", job.ID),
        logger.Duration("duration", time.Since(start)),
    )
}
```

---

**Still have questions? Check the main documentation or raise an issue!** 📚
