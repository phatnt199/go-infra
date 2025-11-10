package authentication

import (
	"github.com/go-playground/validator/v10"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/config"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/handlers"
	"github.com/phatnt199/go-infra/pkg/component/authentication/middleware"
	"github.com/phatnt199/go-infra/pkg/component/authentication/services"
	"github.com/phatnt199/go-infra/pkg/logger"
	"github.com/phatnt199/go-infra/pkg/logger/empty"
)

// AuthenticationPlugin represents the authentication component
type AuthenticationPlugin struct {
	config         *config.AuthConfig
	logger         logger.Logger
	tokenService   authContracts.ITokenService
	authService    authContracts.IAuthService
	userProvider   authContracts.IUserProvider
	passwordHasher authContracts.IPasswordHasher
	validator      *validator.Validate
	middleware     *middleware.JWTMiddleware
}

// NewAuthenticationPlugin creates a new authentication plugin
func NewAuthenticationPlugin(opts ...config.AuthOption) *AuthenticationPlugin {
	config := config.DefaultAuthConfig()

	// Apply options
	for _, opt := range opts {
		opt(config)
	}

	plugin := &AuthenticationPlugin{
		config:    config,
		logger:    empty.EmptyLogger,
		validator: validator.New(),
	}

	plugin.initializeDefaults()
	return plugin
}

// initializeDefaults initializes default implementations if not provided
func (p *AuthenticationPlugin) initializeDefaults() {
	// Initialize token service if not provided
	if p.config.TokenService == nil {
		p.tokenService = services.NewTokenService(
			p.config.JWTSecret,
			p.config.JWTIssuer,
			p.config.JWTAudience,
			p.config.JWTAlgorithm,
			p.config.AccessTokenExpiry,
			p.config.RefreshTokenExpiry,
		)
	} else {
		p.tokenService = p.config.TokenService.(authContracts.ITokenService)
	}

	// Initialize password hasher if not provided
	if p.config.PasswordHasher == nil {
		p.passwordHasher = config.NewPasswordHasher(p.config.PasswordHashAlgorithm)
	} else {
		p.passwordHasher = p.config.PasswordHasher.(authContracts.IPasswordHasher)
	}

	// Get user provider
	if p.config.UserProvider == nil {
		p.logger.Fatal("UserProvider must be provided")
	}
	p.userProvider = p.config.UserProvider.(authContracts.IUserProvider)

	// Initialize auth service if not provided
	if p.config.AuthService == nil {
		p.authService = services.NewAuthService(
			p.userProvider,
			p.tokenService,
			p.passwordHasher,
			p.config,
			p.logger,
		)
	} else {
		p.authService = p.config.AuthService.(authContracts.IAuthService)
	}

	// Initialize middleware
	p.middleware = middleware.NewJWTMiddleware(p.tokenService, p.config.PublicPaths)
}

// Register registers the auth plugin with an HTTP server
func (p *AuthenticationPlugin) Register(server contracts.HttpServer) error {
	// Register global middleware if enabled
	// if p.config.EnableMiddleware {
	// 	server.AddMiddlewares(p.middleware.Handle())
	// }

	// Register handlers if enabled
	if p.config.EnableDefaultHandlers {
		p.registerHandlers(server)
	}

	return nil
}

// RegisterGroup registers auth endpoints on a specific route group
func (p *AuthenticationPlugin) RegisterGroup(group contracts.RouteGroup) {
	p.registerHandlersOnGroup(group)
}

// registerHandlers registers default auth handlers on the server
func (p *AuthenticationPlugin) registerHandlers(server contracts.HttpServer) {
	// Get the RouteBuilder from the server and create a group with the base path
	basePath := p.config.BasePath
	if basePath == "" || basePath == "/" {
		basePath = ""
	}

	// Create a route group with the base path
	group := server.RouteBuilder().Group(basePath)
	p.registerHandlersOnGroup(group)
}

// registerHandlersOnGroup registers handlers on a route group
func (p *AuthenticationPlugin) registerHandlersOnGroup(group contracts.RouteGroup) {
	// Sign in (public)
	signInHandler := handlers.NewSignInHandler(p.authService, p.validator, p.logger)
	group.POST("/signin", signInHandler.Handle)

	// Sign up (public)
	signUpHandler := handlers.NewSignUpHandler(p.authService, p.validator, p.logger)
	group.POST("/signup", signUpHandler.Handle)

	// Change password (protected)
	changePasswordHandler := handlers.NewChangePasswordHandler(p.authService, p.validator, p.logger)
	group.POST("/change-password", changePasswordHandler.Handle, p.RequireAuth())

	// Profile (protected)
	profileHandler := handlers.NewProfileHandler(p.authService, p.logger)
	group.GET("/profile", profileHandler.GetProfile, p.RequireAuth())
}

// RequireAuth returns the JWT authentication middleware
func (p *AuthenticationPlugin) RequireAuth() contracts.MiddlewareFunc {
	return p.middleware.Handle()
}

// GetAuthService returns the auth service
func (p *AuthenticationPlugin) GetAuthService() authContracts.IAuthService {
	return p.authService
}

// GetTokenService returns the token service
func (p *AuthenticationPlugin) GetTokenService() authContracts.ITokenService {
	return p.tokenService
}

// GetUserProvider returns the user provider
func (p *AuthenticationPlugin) GetUserProvider() authContracts.IUserProvider {
	return p.userProvider
}

// GetPasswordHasher returns the password hasher
func (p *AuthenticationPlugin) GetPasswordHasher() authContracts.IPasswordHasher {
	return p.passwordHasher
}

// GetConfig returns the auth config
func (p *AuthenticationPlugin) GetConfig() *config.AuthConfig {
	return p.config
}

// SetLogger sets the logger for the plugin
func (p *AuthenticationPlugin) SetLogger(logger logger.Logger) {
	p.logger = logger
}
