package config

import "time"

// AuthConfig contains configuration for the auth component
type AuthConfig struct {
	// JWT Configuration
	JWTSecret          string
	JWTIssuer          string
	JWTAudience        string
	JWTAlgorithm       string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration

	// Password Configuration
	PasswordHashAlgorithm string // bcrypt, argon2, etc.
	MinPasswordLength     int
	RequireSpecialChar    bool

	// User Provider
	UserProvider interface{} // IUserProvider implementation

	// Token Service
	TokenService interface{} // ITokenService implementation

	// Password Hasher
	PasswordHasher interface{} // IPasswordHasher implementation

	// Auth Service
	AuthService interface{} // IAuthService implementation

	// Middleware
	EnableMiddleware bool
	PublicPaths      []string

	// Handlers
	EnableDefaultHandlers bool
	BasePath              string

	// Supported schemes
	SupportedIdentifierSchemes []string // username, email, phone, etc.
	SupportedCredentialSchemes []string // basic, token, etc.
}

// DefaultAuthConfig returns a default auth configuration
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTSecret:          "change-this-secret-in-production",
		JWTIssuer:          "authentication-service",
		JWTAudience:        "authentication-api",
		JWTAlgorithm:       "HS256",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,

		PasswordHashAlgorithm: "bcrypt",
		MinPasswordLength:     8,
		RequireSpecialChar:    false,

		EnableMiddleware:      true,
		EnableDefaultHandlers: true,
		BasePath:              "/auth",

		PublicPaths: []string{
			"/auth/signin",
			"/auth/signup",
			"/health",
			"/swagger",
			"/docs",
		},

		SupportedIdentifierSchemes: []string{"username", "email", "phone"},
		SupportedCredentialSchemes: []string{"basic", "token"},
	}
}

// AuthOption is a functional option for configuring AuthPlugin
type AuthOption func(*AuthConfig)

// WithJWTSecret sets the JWT secret key
func WithJWTSecret(secret string) AuthOption {
	return func(c *AuthConfig) {
		c.JWTSecret = secret
	}
}

// WithJWTIssuer sets the JWT issuer
func WithJWTIssuer(issuer string) AuthOption {
	return func(c *AuthConfig) {
		c.JWTIssuer = issuer
	}
}

// WithJWTAudience sets the JWT audience
func WithJWTAudience(audience string) AuthOption {
	return func(c *AuthConfig) {
		c.JWTAudience = audience
	}
}

// WithAccessTokenExpiry sets the access token expiry duration
func WithAccessTokenExpiry(expiry time.Duration) AuthOption {
	return func(c *AuthConfig) {
		c.AccessTokenExpiry = expiry
	}
}

// WithRefreshTokenExpiry sets the refresh token expiry duration
func WithRefreshTokenExpiry(expiry time.Duration) AuthOption {
	return func(c *AuthConfig) {
		c.RefreshTokenExpiry = expiry
	}
}

// WithPasswordPolicy sets password policy requirements
func WithPasswordPolicy(minLength int, requireSpecialChar bool) AuthOption {
	return func(c *AuthConfig) {
		c.MinPasswordLength = minLength
		c.RequireSpecialChar = requireSpecialChar
	}
}

// WithUserProvider sets a custom user provider
func WithUserProvider(provider interface{}) AuthOption {
	return func(c *AuthConfig) {
		c.UserProvider = provider
	}
}

// WithTokenService sets a custom token service
func WithTokenService(service interface{}) AuthOption {
	return func(c *AuthConfig) {
		c.TokenService = service
	}
}

// WithPasswordHasher sets a custom password hasher
func WithPasswordHasher(hasher interface{}) AuthOption {
	return func(c *AuthConfig) {
		c.PasswordHasher = hasher
	}
}

// WithAuthService sets a custom auth service
func WithAuthService(service interface{}) AuthOption {
	return func(c *AuthConfig) {
		c.AuthService = service
	}
}

// WithBasePath sets the base path for auth endpoints
func WithBasePath(path string) AuthOption {
	return func(c *AuthConfig) {
		c.BasePath = path
	}
}

// WithPublicPaths sets paths that don't require authentication
func WithPublicPaths(paths []string) AuthOption {
	return func(c *AuthConfig) {
		c.PublicPaths = paths
	}
}

// WithMiddleware enables or disables global middleware
func WithMiddleware(enabled bool) AuthOption {
	return func(c *AuthConfig) {
		c.EnableMiddleware = enabled
	}
}

// WithDefaultHandlers enables or disables default handlers
func WithDefaultHandlers(enabled bool) AuthOption {
	return func(c *AuthConfig) {
		c.EnableDefaultHandlers = enabled
	}
}

// WithSupportedSchemes sets supported identifier and credential schemes
func WithSupportedSchemes(identifierSchemes, credentialSchemes []string) AuthOption {
	return func(c *AuthConfig) {
		c.SupportedIdentifierSchemes = identifierSchemes
		c.SupportedCredentialSchemes = credentialSchemes
	}
}
