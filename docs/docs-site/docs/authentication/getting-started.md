---
sidebar_position: 1
---

# Getting Started

Learn how to add authentication to your go-infra application.

## Overview

The authentication component provides a complete, production-ready authentication system with:

- ✅ **Multiple identifier schemes** - Username, email, phone
- ✅ **Multiple credential schemes** - Password, OAuth, API keys
- ✅ **JWT tokens** - Access and refresh tokens
- ✅ **Password hashing** - Bcrypt and Argon2
- ✅ **Middleware** - Protect routes easily
- ✅ **Role-based access** - Built-in authorization

## Quick Setup

### 1. Add Authentication Module

```go
import (
    authComponent "github.com/phatnt199/go-infra/pkg/component/authentication"
)

func main() {
    // Create authentication component
    authComp, err := authComponent.NewComponentWithDefaultImplementation(
        db,
        jwtConfig,
        authComponent.WithLogger(logger),
    )

    // Get services
    authService := authComp.GetAuthService()
    jwtMiddleware := authComp.GetJWTMiddleware()
}
```

### 2. Configuration

```bash title=".env"
# JWT Configuration
JWT_SECRET=your-super-secret-key-change-in-production
JWT_ISSUER=myapp
JWT_AUDIENCE=myapp-users
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d
```

### 3. Database Tables

The component requires these tables:

- `ms_users` - Core user entity
- `ms_user_identifiers` - Usernames, emails, phones
- `ms_user_credentials` - Passwords and credentials
- `ms_user_profiles` - User profile information

Run migrations:

```bash
go run cmd/migration/main.go
```

## Sign Up

```go
type SignUpRequest struct {
    Identifier contracts.Identifier `json:"identifier"`
    Credential contracts.Credential `json:"credential"`
    Firstname  string               `json:"firstname"`
    Lastname   string               `json:"lastname"`
    Email      string               `json:"email"`
}

func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
    var req SignUpRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    // Sign up user
    result, err := h.authService.SignUp(c.Context(), &contracts.SignUpModel{
        Identifier: &req.Identifier,
        Credential: &req.Credential,
        Profile: &contracts.ProfileSignUpModel{
            Firstname: req.Firstname,
            Lastname:  req.Lastname,
            Email:     req.Email,
        },
    })

    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(201).JSON(result)
}
```

### Request Example

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": {
      "scheme": "username",
      "value": "john_doe"
    },
    "credential": {
      "scheme": "basic",
      "value": "SecurePassword123!"
    },
    "firstname": "John",
    "lastname": "Doe",
    "email": "john@example.com"
  }'
```

## Sign In

```go
type SignInRequest struct {
    Identifier contracts.Identifier `json:"identifier"`
    Credential contracts.Credential `json:"credential"`
}

func (h *AuthHandler) SignIn(c *fiber.Ctx) error {
    var req SignInRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    // Authenticate user
    result, err := h.authService.SignIn(c.Context(), &contracts.SignInModel{
        Identifier: &req.Identifier,
        Credential: &req.Credential,
    })

    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
    }

    return c.JSON(result)
}
```

### Request Example

```bash
curl -X POST http://localhost:8080/api/v1/auth/signin \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": {
      "scheme": "username",
      "value": "john_doe"
    },
    "credential": {
      "scheme": "basic",
      "value": "SecurePassword123!"
    }
  }'
```

### Response

```json
{
	"userId": "123e4567-e89b-12d3-a456-426614174000",
	"username": "john_doe",
	"email": "john@example.com",
	"token": {
		"value": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		"scheme": "bearer",
		"expiresAt": "2024-01-01T01:00:00Z",
		"type": "access"
	}
}
```

## Protect Routes

### Using Middleware

```go
func setupRoutes(
    server contracts.HttpServer,
    authComp *authComponent.Component,
    handler *UserHandler,
) {
    server.RouteBuilder().RegisterHandler(func(router interface{}) {
        app := router.(*fiber.App)

        // Public routes
        app.Post("/auth/signin", authHandler.SignIn)
        app.Post("/auth/signup", authHandler.SignUp)

        // Protected routes
        api := app.Group("/api/v1")
        api.Use(authComp.GetJWTMiddleware().Handle())

        api.Get("/profile", handler.GetProfile)
        api.Put("/profile", handler.UpdateProfile)
        api.Get("/users", handler.GetUsers)
    })
}
```

### Access User from Token

```go
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
    // Get user ID from JWT claims
    userID := c.Locals("user_id").(string)

    // Fetch user
    user, err := h.service.GetUserByID(c.Context(), userID)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "User not found"})
    }

    return c.JSON(user)
}
```

## Identifier Schemes

### Username

```json
{
	"identifier": {
		"scheme": "username",
		"value": "john_doe"
	}
}
```

### Email

```json
{
	"identifier": {
		"scheme": "email",
		"value": "john@example.com"
	}
}
```

### Phone

```json
{
	"identifier": {
		"scheme": "phone",
		"value": "+1234567890"
	}
}
```

## Credential Schemes

### Basic (Password)

```json
{
	"credential": {
		"scheme": "basic",
		"value": "MySecurePassword123!"
	}
}
```

### OAuth

```json
{
	"credential": {
		"scheme": "oauth",
		"value": "google_token_here"
	}
}
```

### API Key

```json
{
	"credential": {
		"scheme": "apikey",
		"value": "api_key_here"
	}
}
```

## Change Password

```go
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(string)

    type ChangePasswordRequest struct {
        OldCredential contracts.Credential `json:"oldCredential"`
        NewCredential contracts.Credential `json:"newCredential"`
    }

    var req ChangePasswordRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    err := h.authService.ChangePassword(c.Context(), &contracts.ChangePasswordModel{
        UserID:        userID,
        OldCredential: &req.OldCredential,
        NewCredential: &req.NewCredential,
    })

    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Password changed successfully"})
}
```

## Token Refresh

```go
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
    type RefreshRequest struct {
        RefreshToken string `json:"refreshToken"`
    }

    var req RefreshRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    // Refresh token
    newToken, err := h.jwtManager.RefreshToken(req.RefreshToken)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid refresh token"})
    }

    return c.JSON(fiber.Map{
        "accessToken": newToken,
    })
}
```

## Role-Based Access

```go
// Check for specific role
func AdminOnly() fiber.Handler {
    return func(c *fiber.Ctx) error {
        roles := c.Locals("roles").([]string)

        hasAdmin := false
        for _, role := range roles {
            if role == "admin" {
                hasAdmin = true
                break
            }
        }

        if !hasAdmin {
            return c.Status(403).JSON(fiber.Map{
                "error": "Forbidden: Admin access required",
            })
        }

        return c.Next()
    }
}

// Use in routes
api.Get("/admin/dashboard", AdminOnly(), handler.GetDashboard)
```

## Custom User Provider

If you have your own user model:

```go
type MyUserProvider struct {
    db *gorm.DB
}

func (p *MyUserProvider) GetUserByIdentifier(
    ctx context.Context,
    scheme, value string,
) (*contracts.UserAuthInfo, error) {
    var user MyUser
    err := p.db.Where("email = ?", value).First(&user).Error
    if err != nil {
        return nil, err
    }

    return &contracts.UserAuthInfo{
        UserID:   user.ID,
        Username: user.Username,
        Email:    user.Email,
        Roles:    user.Roles,
    }, nil
}

// Create component with custom provider
authComp := authComponent.NewComponent(
    &MyUserProvider{db: db},
    authComponent.WithLogger(logger),
)
```

## Complete Example

See the [Authentication Service Example](../examples/authentication-service) for a complete working implementation with:

- Sign up and sign in
- Profile management
- Password change
- Protected routes
- Admin routes
- Swagger documentation

## Security Best Practices

### 1. Strong Secrets

```bash
# Generate secure secret
openssl rand -base64 32
```

### 2. Short Token Expiry

```bash
JWT_ACCESS_EXPIRY=15m   # 15 minutes
JWT_REFRESH_EXPIRY=7d   # 7 days
```

### 3. HTTPS Only

Always use HTTPS in production.

### 4. Rate Limiting

Implement rate limiting on auth endpoints:

```go
import "github.com/gofiber/fiber/v2/middleware/limiter"

auth.Use(limiter.New(limiter.Config{
    Max:        5,
    Expiration: 1 * time.Minute,
}))
```

### 5. Password Requirements

```go
func ValidatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }

    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

    if !hasUpper || !hasLower || !hasNumber {
        return errors.New("password must contain uppercase, lowercase, and number")
    }

    return nil
}
```

## Next Steps

- **[JWT Tokens](./jwt)** - Deep dive into JWT
- **[Password Hashing](./password-hashing)** - Secure password storage
- **[Authorization](./authorization)** - Role-based access control
- **[Complete Example](../examples/authentication-service)** - Full implementation
