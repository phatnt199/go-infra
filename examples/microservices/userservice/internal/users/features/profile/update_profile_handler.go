package profile

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

type UpdateProfileHandlerParams struct {
	fx.In

	Logger      logger.Logger
	UserGroup   contracts.RouteGroup `name:"user-routes"`
	Validator   *validator.Validate
	UserService *services.UserService
}

type updateProfileHandler struct {
	UpdateProfileHandlerParams
}

// NewUpdateProfileHandler creates a new update profile handler
func NewUpdateProfileHandler(params UpdateProfileHandlerParams) route.Endpoint {
	return &updateProfileHandler{
		UpdateProfileHandlerParams: params,
	}
}

func (h *updateProfileHandler) MapEndpoint() {
	h.UserGroup.PUT("/:id/profile", h.handler())
}

// @Summary Update user profile
// @Description Update user profile information
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body requests.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} responses.MessageResponse
// @Failure 400 {object} responses.MessageResponse
// @Failure 404 {object} responses.MessageResponse
// @Router /api/v1/users/{id}/profile [put]
func (h *updateProfileHandler) handler() contracts.HandlerFunc {
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

		var req requests.UpdateProfileRequest
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
		serviceReq := &services.UpdateProfileRequest{
			UserID:    userID,
			Firstname: req.Firstname,
			Lastname:  req.Lastname,
			Birthday:  req.Birthday,
			Locale:    req.Locale,
		}

		err = h.UserService.UpdateProfile(ctx, serviceReq)
		if err != nil {
			h.Logger.Errorf("Failed to update profile: %v", err)
			return c.JSON(http.StatusBadRequest, &responses.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, &responses.MessageResponse{
			Success: true,
			Message: "Profile updated successfully",
		})
	}
}
