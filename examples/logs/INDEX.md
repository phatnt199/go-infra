# Logger Examples - Index

This directory contains comprehensive documentation and examples for using the logging system in the go-infra project.

## 📂 Files in This Directory

| File                   | Purpose                                  | Read Time |
| ---------------------- | ---------------------------------------- | --------- |
| **README.md**          | Start here! Overview and quick reference | 5 min     |
| **HOW_TO_USE_LOGS.md** | Complete guide with all features         | 20 min    |
| **CODE_EXAMPLES.md**   | Copy-paste ready code examples           | 15 min    |
| **TROUBLESHOOTING.md** | Common issues & FAQ                      | 10 min    |

---

## 🎯 Start Here

### For First-Time Users

1. Read: [README.md](./README.md) - Quick overview
2. Read: [HOW_TO_USE_LOGS.md - Quick Start](./HOW_TO_USE_LOGS.md#quick-start)
3. Try: One of the examples from [CODE_EXAMPLES.md](./CODE_EXAMPLES.md#basic-examples)

### For Experienced Users

- Implement pattern from [CODE_EXAMPLES.md](./CODE_EXAMPLES.md#common-patterns)
- Reference: [TROUBLESHOOTING.md - FAQ](./TROUBLESHOOTING.md#frequently-asked-questions)

### If Something Doesn't Work

- Check: [TROUBLESHOOTING.md - Common Issues](./TROUBLESHOOTING.md#common-issues)
- FAQ: [TROUBLESHOOTING.md - FAQ](./TROUBLESHOOTING.md#frequently-asked-questions)

---

## 📖 Documentation Structure

### README.md (This file)

**Quick reference and navigation**

- Quick start
- Common use cases
- Quick reference (commands and code)
- Architecture overview
- Best practices
- Common patterns

### HOW_TO_USE_LOGS.md

**Comprehensive guide (20 min read)**

1. Quick Start - Import and basic usage
2. Basic Usage - Simple logging methods
3. Log Levels - When to use each level
4. Structured Logging - Adding context with fields
5. Configuration - Setting up the logger
6. Common Patterns - Real usage patterns
7. Real-World Examples - Complete service examples
8. Best Practices - Do's and don'ts
9. Output Examples - See what logs look like
10. Troubleshooting - Common issues

**Topics Covered:**

- All logging methods (Debug, Info, Warn, Error, Fatal)
- Field types (String, Int, Int64, Float64, Bool, Time, Duration, Error, Any)
- Configuration via code and environment variables
- Dependency injection pattern
- Error handling best practices
- Request context propagation
- Performance considerations

### CODE_EXAMPLES.md

**Copy-paste ready examples (15 min read)**

1. Basic Examples - Minimal examples

   - Application entry point
   - Different log levels
   - Structured logging with fields

2. Service Examples - Real service implementations

   - User service with lifecycle
   - Order processing service

3. Handler Examples - HTTP handler patterns

   - Create user handler with request logging
   - HTTP middleware for logging
   - Error recovery middleware

4. Database Examples - Data layer patterns

   - Database connection and migration
   - Query logging with timing
   - Bulk operations

5. Error Handling Examples - Error patterns

   - Error categorization and logging
   - Error wrapping with context
   - Panic recovery

6. Testing Examples - Test patterns

   - Creating test logger
   - Testing services with logging
   - Mock implementations

7. Templates - Copy-paste ready code
   - Service template
   - Handler template
   - Repository template

### TROUBLESHOOTING.md

**Issues, questions, and optimization (25 min read)**

**8 Common Issues with Solutions:**

1. Logger not showing output
2. Logger is nil error
3. Logs not showing request details
4. Performance degradation
5. Duplicate log entries
6. Lost logs on panic
7. JSON logs not readable
8. Wrong caller information

**10 FAQ Questions:**

1. Can I use logger from packages?
2. Should I create logger in each file?
3. How to add request ID to all logs?
4. Can I use different levels per package?
5. How to handle sensitive data?
6. Should I catch and log errors?
7. How to test code with logging?
8. Can I send logs to external services?
9. Performance overhead of structured logging?
10. Why use go-infra logger vs Zap directly?

**Performance Tips:**

- Use debug level for development only
- Avoid reflection in hot paths
- Disable caller info if needed
- Batch write large amounts

**Integration Guides:**

- HTTP server integration
- Database query integration
- Background job integration

---

## 🔑 Key Concepts

### The Logger Interface

```go
type Logger interface {
    Debug(args ...interface{})
    Debugf(template string, args ...interface{})
    Debugw(msg string, fields Fields)

    Info(args ...interface{})
    Infof(template string, args ...interface{})
    Infow(msg string, fields Fields)

    Warn(args ...interface{})
    Warnf(template string, args ...interface{})
    WarnMsg(msg string, err error)

    Error(args ...interface{})
    Errorf(template string, args ...interface{})
    Errorw(msg string, fields Fields)
    Err(msg string, err error)

    Fatal(args ...interface{})
    Fatalf(template string, args ...interface{})

    // For request ID and context
    WithName(name string)

    // For gRPC
    GrpcMiddlewareAccessLogger(...)
    GrpcClientInterceptorLogger(...)
}
```

### Field Types

```go
logger.String(key, value)
logger.Int(key, value)
logger.Int64(key, value)
logger.Float64(key, value)
logger.Bool(key, value)
logger.Time(key, value)
logger.Duration(key, value)
logger.Err(error)
logger.Any(key, value)
```

### Log Levels (Severity)

1. **DEBUG** - Detailed diagnostic information
2. **INFO** - General informational messages
3. **WARN** - Warning messages for potential issues
4. **ERROR** - Error events that need investigation
5. **FATAL** - Severe errors that stop the application

---

## 🚀 5-Minute Quick Start

### Step 1: Create a Logger-Enabled Service

```go
package myservice

import "github.com/phatnt199/go-infra/pkg/logger"

type MyService struct {
    logger logger.Logger
}

func NewMyService(logger logger.Logger) *MyService {
    return &MyService{logger: logger}
}

func (s *MyService) DoSomething() error {
    s.logger.Info("Starting operation")

    if err := s.performWork(); err != nil {
        s.logger.Errorw("Operation failed",
            logger.String("operation", "doSomething"),
            logger.Err(err),
        )
        return err
    }

    s.logger.Info("Operation completed")
    return nil
}

func (s *MyService) performWork() error {
    // Your implementation
    return nil
}
```

### Step 2: Use Dependency Injection

```go
// In your fx.Options
fx.Provide(NewMyService)  // Receives logger automatically
```

### Step 3: Set Environment

```bash
export APP_ENV=development
export LogConfig_LogLevel=debug
```

### Step 4: Run and See Logs

```bash
go run main.go
# You'll see colorful console logs in development
# or structured JSON logs in production
```

---

## 💡 Pro Tips

### Tip 1: Add Request ID to Every Log

```go
requestID := r.Header.Get("X-Request-ID")
logger := appLogger.WithFields(
    logger.String("request_id", requestID),
)
```

### Tip 2: Use Different Levels Strategically

- **Debug**: Loop internals, parsed data
- **Info**: Important business events
- **Warn**: Recoverable issues
- **Error**: Actual errors with recovery
- **Fatal**: Application cannot continue

### Tip 3: Structure Your Data

```go
// Good - searchable and parseable
logger.Infow("Payment processed",
    logger.String("user_id", userID),
    logger.Float64("amount", amount),
    logger.Int64("duration_ms", duration),
)
```

### Tip 4: Log at Boundaries

- HTTP handlers
- Service methods
- Repository operations
- Event handlers

### Tip 5: Always Defer Sync()

```go
func main() {
    logger := initLogger()
    defer logger.Sync()  // NEVER SKIP THIS
    // ...
}
```

---

## 📊 Logger Comparison

| Feature     | Debug Level     | Info Level      | Production |
| ----------- | --------------- | --------------- | ---------- |
| Format      | Colored Console | Colored Console | JSON       |
| Caller Info | Full path       | Full path       | Short path |
| Debug logs  | ✅ Shown        | ❌ Hidden       | ❌ Hidden  |
| Info logs   | ✅ Shown        | ✅ Shown        | ✅ Shown   |
| Warn logs   | ✅ Shown        | ✅ Shown        | ✅ Shown   |
| Error logs  | ✅ Shown        | ✅ Shown        | ✅ Shown   |
| Performance | Good            | Better          | Best       |

---

## 🎓 Learning Path

### Level 1: Beginner (15 minutes)

- [ ] Read: [README.md](./README.md)
- [ ] Read: [HOW_TO_USE_LOGS.md - Quick Start](./HOW_TO_USE_LOGS.md#quick-start)
- [ ] Try: Run [CODE_EXAMPLES.md - Example 1](./CODE_EXAMPLES.md#example-1-simple-application-entry-point)

### Level 2: Intermediate (45 minutes)

- [ ] Read: [HOW_TO_USE_LOGS.md - Complete](./HOW_TO_USE_LOGS.md)
- [ ] Try: Implement one of [CODE_EXAMPLES.md - Service Examples](./CODE_EXAMPLES.md#service-examples)
- [ ] Reference: [CODE_EXAMPLES.md - Middleware Examples](./CODE_EXAMPLES.md#example-7-logging-middleware)

### Level 3: Advanced (90 minutes)

- [ ] Implement: Production service with all logging patterns
- [ ] Optimize: Performance using tips from [TROUBLESHOOTING.md](./TROUBLESHOOTING.md#performance-tips)
- [ ] Reference: [CODE_EXAMPLES.md - Database Examples](./CODE_EXAMPLES.md#database-examples)
- [ ] Reference: [CODE_EXAMPLES.md - Error Handling](./CODE_EXAMPLES.md#error-handling-examples)

---

## ❓ Quick FAQ

**Q: Where do I start?**
A: Read [README.md](./README.md), then the quick start section

**Q: I'm seeing nil logger errors**
A: Check [TROUBLESHOOTING.md - Issue 2](./TROUBLESHOOTING.md#issue-2-logger-is-nil-error)

**Q: How do I add user/request ID to logs?**
A: See [HOW_TO_USE_LOGS.md - Context-Based Logging](./HOW_TO_USE_LOGS.md#pattern-2-transaction-logging-real-example-from-project) and [CODE_EXAMPLES.md - Example 6](./CODE_EXAMPLES.md#example-6-http-handler-with-request-logging)

**Q: Performance is slow**
A: Check [TROUBLESHOOTING.md - Issue 4](./TROUBLESHOOTING.md#issue-4-performance-degradation-with-logging)

**Q: Should I log sensitive data?**
A: No! Check [TROUBLESHOOTING.md - Q5](./TROUBLESHOOTING.md#q5-how-do-i-handle-sensitive-data-in-logs)

---

## 📚 Files Guide

```
examples/logs/
├── README.md                    # This file
├── HOW_TO_USE_LOGS.md          # Complete guide (START HERE!)
├── CODE_EXAMPLES.md            # Code samples
└── TROUBLESHOOTING.md          # Issues & FAQ
```

---

## 🔗 External Resources

- [Logger Package](../../pkg/logger)
- [Zap Documentation](https://pkg.go.dev/go.uber.org/zap)
- [Go-Infra Project](https://github.com/phatnt199/go-infra)
- [Structured Logging Best Practices](https://www.kartar.net/2015/12/structured-logging/)

---

## 📝 Contributing

Found an issue or have a suggestion?

1. Check existing docs
2. Check [TROUBLESHOOTING.md](./TROUBLESHOOTING.md)
3. Raise an issue or create a PR

---

**Last Updated:** 2025-01-15
**Logger Package Version:** Latest
**Status:** Complete ✅

---

**Happy Logging! 🎉**
