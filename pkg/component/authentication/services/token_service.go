package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
)

// tokenService implements ITokenService
type tokenService struct {
	secret             string
	issuer             string
	audience           string
	algorithm          string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

// NewTokenService creates a new token service
func NewTokenService(
	secret, issuer, audience, algorithm string,
	accessExpiry, refreshExpiry time.Duration,
) contracts.ITokenService {
	return &tokenService{
		secret:             secret,
		issuer:             issuer,
		audience:           audience,
		algorithm:          algorithm,
		accessTokenExpiry:  accessExpiry,
		refreshTokenExpiry: refreshExpiry,
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID   string   `json:"userId"`
	Username string   `json:"username"`
	Roles    []string `json:"roles,omitempty"`
	Type     string   `json:"type"` // access or refresh
	jwt.RegisteredClaims
}

// GenerateAccessToken generates an access token
func (s *tokenService) GenerateAccessToken(userID, username string, roles []string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Audience:  jwt.ClaimStrings{s.audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// GenerateRefreshToken generates a refresh token
func (s *tokenService) GenerateRefreshToken(userID, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Type:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Audience:  jwt.ClaimStrings{s.audience},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// GenerateTokenPair generates both access and refresh tokens
func (s *tokenService) GenerateTokenPair(userID, username string, roles []string) (string, string, error) {
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
func (s *tokenService) ValidateAccessToken(tokenString string) (string, string, []string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
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
func (s *tokenService) ValidateRefreshToken(tokenString string) (string, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
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
func (s *tokenService) GetAccessTokenExpiry() time.Duration {
	return s.accessTokenExpiry
}

// GetRefreshTokenExpiry returns refresh token expiry duration
func (s *tokenService) GetRefreshTokenExpiry() time.Duration {
	return s.refreshTokenExpiry
}
