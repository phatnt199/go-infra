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

type SignInHandlerParams struct {
	fx.In

	Logger      logger.Logger
	AuthGroup   contracts.RouteGroup `name:"auth-routes"`
	Validator   *validator.Validate
	UserService *services.UserService
}

type signInHandler struct {
	SignInHandlerParams
}

// NewSignInHandler creates a new signin handler
func NewSignInHandler(params SignInHandlerParams) route.Endpoint {
	return &signInHandler{
		SignInHandlerParams: params,
	}
}

func (h *signInHandler) MapEndpoint() {
	h.AuthGroup.POST("/signin", h.handler())
}

// @Summary Sign in
// @Description Authenticate user and get access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body requests.SignInRequest true "Signin request"
// @Success 200 {object} responses.AuthResponse
// @Failure 400 {object} responses.MessageResponse
// @Failure 401 {object} responses.MessageResponse
// @Router /api/v1/auth/signin [post]
func (h *signInHandler) handler() contracts.HandlerFunc {
	return func(c contracts.Context) error {
		ctx := c.Request().Context()

		var req requests.SignInRequest
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

		// Convert DTO to service request
		serviceReq := &services.SignInRequest{
			Username: req.Username,
			Password: req.Password,
		}

		result, err := h.UserService.SignIn(ctx, serviceReq)
		if err != nil {
			h.Logger.Warnf("Failed to sign in: %v", err)
			return c.JSON(http.StatusUnauthorized, &responses.MessageResponse{
				Success: false,
				Message: "Invalid username or password",
			})
		}

		userID, _ := uuid.FromString(result.UserID)
		return c.JSON(http.StatusOK, &responses.AuthResponse{
			UserID:       userID,
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
		})
	}
}
