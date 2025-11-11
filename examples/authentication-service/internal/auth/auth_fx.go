package auth

import (
	"time"

	"github.com/phatnt199/go-infra/examples/authentication-service/config"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/auth/provider"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/shared/data/dbcontext"
	authComponent "github.com/phatnt199/go-infra/pkg/component/authentication"
	authConfig "github.com/phatnt199/go-infra/pkg/component/authentication/config"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"github.com/phatnt199/go-infra/pkg/logger"
	"go.uber.org/fx"
)

// ProvideAuthenticationOptions holds dependencies for authentication
type ProvideAuthenticationOptions struct {
	fx.In

	DBContext  dbcontext.AuthGormDBContext
	AuthConfig *config.AuthOptions
	Logger     logger.Logger
}

// ProvideAuthenticationResult holds the provided authentication services
type ProvideAuthenticationResult struct {
	fx.Out

	AuthComponent *authComponent.Component
	AuthService   authContracts.IAuthService
	TokenService  authContracts.ITokenService
	UserProvider  authContracts.IUserProvider
}

// provideAuthentication creates and configures the authentication component
func provideAuthentication(opts ProvideAuthenticationOptions) ProvideAuthenticationResult {
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
			RequireUpper:   false,
			RequireLower:   false,
			RequireNumber:  false,
			RequireSpecial: false,
		},
		Schemes: authConfig.SchemesConfig{
			IdentifierSchemes: []string{"username", "email", "phone"},
			CredentialSchemes: []string{"basic"},
		},
	}

	// Create authentication component with custom config
	authComp := authComponent.NewComponentWithConfig(
		userProvider,
		authCfg,
		authComponent.WithLogger(opts.Logger),
	)

	return ProvideAuthenticationResult{
		AuthComponent: authComp,
		AuthService:   authComp.GetAuthService(),
		TokenService:  authComp.GetTokenService(),
		UserProvider:  userProvider,
	}
}

// Module provides the auth module with fx
var Module = fx.Module(
	"authfx",
	fx.Provide(provideAuthentication),
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
