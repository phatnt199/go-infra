package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/models"
	"github.com/phatnt199/go-infra/pkg/component/authentication/strategies"
)

// JWTMiddleware handles JWT authentication using strategy pattern
type JWTMiddleware struct {
	strategy *strategies.JWTStrategy
}

// NewJWTMiddleware creates a new JWT middleware with strategy
func NewJWTMiddleware(strategy *strategies.JWTStrategy) *JWTMiddleware {
	return &JWTMiddleware{
		strategy: strategy,
	}
}

// Handle is the middleware handler function
// This uses the JWT strategy to authenticate requests
func (m *JWTMiddleware) Handle() contracts.MiddlewareFunc {
	return func(next contracts.HandlerFunc) contracts.HandlerFunc {
		return func(c contracts.Context) error {
			// Authenticate using JWT strategy
			userID, username, roles, err := m.strategy.Authenticate(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, &models.MessageResponse{
					Success: false,
					Message: "Unauthorized: " + err.Error(),
				})
			}

			// Store user info in context
			c.Set("userID", userID)
			c.Set("username", username)
			c.Set("roles", roles)

			// Store in request context for non-handler access
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, "userID", userID)
			ctx = context.WithValue(ctx, "username", username)
			ctx = context.WithValue(ctx, "roles", roles)

			return next(c)
		}
	}
}

// HandleWithSkipPaths returns a middleware that skips authentication for certain paths
func HandleWithSkipPaths(strategy *strategies.JWTStrategy, skipPaths []string) contracts.MiddlewareFunc {
	return func(next contracts.HandlerFunc) contracts.HandlerFunc {
		return func(c contracts.Context) error {
			// Check if path should skip authentication
			path := c.Request().URL.Path
			for _, skipPath := range skipPaths {
				if strings.HasPrefix(path, skipPath) {
					return next(c)
				}
			}

			// Authenticate using JWT strategy
			userID, username, roles, err := strategy.Authenticate(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, &models.MessageResponse{
					Success: false,
					Message: "Unauthorized: " + err.Error(),
				})
			}

			// Store user info in context
			c.Set("userID", userID)
			c.Set("username", username)
			c.Set("roles", roles)

			// Store in request context for non-handler access
			ctx := c.Request().Context()
			ctx = context.WithValue(ctx, "userID", userID)
			ctx = context.WithValue(ctx, "username", username)
			ctx = context.WithValue(ctx, "roles", roles)

			return next(c)
		}
	}
}
