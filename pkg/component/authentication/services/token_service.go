package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/phatnt199/go-infra/pkg/component/authentication/config"
	"github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
)

// TokenService implements ITokenService
type TokenService struct {
	config *config.JWTConfig
}

// NewTokenService creates a new token service from JWT config
func NewTokenService(cfg *config.JWTConfig) contracts.ITokenService {
	return &TokenService{
		config: cfg,
	}
}

// Claims represents JWT claims (same as strategies.JWTClaims for consistency)
type Claims struct {
	UserID   string   `json:"userId"`
	Username string   `json:"username"`
	Roles    []string `json:"roles,omitempty"`
	ClientID string   `json:"clientId,omitempty"`
	Type     string   `json:"type"` // access or refresh
	jwt.RegisteredClaims
}

// GenerateAccessToken generates an access token
func (s *TokenService) GenerateAccessToken(userID, username string, roles []string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.Issuer,
			Audience:  jwt.ClaimStrings{s.config.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.AccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Secret))
}

// GenerateRefreshToken generates a refresh token
func (s *TokenService) GenerateRefreshToken(userID, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Type:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.config.Issuer,
			Audience:  jwt.ClaimStrings{s.config.Audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Secret))
}

// GenerateTokenPair generates both access and refresh tokens
func (s *TokenService) GenerateTokenPair(userID, username string, roles []string) (string, string, error) {
	accessToken, err := s.GenerateAccessToken(userID, username, roles)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.GenerateRefreshToken(userID, username)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateAccessToken validates an access token and returns claims
func (s *TokenService) ValidateAccessToken(tokenString string) (string, string, []string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return "", "", nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.Type != "access" {
			return "", "", nil, fmt.Errorf("invalid token type")
		}
		return claims.UserID, claims.Username, claims.Roles, nil
	}

	return "", "", nil, fmt.Errorf("invalid token")
}

// ValidateRefreshToken validates a refresh token
func (s *TokenService) ValidateRefreshToken(tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.Type != "refresh" {
			return "", "", fmt.Errorf("invalid token type")
		}
		return claims.UserID, claims.Username, nil
	}

	return "", "", fmt.Errorf("invalid token")
}

// GetAccessTokenExpiry returns access token expiry duration
func (s *TokenService) GetAccessTokenExpiry() time.Duration {
	return s.config.AccessExpiry
}

// GetRefreshTokenExpiry returns refresh token expiry duration
func (s *TokenService) GetRefreshTokenExpiry() time.Duration {
	return s.config.RefreshExpiry
}
