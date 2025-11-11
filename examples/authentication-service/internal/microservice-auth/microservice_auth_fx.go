package microserviceauth

import (
	"time"

	"github.com/phatnt199/go-infra/examples/authentication-service/config"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/shared/data/dbcontext"
	authComponent "github.com/phatnt199/go-infra/pkg/component/authentication"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	defaultModels "github.com/phatnt199/go-infra/pkg/component/authentication/implementations/default/models"
	"github.com/phatnt199/go-infra/pkg/crypto"
	"github.com/phatnt199/go-infra/pkg/logger"
	"go.uber.org/fx"
)

// ProvideMicroserviceAuthOptions holds dependencies for microservice auth
type ProvideMicroserviceAuthOptions struct {
	fx.In

	DBContext  dbcontext.AuthGormDBContext
	AuthConfig *config.AuthOptions
	Logger     logger.Logger
}

// ProvideMicroserviceAuthResult holds the provided microservice auth services
type ProvideMicroserviceAuthResult struct {
	fx.Out

	AuthComponent *authComponent.Component
	AuthService   authContracts.IAuthService
}

// provideMicroserviceAuth creates microservice-style authentication using component's default implementation
// This demonstrates using the authentication component with its built-in default implementation
func provideMicroserviceAuth(opts ProvideMicroserviceAuthOptions) (ProvideMicroserviceAuthResult, error) {
	opts.Logger.Info("Initializing MICROSERVICE authentication using component's default implementation")

	// Create JWT config
	jwtConfig := &crypto.JWTConfig{
		Secret:             opts.AuthConfig.JWT.Secret,
		Algorithm:          crypto.AlgorithmHS256,
		Issuer:             opts.AuthConfig.JWT.Issuer,
		Audience:           opts.AuthConfig.JWT.Audience,
		AccessTokenExpiry:  parseDuration(opts.AuthConfig.JWT.AccessExpiry),
		RefreshTokenExpiry: parseDuration(opts.AuthConfig.JWT.RefreshExpiry),
	}

	// Get GORM DB from context
	db := opts.DBContext.DB()

	// Auto-migrate default models
	err := db.AutoMigrate(
		&defaultModels.User{},
		&defaultModels.UserIdentifier{},
		&defaultModels.UserCredential{},
		&defaultModels.UserProfile{},
	)
	if err != nil {
		opts.Logger.Warnf("Failed to auto-migrate auth tables: %v", err)
	}

	// Create authentication component with default implementation
	// This provides a complete ready-to-use auth system following microservice pattern
	authComp, err := authComponent.NewComponentWithDefaultImplementation(
		db,
		jwtConfig,
		authComponent.WithLogger(opts.Logger),
	)
	if err != nil {
		return ProvideMicroserviceAuthResult{}, err
	}

	// Get auth service from component
	authService := authComp.GetAuthService()

	opts.Logger.Info("Microservice authentication initialized successfully using component's default implementation")

	return ProvideMicroserviceAuthResult{
		AuthComponent: authComp,
		AuthService:   authService,
	}, nil
}

// Module provides the microservice auth module with fx
var Module = fx.Module(
	"microserviceauthfx",
	fx.Provide(provideMicroserviceAuth),
)

// parseDuration helper function
func parseDuration(durationStr string) time.Duration {
	if durationStr == "" {
		return 15 * time.Minute
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 15 * time.Minute
	}
	return duration
}
