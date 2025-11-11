package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/phatnt199/go-infra/examples/authentication-service/config"
	"github.com/phatnt199/go-infra/examples/authentication-service/docs"
	microserviceauth "github.com/phatnt199/go-infra/examples/authentication-service/internal/microservice-auth"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/shared/data"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/shared/handlers"
	"github.com/phatnt199/go-infra/pkg/adapter/fxapp"
	"github.com/phatnt199/go-infra/pkg/adapter/fxapp/contracts"
	httpContracts "github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	fiberAdapter "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
	authComponent "github.com/phatnt199/go-infra/pkg/component/authentication"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
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
	// builder.ProvideModule(auth.Module)
	// builder.ProvideModule(customauth.Module)
	builder.ProvideModule(microserviceauth.Module)

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
	registerAuthentication := func(
		authComp *authComponent.Component,
		authService authContracts.IAuthService,
		server httpContracts.HttpServer,
		log logger.Logger,
	) {
		// Create the auth route group with the basePath
		basePath := server.Cfg().GetBasePath()
		v1Group := server.RouteBuilder().Group(basePath)

		// Get JWT middleware from component for protected routes
		jwtMiddleware := authComp.GetJWTMiddleware()

		// Create local handlers with swagger annotations
		authHandlers := handlers.NewAuthHandlers(authService, validator.New(), log)
		protectedHandlers := handlers.NewProtectedHandlers(log)

		// ============ PUBLIC AUTH ENDPOINTS ============
		authGroup := v1Group.Group("/auth")

		// Sign Up: POST /api/v1/auth/signup (public)
		authGroup.POST("/signup", authHandlers.SignUp)

		// Sign In: POST /api/v1/auth/signin (public)
		authGroup.POST("/signin", authHandlers.SignIn)

		// ============ PROTECTED AUTH ENDPOINTS ============
		// These require authentication via JWT token

		// Change Password: PUT /api/v1/auth/change-password (authenticated)
		authGroup.PUT("/change-password", authHandlers.ChangePassword, jwtMiddleware.Handle())

		// Get Profile: GET /api/v1/auth/profile (authenticated)
		authGroup.GET("/profile", authHandlers.GetProfile, jwtMiddleware.Handle())

		// Update Profile: PUT /api/v1/auth/profile (authenticated)
		authGroup.PUT("/profile", authHandlers.UpdateProfile, jwtMiddleware.Handle())

		// ============ PROTECTED FEATURE ENDPOINTS ============
		// Example endpoints demonstrating authentication usage
		protectedGroup := v1Group.Group("/protected", jwtMiddleware.Handle())

		// Dashboard: GET /api/v1/protected/dashboard (authenticated)
		protectedGroup.GET("/dashboard", protectedHandlers.GetUserDashboard)

		// Settings: GET /api/v1/protected/settings (authenticated)
		protectedGroup.GET("/settings", protectedHandlers.GetUserSettings)

		// Update Settings: PUT /api/v1/protected/settings (authenticated)
		protectedGroup.PUT("/settings", protectedHandlers.UpdateUserSettings)

		// ============ ADMIN ENDPOINTS ============
		// Example endpoints requiring admin role
		adminGroup := v1Group.Group("/admin", jwtMiddleware.Handle())

		// Admin Dashboard: GET /api/v1/admin/dashboard (authenticated + admin role)
		adminGroup.GET("/dashboard", protectedHandlers.AdminOnlyEndpoint)

		log.Info("Authentication and protected endpoints registered successfully")
		log.Info("Public endpoints: /api/v1/auth/signup, /api/v1/auth/signin")
		log.Info("Protected endpoints: /api/v1/auth/*, /api/v1/protected/*, /api/v1/admin/*")
	}
	app.RegisterHook(registerAuthentication)

	return &Application{
		Application: app,
	}
}

func (a *Application) Run() {
	a.Application.Run()
}
