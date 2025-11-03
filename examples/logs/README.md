# Logs Documentation

Complete guide to using the logging system in the go-infra project.

## 📚 Documentation Files

### 1. **[HOW_TO_USE_LOGS.md](./HOW_TO_USE_LOGS.md)** - Main Guide

Start here! Comprehensive guide covering:

- Quick start
- Basic usage
- Log levels
- Structured logging
- Configuration
- Common patterns
- Real-world examples
- Best practices

### 2. **[CODE_EXAMPLES.md](./CODE_EXAMPLES.md)** - Code Samples

Practical, copy-paste ready examples:

- Basic examples
- Service implementations
- HTTP handlers
- Database operations
- Error handling
- Middleware
- Testing with logger
- Quick templates

### 3. **[TROUBLESHOOTING.md](./TROUBLESHOOTING.md)** - FAQ & Issues

Solutions for common problems:

- 8 common issues with solutions
- 10 frequently asked questions
- Performance optimization tips
- Integration guides

---

## 🚀 Quick Start

### 1. Import the Logger

```go
import "github.com/phatnt199/go-infra/pkg/logger"
```

### 2. Use in Your Service

```go
type MyService struct {
    logger logger.Logger
}

func NewMyService(logger logger.Logger) *MyService {
    return &MyService{logger: logger}
}

func (s *MyService) DoWork() {
    s.logger.Info("Starting work")
    // ... your code ...
    s.logger.Info("Work completed")
}
```

### 3. Set Environment

```bash
export APP_ENV=development
export LogConfig_LogLevel=debug
```

---

## 📖 Common Use Cases

| Use Case             | Documentation                                                                                     |
| -------------------- | ------------------------------------------------------------------------------------------------- |
| Get started quickly  | [HOW_TO_USE_LOGS.md - Quick Start](./HOW_TO_USE_LOGS.md#quick-start)                              |
| Log levels explained | [HOW_TO_USE_LOGS.md - Log Levels](./HOW_TO_USE_LOGS.md#log-levels)                                |
| Structured logging   | [HOW_TO_USE_LOGS.md - Structured Logging](./HOW_TO_USE_LOGS.md#structured-logging)                |
| Service example      | [CODE_EXAMPLES.md - Service Examples](./CODE_EXAMPLES.md#service-examples)                        |
| Handler example      | [CODE_EXAMPLES.md - Handler Examples](./CODE_EXAMPLES.md#handler-examples)                        |
| Database logging     | [CODE_EXAMPLES.md - Database Examples](./CODE_EXAMPLES.md#database-examples)                      |
| Logger not working   | [TROUBLESHOOTING.md - Issue 1](./TROUBLESHOOTING.md#issue-1-logger-not-showing-any-output)        |
| Performance issues   | [TROUBLESHOOTING.md - Issue 4](./TROUBLESHOOTING.md#issue-4-performance-degradation-with-logging) |
| Question about usage | [TROUBLESHOOTING.md - FAQ](./TROUBLESHOOTING.md#frequently-asked-questions)                       |

---

## 🎯 Quick Reference

### Log Methods

```go
// Simple logging
logger.Debug(msg)
logger.Info(msg)
logger.Warn(msg)
logger.Error(msg)
logger.Fatal(msg)

// Formatted logging
logger.Debugf("format: %v", value)
logger.Infof("format: %v", value)
logger.Warnf("format: %v", value)
logger.Errorf("format: %v", value)
logger.Fatalf("format: %v", value)

// Structured logging with fields
logger.Debugw(msg, field1, field2, ...)
logger.Infow(msg, field1, field2, ...)
logger.Warnw(msg, field1, field2, ...)
logger.Errorw(msg, field1, field2, ...)
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

### Configuration

```bash
# Log level
export LogConfig_LogLevel=debug     # debug, info, warn, error, fatal

# Application environment
export APP_ENV=development          # development, production, staging

# Logger type
export LogConfig_LogType=0           # 0=Zap, 1=Logrus

# Enable caller information
export LogConfig_CallerEnabled=true  # true, false
```

---

## 🏗️ Architecture

The logging system is built on:

- **Zap**: High-performance structured logging library by Uber
- **Go.uber.org/fx**: Dependency injection framework
- **Hidden Implementation**: Users only see the clean logger interface

```
Your Application
    ↓
pkg/logger (clean interface)
    ↓
pkg/logger/zap (Zap implementation)
    ↓
go.uber.org/zap (underlying library)
```

---

## 📊 Output Examples

### Development Console Output

```
[INFO]  | 2025-01-15T10:30:45.123+0700 | main.go:42 | Application started
[DEBUG] | 2025-01-15T10:30:46.456+0700 | handler.go:15 | Processing request
[INFO]  | 2025-01-15T10:30:47.789+0700 | service.go:28 | User created | user_id=usr-456
[ERROR] | 2025-01-15T10:30:49.234+0700 | db.go:72 | Database error | error=connection timeout
```

### Production JSON Output

```json
{"level":"info","timestamp":"2025-01-15T10:30:47.789+0700","caller":"service.go:28","message":"User created","user_id":"usr-456"}
{"level":"error","timestamp":"2025-01-15T10:30:49.234+0700","caller":"db.go:72","message":"Database error","error":"connection timeout"}
```

---

## ⚡ Performance

The go-infra logger is optimized for performance:

- **Zero-allocation logging**: Uses Zap's allocation-free design
- **Minimal overhead**: ~1-2% compared to printf
- **High throughput**: 4+ million logs/second in benchmarks
- **Structured logging**: Efficient JSON encoding for production

---

## 🔗 Related Resources

- [Logger Package Documentation](../../pkg/logger/README.md)
- [Zap Logger](https://pkg.go.dev/go.uber.org/zap)
- [Go-Infra Repository](https://github.com/phatnt199/go-infra)
- [Structured Logging Best Practices](https://www.kartar.net/2015/12/structured-logging/)

---

## 📝 Common Patterns

### Pattern 1: Service with Dependency Injection

```go
type UserService struct {
    logger logger.Logger
    repo   UserRepository
}

func NewUserService(logger logger.Logger, repo UserRepository) *UserService {
    return &UserService{logger: logger, repo: repo}
}

func (s *UserService) CreateUser(email string) error {
    s.logger.Infof("Creating user: %s", email)
    // ...
}
```

### Pattern 2: HTTP Handler with Request Logging

```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    requestID := r.Header.Get("X-Request-ID")

    h.logger.Infow("Processing create user request",
        logger.String("request_id", requestID),
        logger.String("path", r.URL.Path),
    )

    // Process request...
}
```

### Pattern 3: Database Operations

```go
func (r *UserRepository) Create(user *User) error {
    r.logger.Debugf("Creating user: %s", user.Email)

    if err := r.db.Create(user).Error; err != nil {
        r.logger.Errorw("Create failed",
            logger.String("email", user.Email),
            logger.Err(err),
        )
        return err
    }

    r.logger.Infof("User created: %s", user.ID)
    return nil
}
```

### Pattern 4: Error Handling

```go
func (s *Service) Operation() error {
    if err := doSomething(); err != nil {
        s.logger.Errorw("Operation failed",
            logger.String("operation", "doSomething"),
            logger.Err(err),
        )
        return err
    }
    s.logger.Info("Operation completed successfully")
    return nil
}
```

---

## ✅ Best Practices

1. **Always defer Sync()** - Flush logs before exit
2. **Use structured logging** - Add context via fields, not string formatting
3. **Log at right level** - Debug for details, Info for events, Error for issues
4. **Include IDs and context** - User ID, request ID, operation name
5. **Don't log sensitive data** - Skip passwords, tokens, API keys
6. **Log once per error** - At the handler/middleware level
7. **Inject via dependency injection** - Don't create global loggers
8. **Use debug level for development** - Reduces noise in production

---

## ❓ Need Help?

- **How do I...?** → Check [CODE_EXAMPLES.md](./CODE_EXAMPLES.md)
- **Something's not working** → Check [TROUBLESHOOTING.md](./TROUBLESHOOTING.md)
- **How should I use it?** → Check [HOW_TO_USE_LOGS.md](./HOW_TO_USE_LOGS.md)

---

**Happy Logging! 📝✨**
