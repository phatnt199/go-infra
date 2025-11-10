package authentication

import (
	"github.com/go-playground/validator/v10"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/config"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"go.uber.org/fx"
)

// ProvideAuthenticationOptions provides authentication options for fx injection
type ProvideAuthenticationOptions struct {
	fx.In

	UserProvider authContracts.IUserProvider
	Validator    *validator.Validate `optional:"true"`
}

// ProvideAuthenticationResult holds the results of authentication provision
type ProvideAuthenticationResult struct {
	fx.Out

	AuthenticationPlugin *AuthenticationPlugin
	AuthService          authContracts.IAuthService
	TokenService         authContracts.ITokenService
	PasswordHasher       authContracts.IPasswordHasher
	UserProvider         authContracts.IUserProvider
}

// provideAuthenticationPlugin creates and configures the authentication plugin
func provideAuthenticationPlugin(opts ProvideAuthenticationOptions) ProvideAuthenticationResult {
	authenticationPlugin := NewAuthenticationPlugin(
		config.WithJWTSecret("your-super-secret-jwt-key-change-in-production"),
		config.WithUserProvider(opts.UserProvider),
		config.WithPasswordPolicy(8, false),
		config.WithBasePath("/auth"),
	)

	return ProvideAuthenticationResult{
		AuthenticationPlugin: authenticationPlugin,
		AuthService:          authenticationPlugin.GetAuthService(),
		TokenService:         authenticationPlugin.GetTokenService(),
		PasswordHasher:       authenticationPlugin.GetPasswordHasher(),
		UserProvider:         authenticationPlugin.GetUserProvider(),
	}
}

// Module provides the authentication component with fx
var Module = fx.Module(
	"authenticationfx",

	// Provide validator
	fx.Provide(func() *validator.Validate {
		return validator.New()
	}),

	// Provide authentication plugin and related services
	fx.Provide(provideAuthenticationPlugin),

	// Provide auth route group
	fx.Provide(
		fx.Annotate(func(authServer contracts.HttpServer) contracts.RouteGroup {
			basePath := authServer.Cfg().GetBasePath()
			v1 := authServer.RouteBuilder().Group(basePath)
			return v1.Group("/auth")
		}, fx.ResultTags(`name:"auth-routes"`))),
)
