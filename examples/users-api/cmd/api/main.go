package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	fxapp "github.com/phatnt199/go-infra/pkg/adapter/fxapp"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	customfiber "github.com/phatnt199/go-infra/pkg/adapter/http/fiber_adapter"
	postgresgorm "github.com/phatnt199/go-infra/pkg/infra/postgres/gorm"
	"github.com/phatnt199/go-infra/pkg/logger"

	"github.com/phatnt199/go-infra/examples/users-api/internal/modules"

	_ "github.com/phatnt199/go-infra/examples/users-api/docs" // Import swagger docs
)

// @title Users API
// @version 1.0
// @description This is a sample users API built with go-infra
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load .env file (optional - will use system env if not found)
	// _ = godotenv.Load("../../.env")

	// Create application builder with fx dependency injection
	appBuilder := fxapp.NewApplicationBuilder()
	appLogger := appBuilder.Logger()

	appLogger.Info("Starting Users API application with fx framework")

	// Register Fiber HTTP server module
	appBuilder.ProvideModule(customfiber.Module)

	// Register PostgreSQL database module
	appBuilder.ProvideModule(postgresgorm.Module)

	// Register custom module that provides repositories, handlers, and routes
	appBuilder.ProvideModule(modules.Module)

	// Build the application
	app := appBuilder.Build()

	// Register custom Swagger setup hook
	app.RegisterHook(setupSwagger)

	// Run application
	app.Run()
}

// setupSwagger configures Swagger documentation endpoint
func setupSwagger(
	server contracts.HttpServer,
	logger logger.Logger,
) {
	server.RouteBuilder().RegisterHandler(func(instance interface{}) {
		if app, ok := instance.(*fiber.App); ok {
			app.Get("/swagger/*", swagger.HandlerDefault)
			logger.Info("Swagger documentation available at /swagger/index.html")
		}
	})
}
