package profile

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/dtos/v1/responses"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/services"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/core/web/route"
	"github.com/phatnt199/go-infra/pkg/logger"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/fx"
)

type EnableUserHandlerParams struct {
	fx.In

	Logger      logger.Logger
	UserGroup   contracts.RouteGroup `name:"user-routes"`
	Validator   *validator.Validate
	UserService *services.UserService
}

type enableUserHandler struct {
	EnableUserHandlerParams
}

// NewEnableUserHandler creates a new enable user handler
func NewEnableUserHandler(params EnableUserHandlerParams) route.Endpoint {
	return &enableUserHandler{
		EnableUserHandlerParams: params,
	}
}

func (h *enableUserHandler) MapEndpoint() {
	h.UserGroup.POST("/:id/enable", h.handler())
}

// @Summary Enable user
// @Description Activate a user account
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} responses.MessageResponse
// @Failure 400 {object} responses.MessageResponse
// @Failure 404 {object} responses.MessageResponse
// @Router /api/v1/users/{id}/enable [post]
func (h *enableUserHandler) handler() contracts.HandlerFunc {
	return func(c contracts.Context) error {
		ctx := c.Request().Context()

		userIDStr := c.Param("id")
		userID, err := uuid.FromString(userIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, &responses.MessageResponse{
				Success: false,
				Message: "Invalid user ID",
			})
		}

		err = h.UserService.EnableUser(ctx, userID)
		if err != nil {
			h.Logger.Errorf("Failed to enable user: %v", err)
			return c.JSON(http.StatusBadRequest, &responses.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, &responses.MessageResponse{
			Success: true,
			Message: "User enabled successfully",
		})
	}
}
