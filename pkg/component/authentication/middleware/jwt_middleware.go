package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/models"
)

// JWTMiddleware handles JWT authentication
type JWTMiddleware struct {
	tokenService authContracts.ITokenService
	publicPaths  []string
}

// NewJWTMiddleware creates a new JWT middleware
func NewJWTMiddleware(tokenService authContracts.ITokenService, publicPaths []string) *JWTMiddleware {
	return &JWTMiddleware{
		tokenService: tokenService,
		publicPaths:  publicPaths,
	}
}

// Handle is the middleware handler function
func (m *JWTMiddleware) Handle() contracts.MiddlewareFunc {
	return func(next contracts.HandlerFunc) contracts.HandlerFunc {
		return func(c contracts.Context) error {
			// Check if path is public
			path := c.Request().URL.Path
			if m.isPublicPath(path) {
				return next(c)
			}

			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, &models.MessageResponse{
					Success: false,
					Message: "Missing authorization header",
				})
			}

			// Extract bearer token
			token, err := extractBearerToken(authHeader)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, &models.MessageResponse{
					Success: false,
					Message: "Invalid authorization header format",
				})
			}

			// Validate token
			userID, username, roles, err := m.tokenService.ValidateAccessToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, &models.MessageResponse{
					Success: false,
					Message: "Invalid or expired token",
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

// isPublicPath checks if the path is in the public paths list
func (m *JWTMiddleware) isPublicPath(path string) bool {
	for _, publicPath := range m.publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}
	return false
}

// extractBearerToken extracts the token from the Authorization header
func extractBearerToken(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", http.ErrNotSupported
	}
	return parts[1], nil
}
