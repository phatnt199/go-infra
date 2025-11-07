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

type GetUserHandlerParams struct {
	fx.In

	Logger      logger.Logger
	UserGroup   contracts.RouteGroup `name:"user-routes"`
	Validator   *validator.Validate
	UserService *services.UserService
}

type getUserHandler struct {
	GetUserHandlerParams
}

// NewGetUserHandler creates a new get user handler
func NewGetUserHandler(params GetUserHandlerParams) route.Endpoint {
	return &getUserHandler{
		GetUserHandlerParams: params,
	}
}

func (h *getUserHandler) MapEndpoint() {
	h.UserGroup.GET("/:id", h.handler())
}

// @Summary Get user details
// @Description Get full user details (without password)
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} responses.UserFullDetailsResponse
// @Failure 400 {object} responses.MessageResponse
// @Failure 404 {object} responses.MessageResponse
// @Router /api/v1/users/{id} [get]
func (h *getUserHandler) handler() contracts.HandlerFunc {
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

		result, err := h.UserService.GetUserFullDetails(ctx, userID)
		if err != nil {
			h.Logger.Errorf("Failed to get user: %v", err)
			return c.JSON(http.StatusNotFound, &responses.MessageResponse{
				Success: false,
				Message: "User not found",
			})
		}

		// Convert to response DTO
		response := &responses.UserFullDetailsResponse{
			User: &responses.UserResponse{
				ID:          result.User.ID,
				Status:      result.User.Status,
				UserType:    string(result.User.UserType),
				ActivatedAt: result.User.ActivatedAt,
				LastLoginAt: result.User.LastLoginAt,
				ValidFrom:   result.User.ValidFrom,
				ValidTo:     result.User.ValidTo,
				CreatedAt:   result.User.CreatedAt,
				ModifiedAt:  result.User.ModifiedAt,
			},
		}

		if result.Identifier != nil {
			response.Identifier = &responses.UserIdentifierResponse{
				ID:         result.Identifier.ID,
				Scheme:     result.Identifier.Scheme,
				Identifier: result.Identifier.Identifier,
				Verified:   result.Identifier.Verified,
			}
		}

		if result.Profile != nil {
			response.Profile = &responses.UserProfileResponse{
				ID:         result.Profile.ID,
				Firstname:  result.Profile.Firstname,
				Lastname:   result.Profile.Lastname,
				Birthday:   result.Profile.Birthday,
				Locale:     result.Profile.Locale,
				Details:    result.Profile.Details,
				CreatedAt:  result.Profile.CreatedAt,
				ModifiedAt: result.Profile.ModifiedAt,
			}
		}

		return c.JSON(http.StatusOK, response)
	}
}
