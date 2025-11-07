package auth

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/dtos/v1/requests"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/dtos/v1/responses"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/services"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/core/web/route"
	"github.com/phatnt199/go-infra/pkg/logger"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/fx"
)

type ChangePasswordHandlerParams struct {
	fx.In

	Logger      logger.Logger
	AuthGroup   contracts.RouteGroup `name:"auth-routes"`
	Validator   *validator.Validate
	UserService *services.UserService
}

type changePasswordHandler struct {
	ChangePasswordHandlerParams
}

// NewChangePasswordHandler creates a new change password handler
func NewChangePasswordHandler(params ChangePasswordHandlerParams) route.Endpoint {
	return &changePasswordHandler{
		ChangePasswordHandlerParams: params,
	}
}

func (h *changePasswordHandler) MapEndpoint() {
	h.AuthGroup.POST("/change-password/:id", h.handler())
}

// @Summary Change user password
// @Description Change password for authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body requests.ChangePasswordRequest true "Change password request"
// @Success 200 {object} responses.MessageResponse
// @Failure 400 {object} responses.MessageResponse
// @Router /api/v1/auth/change-password/{id} [post]
func (h *changePasswordHandler) handler() contracts.HandlerFunc {
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

		var req requests.ChangePasswordRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, &responses.MessageResponse{
				Success: false,
				Message: "Invalid request body",
			})
		}

		if err := h.Validator.Struct(req); err != nil {
			return c.JSON(http.StatusBadRequest, &responses.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}

		// Convert to service request
		req.UserID = userID
		serviceReq := &services.ChangePasswordRequest{
			UserID:      userID,
			OldPassword: req.GetOldPassword(),
			NewPassword: req.GetNewPassword(),
		}

		err = h.UserService.ChangePassword(ctx, serviceReq)
		if err != nil {
			h.Logger.Errorf("Failed to change password: %v", err)
			return c.JSON(http.StatusBadRequest, &responses.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, &responses.MessageResponse{
			Success: true,
			Message: "Password changed successfully",
		})
	}
}
