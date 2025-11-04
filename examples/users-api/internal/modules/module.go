package modules

import (
	"github.com/phatnt199/go-infra/examples/users-api/internal/handler"
	"github.com/phatnt199/go-infra/examples/users-api/internal/repository"
	"github.com/phatnt199/go-infra/examples/users-api/internal/routes"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/logger"

	"go.uber.org/fx"
)

// Module provides all dependencies for the application using fx
var Module = fx.Module(
	"users_api_module",
	fx.Provide(repository.NewUserRepository),
	fx.Provide(handler.NewUserHandler),
	fx.Invoke(setupRoutes),
)

// setupRoutes sets up all HTTP routes
func setupRoutes(
	server contracts.HttpServer,
	userHandler *handler.UserHandler,
	logger logger.Logger,
) {
	routes.SetupRoutes(server, userHandler, logger)
}
