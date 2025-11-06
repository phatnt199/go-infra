package auth

import (
	"net/http"
	"time"

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
		result, err := h.UserService.SignIn(ctx, &req)
		if err != nil {
			h.Logger.Warnf("Failed to sign in: %v", err)
			return c.JSON(http.StatusUnauthorized, &responses.MessageResponse{
				Success: false,
				Message: "Invalid username or password",
			})
		}

		userID, _ := uuid.FromString(result.UserID)

		// Build token response structure
		resp := &responses.AuthResponse{
			UserID:    userID,
			Username:  result.Username,
			Status:    result.UserStatus,
			UserType:  result.UserType,
			CreatedAt: result.CreatedAt,
		}

		// Include access token if available
		if result.AccessToken != "" {
			resp.AccessToken = &responses.TokenInfo{
				Value:     result.AccessToken,
				Scheme:    "bearer",
				Type:      "access",
				ExpiresAt: time.Now().Add(result.AccessTokenExpiry),
			}
		}

		// Include refresh token if available
		if result.RefreshToken != "" {
			resp.RefreshToken = &responses.TokenInfo{
				Value:     result.RefreshToken,
				Scheme:    "bearer",
				Type:      "refresh",
				ExpiresAt: time.Now().Add(result.RefreshTokenExpiry),
			}
		}

		return c.JSON(http.StatusOK, resp)
	}
}
