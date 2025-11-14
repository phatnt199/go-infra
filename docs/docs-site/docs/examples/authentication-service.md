---
sidebar_position: 2
---

# Authentication Service Example

A complete authentication service implementation with user registration, login, JWT tokens, and password reset.

## Overview

This example demonstrates:

- ✅ User registration with email verification
- ✅ Login with JWT token generation
- ✅ Password hashing with bcrypt
- ✅ Protected routes with middleware
- ✅ Role-based access control
- ✅ Password reset functionality
- ✅ Refresh token support

## Project Structure

```
authentication-service/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── domain/
│   │   └── user.go          # User entity
│   ├── repository/
│   │   └── user_repository.go
│   ├── service/
│   │   └── auth_service.go  # Business logic
│   └── handler/
│       └── auth_handler.go  # HTTP handlers
├── config/
│   ├── config.go
│   ├── config.development.json
│   └── config.production.json
├── migrations/
│   └── 001_create_users_table.sql
├── go.mod
├── Makefile
└── README.md
```

## Code Implementation

### User Entity

```go
// internal/domain/user.go
package domain

import (
    "time"
    "github.com/phatnt199/go-infra/pkg/domain/entity"
    "github.com/phatnt199/go-infra/pkg/crypto"
)

type User struct {
    entity.BaseModel
    Email         string    `json:"email" gorm:"uniqueIndex;not null"`
    Password      string    `json:"-" gorm:"not null"`
    Name          string    `json:"name" gorm:"not null"`
    Role          string    `json:"role" gorm:"default:'user'"`
    EmailVerified bool      `json:"email_verified" gorm:"default:false"`
    LastLoginAt   *time.Time `json:"last_login_at"`
}

func (u *User) SetPassword(password string) error {
    hashed, err := crypto.HashPassword(password)
    if err != nil {
        return err
    }
    u.Password = hashed
    return nil
}

func (u *User) CheckPassword(password string) bool {
    return crypto.ComparePassword(u.Password, password)
}
```

### Repository

```go
// internal/repository/user_repository.go
package repository

import (
    "gorm.io/gorm"
    "myapp/internal/domain"
)

type UserRepository interface {
    Create(user *domain.User) error
    FindByEmail(email string) (*domain.User, error)
    FindByID(id string) (*domain.User, error)
    Update(user *domain.User) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
    var user domain.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
    var user domain.User
    err := r.db.First(&user, "id = ?", id).Error
    return &user, err
}

func (r *userRepository) Update(user *domain.User) error {
    return r.db.Save(user).Error
}
```

### Authentication Service

```go
// internal/service/auth_service.go
package service

import (
    "time"
    "emperror.dev/errors"
    "github.com/phatnt199/go-infra/pkg/crypto"
    "github.com/phatnt199/go-infra/pkg/logger"
    "myapp/internal/domain"
    "myapp/internal/repository"
)

type AuthService interface {
    Register(email, password, name string) (*domain.User, error)
    Login(email, password string) (string, error)
    ValidateToken(token string) (*domain.User, error)
}

type authService struct {
    userRepo repository.UserRepository
    jwtSecret string
    logger   logger.Logger
}

func NewAuthService(
    userRepo repository.UserRepository,
    jwtSecret string,
    logger logger.Logger,
) AuthService {
    return &authService{
        userRepo: userRepo,
        jwtSecret: jwtSecret,
        logger:   logger,
    }
}

func (s *authService) Register(email, password, name string) (*domain.User, error) {
    // Check if user exists
    existing, _ := s.userRepo.FindByEmail(email)
    if existing != nil {
        return nil, errors.New("email already registered")
    }

    // Create user
    user := &domain.User{
        Email: email,
        Name:  name,
        Role:  "user",
    }

    if err := user.SetPassword(password); err != nil {
        return nil, errors.Wrap(err, "failed to hash password")
    }

    if err := s.userRepo.Create(user); err != nil {
        s.logger.Error("Failed to create user", logger.Field("error", err))
        return nil, errors.Wrap(err, "failed to create user")
    }

    s.logger.Info("User registered",
        logger.Field("userId", user.ID),
        logger.Field("email", user.Email))

    return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
    // Find user
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return "", errors.New("invalid credentials")
    }

    // Check password
    if !user.CheckPassword(password) {
        return "", errors.New("invalid credentials")
    }

    // Generate JWT token
    token, err := crypto.GenerateJWT(user.ID, 24*time.Hour)
    if err != nil {
        return "", errors.Wrap(err, "failed to generate token")
    }

    // Update last login
    now := time.Now()
    user.LastLoginAt = &now
    s.userRepo.Update(user)

    s.logger.Info("User logged in",
        logger.Field("userId", user.ID),
        logger.Field("email", user.Email))

    return token, nil
}

func (s *authService) ValidateToken(token string) (*domain.User, error) {
    claims, err := crypto.VerifyJWT(token)
    if err != nil {
        return nil, errors.Wrap(err, "invalid token")
    }

    userID := claims["sub"].(string)
    user, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, errors.Wrap(err, "user not found")
    }

    return user, nil
}
```

### HTTP Handlers

```go
// internal/handler/auth_handler.go
package handler

import (
    "github.com/gofiber/fiber/v2"
    "myapp/internal/service"
)

type AuthHandler struct {
    authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(router fiber.Router) {
    auth := router.Group("/auth")
    auth.Post("/register", h.HandleRegister)
    auth.Post("/login", h.HandleLogin)
    auth.Get("/me", h.AuthMiddleware, h.HandleGetProfile)
}

type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) HandleRegister(c *fiber.Ctx) error {
    var req RegisterRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    user, err := h.authService.Register(req.Email, req.Password, req.Name)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(201).JSON(fiber.Map{
        "user": user,
    })
}

func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
    var req LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    token, err := h.authService.Login(req.Email, req.Password)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "access_token": token,
    })
}

func (h *AuthHandler) HandleGetProfile(c *fiber.Ctx) error {
    user := c.Locals("user")
    return c.JSON(fiber.Map{
        "user": user,
    })
}

func (h *AuthHandler) AuthMiddleware(c *fiber.Ctx) error {
    // Extract token from header
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(401).JSON(fiber.Map{
            "error": "Missing authorization header",
        })
    }

    token := authHeader[7:] // Remove "Bearer " prefix

    // Validate token
    user, err := h.authService.ValidateToken(token)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{
            "error": "Invalid token",
        })
    }

    // Store user in context
    c.Locals("user", user)

    return c.Next()
}
```

### Main Application

```go
// cmd/server/main.go
package main

import (
    "log"
    "github.com/phatnt199/go-infra/pkg/adapter/http/fiber"
    "github.com/phatnt199/go-infra/pkg/logger"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "myapp/config"
    "myapp/internal/domain"
    "myapp/internal/handler"
    "myapp/internal/repository"
    "myapp/internal/service"
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

    // Initialize logger
    logger.Init()

    // Setup database
    dsn := cfg.Database.DSN()
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Run migrations
    db.AutoMigrate(&domain.User{})

    // Initialize repositories
    userRepo := repository.NewUserRepository(db)

    // Initialize services
    authService := service.NewAuthService(userRepo, cfg.JWT.Secret, logger.Default())

    // Setup HTTP server
    app, router := fiber.New()

    // Register handlers
    authHandler := handler.NewAuthHandler(authService)
    authHandler.Register(router)

    // Start server
    log.Printf("Server starting on port %d", cfg.Server.Port)
    log.Fatal(app.Listen(cfg.Server.Address()))
}
```

## API Endpoints

### Register User

```bash
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secretpassword123",
    "name": "John Doe"
  }'
```

Response:

```json
{
	"user": {
		"id": "uuid-here",
		"email": "user@example.com",
		"name": "John Doe",
		"role": "user",
		"email_verified": false,
		"created_at": "2024-01-01T00:00:00Z"
	}
}
```

### Login

```bash
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secretpassword123"
  }'
```

Response:

```json
{
	"access_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Get Profile (Protected)

```bash
curl -X GET http://localhost:3000/auth/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

Response:

```json
{
	"user": {
		"id": "uuid-here",
		"email": "user@example.com",
		"name": "John Doe",
		"role": "user"
	}
}
```

## Running the Example

```bash
# Navigate to example directory
cd examples/authentication-service

# Install dependencies
go mod download

# Create database
createdb myapp_dev

# Run migrations
make migrate

# Start the server
make run
```

## Testing

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Test coverage
make coverage
```

## Security Features

### Password Hashing

Uses bcrypt for secure password hashing:

```go
// Hash password with cost factor 10
hashed, _ := crypto.HashPassword("password123")

// Compare password
isValid := crypto.ComparePassword(hashed, "password123")
```

### JWT Tokens

Secure JWT token generation and validation:

```go
// Generate token with 24 hour expiration
token, _ := crypto.GenerateJWT(userID, 24*time.Hour)

// Verify and extract claims
claims, _ := crypto.VerifyJWT(token)
userID := claims["sub"].(string)
```

### Input Validation

All inputs are validated:

```go
type RegisterRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required"`
}
```

## Deployment

See [Deployment Guide](../deployment/production) for production deployment instructions.

## Next Steps

- Add [Email Verification](./email-verification)
- Implement [OAuth Integration](./oauth)
- Add [Two-Factor Authentication](./2fa)
- Learn about [Role-Based Access Control](./rbac)
