package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/phatnt199/go-infra/examples/authentication-service/config"
	"github.com/phatnt199/go-infra/examples/authentication-service/docs"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/auth"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/shared/data"
	"github.com/phatnt199/go-infra/pkg/adapter/fxapp"
	"github.com/phatnt199/go-infra/pkg/adapter/fxapp/contracts"
	httpContracts "github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	fiberAdapter "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
	authComponent "github.com/phatnt199/go-infra/pkg/component/authentication"
	"github.com/phatnt199/go-infra/pkg/health"
	postgresgorm "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
	"github.com/phatnt199/go-infra/pkg/logger"
	"github.com/phatnt199/go-infra/pkg/migration/goose"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

type Application struct {
	contracts.Application
}

func NewApplication() *Application {
	// Create application builder with base fx modules
	builder := fxapp.NewApplicationBuilder()

	// Include configuration module
	builder.ProvideModule(config.Module)

	// Include Fiber HTTP server module
	builder.ProvideModule(fiberAdapter.Module)

	// Include PostgreSQL GORM module
	builder.ProvideModule(postgresgorm.Module)

	// Include migration module
	builder.ProvideModule(goose.Module)

	// Include data module
	builder.ProvideModule(data.Module)

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
	registerAuthentication := func(authComp *authComponent.Component, server httpContracts.HttpServer, log logger.Logger) {
		// Create the auth route group with the basePath
		basePath := server.Cfg().GetBasePath()
		v1Group := server.RouteBuilder().Group(basePath)
		authGroup := v1Group.Group("/auth")

		// Create handler factory from component
		handlerFactory := authComp.GetHandlerFactory()

		// Get JWT middleware from component for protected routes
		jwtMiddleware := authComp.GetJWTMiddleware()

		// Register authentication endpoints using handler factory
		// Sign Up: POST /auth/signup (public)
		authGroup.POST("/signup", handlerFactory.SignUp())

		// Sign In: POST /auth/signin (public)
		authGroup.POST("/signin", handlerFactory.SignIn())

		// Change Password: PUT /auth/change-password (requires authentication)
		authGroup.PUT("/change-password", handlerFactory.ChangePassword(), jwtMiddleware.Handle())

		// Get Profile: GET /auth/profile (requires authentication)
		authGroup.GET("/profile", handlerFactory.GetProfile(), jwtMiddleware.Handle())

		// Update Profile: PUT /auth/profile (requires authentication)
		authGroup.PUT("/profile", handlerFactory.UpdateProfile(), jwtMiddleware.Handle())

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
