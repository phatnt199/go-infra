package handler

import (
	"github.com/phatnt199/go-infra/examples/users-api/internal/repository"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/logger"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	repo   *repository.UserRepository
	logger logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(repo *repository.UserRepository, logger logger.Logger) *UserHandler {
	return &UserHandler{
		repo:   repo,
		logger: logger,
	}
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "users"
// @Failure 500 {object} map[string]interface{} "error"
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c contracts.Context) error {
	h.logger.Info("Received request to get all users")

	users, err := h.repo.GetAll(c.Request().Context())
	if err != nil {
		h.logger.Errorw("Failed to get users", logger.Fields{
			"error": err,
		})
		return c.JSON(500, map[string]interface{}{
			"error": "Failed to fetch users",
		})
	}

	h.logger.Infow("Successfully retrieved users", logger.Fields{
		"count": len(users),
	})

	return c.JSON(200, map[string]interface{}{
		"data":  users,
		"count": len(users),
	})
}
