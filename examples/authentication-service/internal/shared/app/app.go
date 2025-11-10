package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/phatnt199/go-infra/examples/authentication-service/docs"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/auth"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/auth/models"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/shared/config"
	"github.com/phatnt199/go-infra/pkg/adapter/fxapp"
	"github.com/phatnt199/go-infra/pkg/adapter/fxapp/contracts"
	httpContracts "github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	fiberAdapter "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
	"github.com/phatnt199/go-infra/pkg/component/authentication"
	"github.com/phatnt199/go-infra/pkg/health"
	"github.com/phatnt199/go-infra/pkg/logger"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type Application struct {
	contracts.Application
}

func NewApplication() *Application {
	// Load custom configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Create application builder with base fx modules
	builder := fxapp.NewApplicationBuilder()

	// Provide custom configuration
	builder.ProvideModule(fx.Module(
		"configfx",
		fx.Provide(func() *config.Config {
			return cfg
		}),
	))

	// Provide database
	builder.ProvideModule(fx.Module(
		"dbfx",
		fx.Provide(func() (*gorm.DB, error) {
			dsn := fmt.Sprintf(
				"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
				cfg.Database.Host,
				cfg.Database.Port,
				cfg.Database.User,
				cfg.Database.Password,
				cfg.Database.Name,
				cfg.Database.SSLMode,
			)

			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				return nil, fmt.Errorf("failed to connect database: %w", err)
			}

			sqlDB, err := db.DB()
			if err != nil {
				return nil, err
			}

			sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
			sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)

			// Auto migrate
			if err := db.AutoMigrate(&models.User{}); err != nil {
				return nil, fmt.Errorf("failed to migrate database: %w", err)
			}

			return db, nil
		}),
	))

	// Include Fiber HTTP server module
	builder.ProvideModule(fiberAdapter.Module)

	// Include health module
	builder.ProvideModule(health.Module)

	// Include auth module
	builder.ProvideModule(auth.Module)

	// Build the application
	app := builder.Build()

	// Register swagger docs handler
	registerSwagger := func(server httpContracts.HttpServer, log logger.Logger) {
		// Get the Fiber app instance
		fiberApp, ok := server.GetServerInstance().(*fiber.App)
		if !ok {
			log.Warn("Failed to get Fiber app instance for swagger setup")
			return
		}

		// Import the docs package to register swagger
		_ = docs.SwaggerInfo

		// Configure swagger info
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Title = "Authentication Service API"
		docs.SwaggerInfo.Description = "Authentication service using identifier/credential pattern"

		// Ensure swagger uses same-origin by clearing Host and setting BasePath to empty
		// since the paths in docs.go already include the full base path
		docs.SwaggerInfo.Host = ""
		docs.SwaggerInfo.BasePath = "" // Register swagger handler
		fiberApp.Get("/swagger/*", fiberSwagger.WrapHandler)
		server.SetupDefaultMiddlewares()
		log.Info("Swagger documentation available at http://localhost:8081/swagger")
	}
	app.RegisterHook(registerSwagger)

	// Register authentication endpoints as a hook
	registerAuthentication := func(authenticationPlugin *authentication.AuthenticationPlugin, server httpContracts.HttpServer, log logger.Logger) {
		authenticationPlugin.SetLogger(log)

		// Create the auth route group with the basePath
		basePath := server.Cfg().GetBasePath()
		v1Group := server.RouteBuilder().Group(basePath)
		authGroup := v1Group.Group("/auth")

		// Register authentication routes on the group
		authenticationPlugin.RegisterGroup(authGroup)

		log.Info("Authentication endpoints registered successfully")
	}
	app.RegisterHook(registerAuthentication)

	return &Application{
		Application: app,
	}
}

func (a *Application) Run() {
	a.Application.Run()
}
