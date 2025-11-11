package customauth

import (
	"time"

	"github.com/phatnt199/go-infra/examples/authentication-service/config"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/auth/provider"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/custom-auth/services"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/shared/data/dbcontext"
	authComponent "github.com/phatnt199/go-infra/pkg/component/authentication"
	authConfig "github.com/phatnt199/go-infra/pkg/component/authentication/config"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	authServices "github.com/phatnt199/go-infra/pkg/component/authentication/services"
	"github.com/phatnt199/go-infra/pkg/logger"
	"go.uber.org/fx"
)

// ProvideCustomAuthenticationOptions holds dependencies for custom authentication
type ProvideCustomAuthenticationOptions struct {
	fx.In

	DBContext  dbcontext.AuthGormDBContext
	AuthConfig *config.AuthOptions
	Logger     logger.Logger
}

// ProvideCustomAuthenticationResult holds the provided custom authentication services
type ProvideCustomAuthenticationResult struct {
	fx.Out

	AuthComponent       *authComponent.Component
	AuthService         authContracts.IAuthService
	CustomAuthService   *services.CustomAuthService `name:"customAuthService"`
	TokenService        authContracts.ITokenService
	UserProvider        authContracts.IUserProvider
	AuditLogger         *services.AuditLogger
	LoginAttemptManager *services.LoginAttemptManager
}

// provideCustomAuthentication creates and configures custom authentication
func provideCustomAuthentication(opts ProvideCustomAuthenticationOptions) ProvideCustomAuthenticationResult {
	opts.Logger.Info("Initializing CUSTOM authentication component")

	// Create user provider with DB context
	userProvider := provider.NewUserProvider(opts.DBContext)

	// Create auth config from options
	authCfg := &authConfig.Config{
		JWT: authConfig.JWTConfig{
			Secret:        opts.AuthConfig.JWT.Secret,
			Issuer:        opts.AuthConfig.JWT.Issuer,
			Audience:      opts.AuthConfig.JWT.Audience,
			Algorithm:     opts.AuthConfig.JWT.Algorithm,
			AccessExpiry:  parseDuration(opts.AuthConfig.JWT.AccessExpiry),
			RefreshExpiry: parseDuration(opts.AuthConfig.JWT.RefreshExpiry),
		},
		Password: authConfig.PasswordConfig{
			HashAlgorithm:  "bcrypt",
			BcryptCost:     12,
			MinLength:      8,
			RequireUpper:   true,
			RequireLower:   true,
			RequireNumber:  true,
			RequireSpecial: true,
		},
		Schemes: authConfig.SchemesConfig{
			IdentifierSchemes: []string{"username", "email", "phone"},
			CredentialSchemes: []string{"basic"},
		},
	}

	// Create token service and password hasher (reuse from component)
	tokenService := authServices.NewTokenService(&authCfg.JWT)
	passwordHasher := authServices.NewBcryptHasher(authCfg.Password.BcryptCost)

	// Create custom services
	auditLogger := services.NewAuditLogger(opts.DBContext, opts.Logger)
	loginAttemptManager := services.NewLoginAttemptManager(
		opts.DBContext,
		opts.Logger,
		5,              // max attempts
		30*time.Minute, // lock duration
	)

	// Create custom auth service
	customAuthService := services.NewCustomAuthService(
		userProvider,
		tokenService,
		passwordHasher,
		authCfg,
		opts.Logger,
		auditLogger,
		loginAttemptManager,
	)

	// Create authentication component with custom auth service
	authComp := authComponent.NewComponentWithConfig(
		userProvider,
		authCfg,
		authComponent.WithLogger(opts.Logger),
		authComponent.WithAuthService(customAuthService),
		authComponent.WithTokenService(tokenService),
		authComponent.WithPasswordHasher(passwordHasher),
	)

	opts.Logger.Info("Custom authentication component initialized successfully")

	return ProvideCustomAuthenticationResult{
		AuthComponent:       authComp,
		AuthService:         customAuthService,
		CustomAuthService:   customAuthService.(*services.CustomAuthService),
		TokenService:        tokenService,
		UserProvider:        userProvider,
		AuditLogger:         auditLogger,
		LoginAttemptManager: loginAttemptManager,
	}
}

// Module provides the custom auth module with fx
var Module = fx.Module(
	"customauthfx",
	fx.Provide(provideCustomAuthentication),
)

// parseDuration helper function
func parseDuration(durationStr string) time.Duration {
	if durationStr == "" {
		return 15 * time.Minute
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		// Return default if parsing fails
		return 15 * time.Minute
	}
	return duration
}
