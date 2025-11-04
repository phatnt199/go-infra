package routes

import (
	"github.com/phatnt199/go-infra/examples/users-api/internal/handler"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/logger"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	server contracts.HttpServer,
	userHandler *handler.UserHandler,
	logger logger.Logger,
) {
	logger.Info("Setting up application routes")

	// API group with v1 versioning
	server.ConfigGroup("/api/v1", func(group contracts.RouteGroup) {
		// Users routes - CRUD operations
		// GET /api/v1/users - Get all users
		group.GET("/users", userHandler.GetAllUsers)
		// POST /api/v1/users - Create a new user
		group.POST("/users", userHandler.CreateUser)
		// GET /api/v1/users/:id - Get a user by ID
		group.GET("/users/:id", userHandler.GetUserByID)
		// PUT /api/v1/users/:id - Update a user
		group.PUT("/users/:id", userHandler.UpdateUser)
		// DELETE /api/v1/users/:id - Delete a user
		group.DELETE("/users/:id", userHandler.DeleteUser)
	})

	// Health check endpoint (outside API group)
	server.RouteBuilder().GET("/health", func(c contracts.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":  "ok",
			"service": "users-api",
		})
	})

	logger.Info("Routes setup completed")
}
