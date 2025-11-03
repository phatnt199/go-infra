# Logger - Code Examples

This document provides practical code examples for using the logger in the go-infra project.

## Table of Contents

- [Basic Examples](#basic-examples)
- [Service Examples](#service-examples)
- [Handler Examples](#handler-examples)
- [Database Examples](#database-examples)
- [Error Handling Examples](#error-handling-examples)
- [Middleware Examples](#middleware-examples)
- [Testing Examples](#testing-examples)

---

## Basic Examples

### Example 1: Simple Application Entry Point

```go
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/phatnt199/go-infra/pkg/logger"
	"github.com/phatnt199/go-infra/pkg/logger/config"
	"github.com/phatnt199/go-infra/pkg/logger/models"
	"github.com/phatnt199/go-infra/pkg/logger/zap"
	"github.com/phatnt199/go-infra/pkg/application/environment"
)

func main() {
	// Create logger configuration
	logOptions := &config.LogOptions{
		LogLevel:      "info",
		LogType:       models.Zap,
		CallerEnabled: true,
		EnableTracing: false,
	}

	// Get environment
	env := environment.NewEnvironment("development")

	// Create logger
	appLogger := zap.NewZapLogger(logOptions, env)
	defer appLogger.Sync()

	// Use logger
	appLogger.Info("Application started")
	appLogger.Infof("Running on environment: %s", env.GetEnvironment())

	// Simulate work
	time.Sleep(1 * time.Second)

	appLogger.Info("Application shutdown")
}
```

### Example 2: Different Log Levels

```go
package examples

import "github.com/phatnt199/go-infra/pkg/logger"

func LogLevelExamples(logger logger.Logger) {
	// Debug - Most detailed, for development
	logger.Debug("Starting database query execution")
	logger.Debugf("Query parameters: user_id=%d, limit=%d", 123, 10)

	// Info - General informational messages
	logger.Info("User authenticated successfully")
	logger.Infof("Processed %d records in batch", 1000)

	// Warn - Warning for potential issues
	logger.Warn("Database connection pool is running low")
	logger.Warnf("Response time for API call: %dms (threshold: 100ms)", 250)

	// Error - Error events that need investigation
	logger.Error("Failed to save user profile")
	logger.Errorf("Connection to cache server failed: %s", "connection timeout")

	// Fatal - Severe errors that stop the application
	// Be careful with this - it will exit the program
	// logger.Fatal("Critical: Cannot connect to database")
}
```

### Example 3: Structured Logging with Fields

```go
package examples

import (
	"time"
	"github.com/phatnt199/go-infra/pkg/logger"
)

func StructuredLoggingExample(logger logger.Logger) {
	// Log with string fields
	logger.Infow("User action performed",
		logger.String("user_id", "usr-123"),
		logger.String("action", "login"),
		logger.String("ip_address", "192.168.1.100"),
	)

	// Log with numeric fields
	logger.Infow("API request processed",
		logger.String("endpoint", "/api/users"),
		logger.Int("status_code", 200),
		logger.Int64("response_time_ms", 145),
		logger.Float64("cpu_usage_percent", 45.3),
	)

	// Log with boolean fields
	logger.Infow("Feature flag evaluation",
		logger.String("feature", "new_checkout"),
		logger.Bool("enabled", true),
		logger.Bool("is_beta_user", false),
	)

	// Log with time fields
	logger.Infow("Cache entry created",
		logger.String("key", "user-profile-123"),
		logger.Time("created_at", time.Now()),
		logger.Duration("ttl", 24*time.Hour),
	)
}
```

---

## Service Examples

### Example 4: User Service with Full Lifecycle

```go
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/phatnt199/go-infra/pkg/logger"
)

type User struct {
	ID    string
	Email string
	Name  string
	Active bool
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

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

// CreateUser demonstrates structured logging with validation and error handling
func (s *UserService) CreateUser(ctx context.Context, email, name string) (*User, error) {
	start := time.Now()

	// Log operation start
	s.logger.Infow("Creating new user",
		logger.String("email", email),
		logger.String("name", name),
	)

	// Validate input
	if email == "" {
		s.logger.Errorw("Invalid user creation request - empty email",
			logger.String("name", name),
		)
		return nil, fmt.Errorf("email cannot be empty")
	}

	// Create user object
	user := &User{
		ID:     fmt.Sprintf("usr-%d", time.Now().UnixNano()),
		Email:  email,
		Name:   name,
		Active: true,
	}

	// Save to repository
	if err := s.repo.Create(ctx, user); err != nil {
		s.logger.Errorw("Failed to create user in repository",
			logger.String("email", email),
			logger.String("user_id", user.ID),
			logger.Err(err),
		)
		return nil, err
	}

	// Log success with execution time
	s.logger.Infow("User created successfully",
		logger.String("user_id", user.ID),
		logger.String("email", user.Email),
		logger.Duration("execution_time", time.Since(start)),
	)

	return user, nil
}

// GetUser demonstrates debug logging and not found errors
func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
	s.logger.Debugf("Fetching user with ID: %s", userID)

	user, err := s.repo.FindByID(ctx, userID)

	if err != nil {
		s.logger.Warnw("User not found",
			logger.String("user_id", userID),
			logger.Err(err),
		)
		return nil, fmt.Errorf("user not found: %w", err)
	}

	s.logger.Debugf("User found: %s (%s)", userID, user.Email)
	return user, nil
}

// UpdateUser demonstrates transaction-like logging
func (s *UserService) UpdateUser(ctx context.Context, user *User) error {
	s.logger.Infow("Updating user",
		logger.String("user_id", user.ID),
		logger.String("email", user.Email),
	)

	if err := s.repo.Update(ctx, user); err != nil {
		s.logger.Errorw("Failed to update user",
			logger.String("user_id", user.ID),
			logger.String("database", "postgres"),
			logger.Err(err),
		)
		return err
	}

	s.logger.Infof("User %s updated successfully", user.ID)
	return nil
}

// DeleteUser demonstrates cascade operation logging
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {
	s.logger.Infow("Deleting user",
		logger.String("user_id", userID),
	)

	// Validate user exists
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		s.logger.Errorw("Cannot delete - user not found",
			logger.String("user_id", userID),
			logger.Err(err),
		)
		return err
	}

	// Perform deletion
	if err := s.repo.Delete(ctx, userID); err != nil {
		s.logger.Errorw("Failed to delete user",
			logger.String("user_id", userID),
			logger.String("email", user.Email),
			logger.Err(err),
		)
		return err
	}

	s.logger.Infow("User deleted successfully",
		logger.String("user_id", userID),
		logger.String("email", user.Email),
	)

	return nil
}
```

### Example 5: Order Processing Service with Business Logic

```go
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/phatnt199/go-infra/pkg/logger"
)

type Order struct {
	ID       string
	UserID   string
	Total    float64
	Status   string
	Items    []OrderItem
	CreatedAt time.Time
}

type OrderItem struct {
	ProductID string
	Quantity  int
	Price     float64
}

type OrderService struct {
	logger logger.Logger
}

func NewOrderService(logger logger.Logger) *OrderService {
	return &OrderService{logger: logger}
}

func (s *OrderService) ProcessOrder(ctx context.Context, order *Order) error {
	start := time.Now()

	// Log order processing start
	s.logger.Infow("Processing order",
		logger.String("order_id", order.ID),
		logger.String("user_id", order.UserID),
		logger.Float64("total", order.Total),
		logger.Int("items_count", len(order.Items)),
	)

	// Validate order
	if len(order.Items) == 0 {
		s.logger.Errorw("Order validation failed - no items",
			logger.String("order_id", order.ID),
		)
		return fmt.Errorf("order must contain at least one item")
	}

	// Check inventory
	for i, item := range order.Items {
		s.logger.Debugf("Checking inventory for item %d: product=%s, qty=%d",
			i, item.ProductID, item.Quantity)

		// Simulate inventory check
		if item.Quantity <= 0 {
			s.logger.Errorw("Invalid quantity for order item",
				logger.String("order_id", order.ID),
				logger.String("product_id", item.ProductID),
				logger.Int("quantity", item.Quantity),
			)
			return fmt.Errorf("invalid quantity for product %s", item.ProductID)
		}
	}

	s.logger.Debugf("Inventory check passed for order %s", order.ID)

	// Process payment
	s.logger.Infow("Processing payment",
		logger.String("order_id", order.ID),
		logger.Float64("amount", order.Total),
	)

	// Simulate payment processing
	if order.Total > 1000000 {
		s.logger.Warnw("High-value order",
			logger.String("order_id", order.ID),
			logger.Float64("total", order.Total),
		)
	}

	// Update order status
	order.Status = "completed"

	// Log completion
	s.logger.Infow("Order processed successfully",
		logger.String("order_id", order.ID),
		logger.String("status", order.Status),
		logger.Duration("processing_time", time.Since(start)),
	)

	return nil
}
```

---

## Handler Examples

### Example 6: HTTP Handler with Request Logging

```go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/phatnt199/go-infra/pkg/logger"
)

type CreateUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserHandler struct {
	logger      logger.Logger
	userService UserService
}

func NewUserHandler(logger logger.Logger, userService UserService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: userService,
	}
}

// CreateUser demonstrates HTTP request/response logging
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Extract request ID from header or generate one
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = fmt.Sprintf("req-%d", time.Now().UnixNano())
	}

	// Create scoped logger with request context
	reqLogger := h.logger
	_ = reqLogger // In real scenario, would use WithFields for request ID

	// Log request start
	h.logger.Infow("Received create user request",
		logger.String("request_id", requestID),
		logger.String("method", r.Method),
		logger.String("path", r.URL.Path),
		logger.String("remote_addr", r.RemoteAddr),
	)

	// Parse request body
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("Failed to parse request body",
			logger.String("request_id", requestID),
			logger.Err(err),
		)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Log parsed request
	h.logger.Debugw("Parsed request data",
		logger.String("request_id", requestID),
		logger.String("email", req.Email),
	)

	// Create user
	user, err := h.userService.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		h.logger.Errorw("Failed to create user",
			logger.String("request_id", requestID),
			logger.String("email", req.Email),
			logger.Err(err),
		)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Log response
	h.logger.Infow("User created successfully",
		logger.String("request_id", requestID),
		logger.String("user_id", user.ID),
		logger.String("email", user.Email),
		logger.Duration("response_time", time.Since(start)),
	)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser demonstrates GET handler logging
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	h.logger.Debugf("Get user request for ID: %s", userID)

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		h.logger.Warnw("User not found",
			logger.String("user_id", userID),
			logger.Err(err),
		)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	h.logger.Infof("Returning user %s", userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
```

### Example 7: Logging Middleware

```go
package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/phatnt199/go-infra/pkg/logger"
)

// LoggingMiddleware logs HTTP requests and responses
func LoggingMiddleware(logger logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Generate or extract request ID
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = fmt.Sprintf("req-%d", start.UnixNano())
			}

			// Log request
			logger.Infow("HTTP request received",
				logger.String("request_id", requestID),
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
				logger.String("query", r.URL.RawQuery),
				logger.String("remote_addr", r.RemoteAddr),
				logger.String("user_agent", r.UserAgent()),
			)

			// Call next handler
			next.ServeHTTP(w, r)

			// Log response
			duration := time.Since(start)
			logger.Infow("HTTP request completed",
				logger.String("request_id", requestID),
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
				logger.Duration("duration", duration),
			)
		})
	}
}

// ErrorRecoveryMiddleware catches and logs panics
func ErrorRecoveryMiddleware(logger logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Errorw("Request handler panicked",
						logger.String("method", r.Method),
						logger.String("path", r.URL.Path),
						logger.Any("panic", err),
					)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
```

---

## Database Examples

### Example 8: Database Connection and Migration

```go
package database

import (
	"fmt"
	"time"

	"github.com/phatnt199/go-infra/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

type DatabaseManager struct {
	logger logger.Logger
	db     *gorm.DB
}

func NewDatabaseManager(logger logger.Logger) *DatabaseManager {
	return &DatabaseManager{logger: logger}
}

// Connect demonstrates connection logging
func (m *DatabaseManager) Connect(connectionString string) error {
	m.logger.Infow("Attempting to connect to database",
		logger.String("connection_type", "postgres"),
	)

	start := time.Now()

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		m.logger.Errorw("Failed to connect to database",
			logger.String("error_type", fmt.Sprintf("%T", err)),
			logger.Duration("attempt_duration", time.Since(start)),
			logger.Err(err),
		)
		return err
	}

	m.logger.Infow("Database connection established",
		logger.Duration("connection_time", time.Since(start)),
	)

	m.db = db
	return nil
}

// RunMigrations demonstrates migration logging
func (m *DatabaseManager) RunMigrations(models ...interface{}) error {
	m.logger.Info("Starting database migrations")

	for _, model := range models {
		modelName := fmt.Sprintf("%T", model)
		m.logger.Debugf("Migrating model: %s", modelName)

		if err := m.db.AutoMigrate(model); err != nil {
			m.logger.Errorw("Failed to migrate model",
				logger.String("model", modelName),
				logger.Err(err),
			)
			return err
		}

		m.logger.Infof("Successfully migrated: %s", modelName)
	}

	m.logger.Info("All database migrations completed")
	return nil
}

// Close demonstrates graceful shutdown logging
func (m *DatabaseManager) Close() error {
	m.logger.Info("Closing database connection")

	sqlDB, err := m.db.DB()
	if err != nil {
		m.logger.Errorw("Failed to get database connection",
			logger.Err(err),
		)
		return err
	}

	if err := sqlDB.Close(); err != nil {
		m.logger.Errorw("Failed to close database connection",
			logger.Err(err),
		)
		return err
	}

	m.logger.Info("Database connection closed gracefully")
	return nil
}
```

### Example 9: Query Logging

```go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/phatnt199/go-infra/pkg/logger"
	"gorm.io/gorm"
)

type QueryLogger struct {
	logger logger.Logger
	db     *gorm.DB
}

func NewQueryLogger(logger logger.Logger, db *gorm.DB) *QueryLogger {
	return &QueryLogger{logger: logger, db: db}
}

// FindByID demonstrates query logging
func (q *QueryLogger) FindByID(ctx context.Context, table string, id interface{}) *gorm.DB {
	start := time.Now()

	q.logger.Debugf("Executing query: SELECT * FROM %s WHERE id = ?", table)

	result := q.db.WithContext(ctx).Table(table).Where("id = ?", id)

	duration := time.Since(start)

	if result.Error != nil {
		q.logger.Errorw("Query execution failed",
			logger.String("table", table),
			logger.Any("id", id),
			logger.Duration("query_time", duration),
			logger.Err(result.Error),
		)
	} else {
		q.logger.Debugw("Query executed successfully",
			logger.String("table", table),
			logger.Duration("query_time", duration),
		)
	}

	return result
}

// BulkInsert demonstrates batch operation logging
func (q *QueryLogger) BulkInsert(ctx context.Context, table string, records interface{}) error {
	start := time.Now()

	q.logger.Infof("Starting bulk insert into %s", table)

	result := q.db.WithContext(ctx).Table(table).CreateInBatches(records, 100)

	duration := time.Since(start)

	if result.Error != nil {
		q.logger.Errorw("Bulk insert failed",
			logger.String("table", table),
			logger.Int64("rows_affected", result.RowsAffected),
			logger.Duration("operation_time", duration),
			logger.Err(result.Error),
		)
		return result.Error
	}

	q.logger.Infow("Bulk insert completed",
		logger.String("table", table),
		logger.Int64("rows_inserted", result.RowsAffected),
		logger.Duration("operation_time", duration),
	)

	return nil
}
```

---

## Error Handling Examples

### Example 10: Comprehensive Error Handling

```go
package errorhandling

import (
	"fmt"
	"github.com/phatnt199/go-infra/pkg/logger"
)

type ValidationError struct {
	Field  string
	Reason string
}

type PaymentError struct {
	Code    string
	Message string
}

// HandleErrors demonstrates error categorization and logging
func HandleErrors(logger logger.Logger, err error) {
	if err == nil {
		return
	}

	// Log different error types differently
	switch v := err.(type) {
	case *ValidationError:
		logger.Warnw("Validation error",
			logger.String("field", v.Field),
			logger.String("reason", v.Reason),
		)

	case *PaymentError:
		logger.Errorw("Payment processing failed",
			logger.String("error_code", v.Code),
			logger.String("error_message", v.Message),
		)

	default:
		logger.Errorw("Unknown error occurred",
			logger.String("error_type", fmt.Sprintf("%T", err)),
			logger.Err(err),
		)
	}
}

// WrapError demonstrates error wrapping with logging context
func WrapError(logger logger.Logger, err error, context map[string]interface{}) error {
	if err == nil {
		return nil
	}

	// Log with context
	fields := make([]interface{}, 0)
	for key, value := range context {
		// Convert to logger fields (simplified)
		fields = append(fields, key, value)
	}

	logger.Errorw("Error with context",
		logger.Err(err),
	)

	// Wrap error for propagation
	return fmt.Errorf("operation failed: %w", err)
}

// RecoveryHandler demonstrates panic recovery logging
func RecoveryHandler(logger logger.Logger) {
	if r := recover(); r != nil {
		logger.Errorw("Panic recovered",
			logger.Any("panic_value", r),
		)
		panic(r) // Re-panic if needed
	}
}
```

---

## Testing Examples

### Example 11: Testing with Logger

```go
package tests

import (
	"context"
	"testing"

	"github.com/phatnt199/go-infra/pkg/logger"
	"github.com/phatnt199/go-infra/pkg/logger/config"
	"github.com/phatnt199/go-infra/pkg/logger/models"
	"github.com/phatnt199/go-infra/pkg/logger/zap"
	"github.com/phatnt199/go-infra/pkg/application/environment"
)

// CreateTestLogger creates a logger for testing
func CreateTestLogger(t *testing.T) logger.Logger {
	logOptions := &config.LogOptions{
		LogLevel:      "debug",
		LogType:       models.Zap,
		CallerEnabled: false,
		EnableTracing: false,
	}

	env := environment.NewEnvironment("testing")
	return zap.NewZapLogger(logOptions, env)
}

// TestUserService demonstrates service testing with logging
func TestUserService(t *testing.T) {
	logger := CreateTestLogger(t)
	defer logger.(*zap.zapLogger).Sync()

	// Create mock repository
	mockRepo := &MockUserRepository{}

	// Create service with logger
	service := NewUserService(logger, mockRepo)

	// Test creating user
	ctx := context.Background()
	user, err := service.CreateUser(ctx, "test@example.com", "Test User")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Email != "test@example.com" {
		t.Fatalf("Expected email 'test@example.com', got '%s'", user.Email)
	}

	// Logger will have printed debug information
	logger.Infof("Test passed: user_id=%s", user.ID)
}

// MockUserRepository is a mock implementation for testing
type MockUserRepository struct{}

func (m *MockUserRepository) Create(ctx context.Context, user *User) error {
	user.ID = "mock-id-123"
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*User, error) {
	return &User{
		ID:    id,
		Email: "mock@example.com",
		Name:  "Mock User",
	}, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *User) error {
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	return nil
}
```

---

## Quick Copy-Paste Templates

### Service Template

```go
package services

import (
	"context"
	"github.com/phatnt199/go-infra/pkg/logger"
)

type MyService struct {
	logger logger.Logger
}

func NewMyService(logger logger.Logger) *MyService {
	return &MyService{logger: logger}
}

func (s *MyService) DoSomething(ctx context.Context) error {
	s.logger.Infow("Starting operation",
		logger.String("operation", "do_something"),
	)

	if err := s.performAction(); err != nil {
		s.logger.Errorw("Operation failed",
			logger.String("operation", "do_something"),
			logger.Err(err),
		)
		return err
	}

	s.logger.Infof("Operation completed successfully")
	return nil
}

func (s *MyService) performAction() error {
	// Implementation
	return nil
}
```

### Handler Template

```go
package handlers

import (
	"net/http"
	"github.com/phatnt199/go-infra/pkg/logger"
)

type MyHandler struct {
	logger logger.Logger
}

func NewMyHandler(logger logger.Logger) *MyHandler {
	return &MyHandler{logger: logger}
}

func (h *MyHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	h.logger.Infow("Handling request",
		logger.String("method", r.Method),
		logger.String("path", r.URL.Path),
	)

	// Process request
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
```

### Repository Template

```go
package repositories

import (
	"context"
	"github.com/phatnt199/go-infra/pkg/logger"
	"gorm.io/gorm"
)

type MyRepository struct {
	logger logger.Logger
	db     *gorm.DB
}

func NewMyRepository(logger logger.Logger, db *gorm.DB) *MyRepository {
	return &MyRepository{logger: logger, db: db}
}

func (r *MyRepository) Create(ctx context.Context, entity interface{}) error {
	r.logger.Debugf("Creating entity")

	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		r.logger.Errorw("Failed to create entity",
			logger.String("database", "postgres"),
			logger.Err(err),
		)
		return err
	}

	r.logger.Info("Entity created successfully")
	return nil
}
```

---

**Happy Coding! 🚀**
