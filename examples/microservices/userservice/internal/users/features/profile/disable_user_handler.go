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

type DisableUserHandlerParams struct {
	fx.In

	Logger      logger.Logger
	UserGroup   contracts.RouteGroup `name:"user-routes"`
	Validator   *validator.Validate
	UserService *services.UserService
}

type disableUserHandler struct {
	DisableUserHandlerParams
}

// NewDisableUserHandler creates a new disable user handler
func NewDisableUserHandler(params DisableUserHandlerParams) route.Endpoint {
	return &disableUserHandler{
		DisableUserHandlerParams: params,
	}
}

func (h *disableUserHandler) MapEndpoint() {
	h.UserGroup.POST("/:id/disable", h.handler())
}

// @Summary Disable user
// @Description Deactivate a user account
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} responses.MessageResponse
// @Failure 400 {object} responses.MessageResponse
// @Failure 404 {object} responses.MessageResponse
// @Router /api/v1/users/{id}/disable [post]
func (h *disableUserHandler) handler() contracts.HandlerFunc {
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

		err = h.UserService.DisableUser(ctx, userID)
		if err != nil {
			h.Logger.Errorf("Failed to disable user: %v", err)
			return c.JSON(http.StatusBadRequest, &responses.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, &responses.MessageResponse{
			Success: true,
			Message: "User disabled successfully",
		})
	}
}
