package auth

import (
	"time"

	"github.com/phatnt199/go-infra/examples/authentication-service/internal/auth/provider"
	appcfg "github.com/phatnt199/go-infra/examples/authentication-service/internal/shared/config"
	"github.com/phatnt199/go-infra/pkg/component/authentication"
	authConfig "github.com/phatnt199/go-infra/pkg/component/authentication/config"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// ProvideAuthenticationOptions provides authentication options for fx injection
type ProvideAuthenticationOptions struct {
	fx.In

	DB     *gorm.DB
	Config *appcfg.Config
}

// ProvideAuthenticationResult holds the results of authentication provision
type ProvideAuthenticationResult struct {
	fx.Out

	AuthenticationPlugin *authentication.AuthenticationPlugin
	AuthService          authContracts.IAuthService
	TokenService         authContracts.ITokenService
	PasswordHasher       authContracts.IPasswordHasher
	UserProvider         authContracts.IUserProvider
}

// provideAuthentication creates and configures the authentication plugin with database provider
func provideAuthentication(opts ProvideAuthenticationOptions) ProvideAuthenticationResult {
	// Create user provider
	userProvider := provider.NewUserProvider(opts.DB)

	// Create authentication plugin with configuration
	authenticationPlugin := authentication.NewAuthenticationPlugin(
		authConfig.WithJWTSecret("your-super-secret-jwt-key-change-in-production"),
		authConfig.WithJWTIssuer("authentication-service"),
		authConfig.WithJWTAudience("authentication-api"),
		authConfig.WithAccessTokenExpiry(15*time.Minute),
		authConfig.WithRefreshTokenExpiry(7*24*time.Hour),
		authConfig.WithUserProvider(userProvider),
		authConfig.WithPasswordPolicy(8, false),
		authConfig.WithBasePath("/auth"),
		// Ensure public paths include the server base path so swagger (/api/v1/...) calls are treated as public
		// authentication.WithPublicPaths([]string{
		// 	opts.Config.HTTP.BasePath + "/auth/signin",
		// 	opts.Config.HTTP.BasePath + "/auth/signup",
		// 	"/health",
		// 	"/swagger",
		// 	"/docs",
		// }),
		authConfig.WithSupportedSchemes(
			[]string{"username", "email", "phone"},
			[]string{"basic", "token"},
		),
	)

	return ProvideAuthenticationResult{
		AuthenticationPlugin: authenticationPlugin,
		AuthService:          authenticationPlugin.GetAuthService(),
		TokenService:         authenticationPlugin.GetTokenService(),
		PasswordHasher:       authenticationPlugin.GetPasswordHasher(),
		UserProvider:         userProvider,
	}
}

// Module provides the auth module with fx
var Module = fx.Module(
	"authfx",

	// Provide authentication plugin and related services
	fx.Provide(provideAuthentication),
)
