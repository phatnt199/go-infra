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
	"go.uber.org/fx"
)

type SignUpHandlerParams struct {
	fx.In

	Logger      logger.Logger
	AuthGroup   contracts.RouteGroup `name:"auth-routes"`
	Validator   *validator.Validate
	UserService *services.UserService
}

type signUpHandler struct {
	SignUpHandlerParams
}

// NewSignUpHandler creates a new signup handler
func NewSignUpHandler(params SignUpHandlerParams) route.Endpoint {
	return &signUpHandler{
		SignUpHandlerParams: params,
	}
}

func (h *signUpHandler) MapEndpoint() {
	h.AuthGroup.POST("/signup", h.handler())
}

// @Summary Sign up a new user
// @Description Create a new user account with username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body requests.SignUpRequest true "Signup request"
// @Success 201 {object} responses.AuthResponse
// @Failure 400 {object} responses.MessageResponse
// @Failure 500 {object} responses.MessageResponse
// @Router /api/v1/auth/signup [post]
func (h *signUpHandler) handler() contracts.HandlerFunc {
	return func(c contracts.Context) error {
		ctx := c.Request().Context()

		var req requests.SignUpRequest
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
		serviceReq := &services.SignUpRequest{
			Username:  req.Username,
			Password:  req.Password,
			Firstname: req.Firstname,
			Lastname:  req.Lastname,
			Locale:    req.Locale,
		}

		result, err := h.UserService.SignUp(ctx, serviceReq)
		if err != nil {
			h.Logger.Errorf("Failed to sign up user: %v", err)
			return c.JSON(http.StatusBadRequest, &responses.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, &responses.AuthResponse{
			UserID:       result.UserID,
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
		})
	}
}
