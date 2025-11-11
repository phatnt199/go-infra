package authentication

import (
	"github.com/go-playground/validator/v10"
	"github.com/phatnt199/go-infra/pkg/component/authentication/config"
	"github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/handlers"
	defaultimpl "github.com/phatnt199/go-infra/pkg/component/authentication/implementations/default"
	"github.com/phatnt199/go-infra/pkg/component/authentication/middleware"
	"github.com/phatnt199/go-infra/pkg/component/authentication/services"
	"github.com/phatnt199/go-infra/pkg/component/authentication/strategies"
	"github.com/phatnt199/go-infra/pkg/crypto"
	"github.com/phatnt199/go-infra/pkg/logger"
	"github.com/phatnt199/go-infra/pkg/logger/empty"
	"gorm.io/gorm"
)

// Component represents the authentication component in framework mode
// This provides infrastructure (strategies, middleware, services) that users can leverage
// Users must provide their own IUserProvider implementation to connect to their user storage
type Component struct {
	config    *config.Config
	logger    logger.Logger
	validator *validator.Validate

	// Core services (can be overridden)
	tokenService   contracts.ITokenService
	passwordHasher services.IPasswordHasher
	authService    contracts.IAuthService

	// Strategy
	jwtStrategy *strategies.JWTStrategy

	// Middleware
	jwtMiddleware *middleware.JWTMiddleware

	// Handlers
	handlerFactory *handlers.HandlerFactory

	// User-provided
	userProvider contracts.IUserProvider
}

// Option is a functional option for configuring the Component
type Option func(*Component)

// NewComponent creates a new authentication component
// This is framework-first: users provide their UserProvider, we provide the infrastructure
func NewComponent(userProvider contracts.IUserProvider, opts ...Option) *Component {
	// Load config from environment (following go-infra pattern)
	cfg := config.LoadFromEnv()

	component := &Component{
		config:       cfg,
		logger:       empty.EmptyLogger,
		validator:    validator.New(),
		userProvider: userProvider,
	}

	// Apply options
	for _, opt := range opts {
		opt(component)
	}

	// Initialize infrastructure
	component.initializeDefaults()

	return component
}

// NewComponentWithConfig creates a component with custom config
func NewComponentWithConfig(userProvider contracts.IUserProvider, cfg *config.Config, opts ...Option) *Component {
	component := &Component{
		config:       cfg,
		logger:       empty.EmptyLogger,
		validator:    validator.New(),
		userProvider: userProvider,
	}

	// Apply options
	for _, opt := range opts {
		opt(component)
	}

	// Initialize infrastructure
	component.initializeDefaults()

	return component
}

// NewComponentWithDefaultImplementation creates a component with default microservice-style implementation
// This provides a ready-to-use authentication system with User/Identifier/Credential/Profile separation
// Users only need to provide a GORM database connection and JWT configuration
func NewComponentWithDefaultImplementation(db *gorm.DB, jwtConfig *crypto.JWTConfig, opts ...Option) (*Component, error) {
	// Create config from JWT config
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        jwtConfig.Secret,
			Issuer:        jwtConfig.Issuer,
			Audience:      jwtConfig.Audience,
			Algorithm:     string(jwtConfig.Algorithm),
			AccessExpiry:  jwtConfig.AccessTokenExpiry,
			RefreshExpiry: jwtConfig.RefreshTokenExpiry,
		},
		Password: config.PasswordConfig{
			HashAlgorithm:  "bcrypt",
			BcryptCost:     12,
			MinLength:      8,
			RequireUpper:   false,
			RequireLower:   false,
			RequireNumber:  false,
			RequireSpecial: false,
		},
		Schemes: config.SchemesConfig{
			IdentifierSchemes: []string{"username", "email", "phone"},
			CredentialSchemes: []string{"basic"},
		},
	}

	component := &Component{
		config:    cfg,
		logger:    empty.EmptyLogger,
		validator: validator.New(),
	}

	// Apply options first (allows overriding logger, etc.)
	for _, opt := range opts {
		opt(component)
	}

	// Create default implementation
	authService, err := defaultimpl.NewDefaultImplementation(db, jwtConfig, component.logger)
	if err != nil {
		return nil, err
	}

	// Set the auth service
	component.authService = authService

	// Initialize infrastructure (middleware, strategies, etc.)
	component.initializeDefaults()

	return component, nil
}

// initializeDefaults initializes default implementations if not provided
func (c *Component) initializeDefaults() {
	// Initialize token service if not provided
	if c.tokenService == nil {
		c.tokenService = services.NewTokenService(&c.config.JWT)
	}

	// Initialize password hasher if not provided
	if c.passwordHasher == nil {
		c.passwordHasher = services.NewBcryptHasher(c.config.Password.BcryptCost)
	}

	// Initialize JWT strategy
	c.jwtStrategy = strategies.NewJWTStrategy(&c.config.JWT)

	// Initialize auth service if not provided
	if c.authService == nil {
		c.authService = services.NewAuthService(
			c.userProvider,
			c.tokenService,
			c.passwordHasher,
			c.config,
			c.logger,
		)
	}

	// Initialize JWT middleware
	c.jwtMiddleware = middleware.NewJWTMiddleware(c.jwtStrategy)

	// Initialize handler factory
	c.handlerFactory = handlers.NewHandlerFactory(c.authService, c.validator, c.logger)
}

// GetAuthService returns the auth service for user to use in their handlers
func (c *Component) GetAuthService() contracts.IAuthService {
	return c.authService
}

// GetTokenService returns the token service
func (c *Component) GetTokenService() contracts.ITokenService {
	return c.tokenService
}

// GetJWTStrategy returns the JWT authentication strategy
// Users can use this in their own middleware
func (c *Component) GetJWTStrategy() *strategies.JWTStrategy {
	return c.jwtStrategy
}

// GetPasswordHasher returns the password hasher
func (c *Component) GetPasswordHasher() services.IPasswordHasher {
	return c.passwordHasher
}

// GetUserProvider returns the user provider
func (c *Component) GetUserProvider() contracts.IUserProvider {
	return c.userProvider
}

// GetValidator returns the validator instance
func (c *Component) GetValidator() *validator.Validate {
	return c.validator
}

// GetConfig returns the configuration
func (c *Component) GetConfig() *config.Config {
	return c.config
}

// GetLogger returns the logger
func (c *Component) GetLogger() logger.Logger {
	return c.logger
}

// GetJWTMiddleware returns the JWT middleware for protecting routes
func (c *Component) GetJWTMiddleware() *middleware.JWTMiddleware {
	return c.jwtMiddleware
}

// GetHandlerFactory returns the handler factory for creating route handlers
func (c *Component) GetHandlerFactory() *handlers.HandlerFactory {
	return c.handlerFactory
}

// Functional Options

// WithLogger sets a custom logger
func WithLogger(logger logger.Logger) Option {
	return func(c *Component) {
		c.logger = logger
	}
}

// WithTokenService sets a custom token service
func WithTokenService(service contracts.ITokenService) Option {
	return func(c *Component) {
		c.tokenService = service
	}
}

// WithPasswordHasher sets a custom password hasher
func WithPasswordHasher(hasher services.IPasswordHasher) Option {
	return func(c *Component) {
		c.passwordHasher = hasher
	}
}

// WithAuthService sets a custom auth service (full override)
func WithAuthService(service contracts.IAuthService) Option {
	return func(c *Component) {
		c.authService = service
	}
}

// WithValidator sets a custom validator
func WithValidator(validator *validator.Validate) Option {
	return func(c *Component) {
		c.validator = validator
	}
}
