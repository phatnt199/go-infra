package strategies

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/config"
)

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID   string   `json:"userId"`
	Username string   `json:"username"`
	Roles    []string `json:"roles,omitempty"`
	ClientID string   `json:"clientId,omitempty"`
	Type     string   `json:"type"` // access or refresh
	jwt.RegisteredClaims
}

// JWTStrategy handles JWT authentication strategy
type JWTStrategy struct {
	config *config.JWTConfig
}

// NewJWTStrategy creates a new JWT strategy
func NewJWTStrategy(cfg *config.JWTConfig) *JWTStrategy {
	return &JWTStrategy{
		config: cfg,
	}
}

// ExtractToken extracts the JWT token from the request
// Looks for token in Authorization header (Bearer scheme)
func (s *JWTStrategy) ExtractToken(c contracts.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header is required")
	}

	// Check for Bearer scheme
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid authorization header format")
	}

	scheme := strings.ToLower(parts[0])
	if scheme != "bearer" {
		return "", fmt.Errorf("authorization scheme must be Bearer")
	}

	token := parts[1]
	if token == "" {
		return "", fmt.Errorf("token is required")
	}

	return token, nil
}

// ValidateToken validates the JWT token and returns claims
func (s *JWTStrategy) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// Authenticate performs JWT authentication on the request
// Returns userID, username, roles, and any error
func (s *JWTStrategy) Authenticate(c contracts.Context) (userID, username string, roles []string, err error) {
	// Extract token from request
	tokenString, err := s.ExtractToken(c)
	if err != nil {
		return "", "", nil, err
	}

	// Validate token
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", "", nil, err
	}

	// Check token type
	if claims.Type != "access" {
		return "", "", nil, fmt.Errorf("invalid token type: expected access token")
	}

	return claims.UserID, claims.Username, claims.Roles, nil
}
