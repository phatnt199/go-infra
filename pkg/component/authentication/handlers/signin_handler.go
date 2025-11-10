package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/models"
	"github.com/phatnt199/go-infra/pkg/logger"
)

// SignInHandler handles user sign in
type SignInHandler struct {
	authService authContracts.IAuthService
	validator   *validator.Validate
	logger      logger.Logger
}

// NewSignInHandler creates a new sign in handler
func NewSignInHandler(authService authContracts.IAuthService, validator *validator.Validate, logger logger.Logger) *SignInHandler {
	return &SignInHandler{
		authService: authService,
		validator:   validator,
		logger:      logger,
	}
}

// Handle processes sign in request
// @Summary Sign in user
// @Description Authenticate with identifier/credential pattern
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.SignInRequest true "Sign in request"
// @Success 200 {object} models.AuthResponse "Success"
// @Failure 400 {object} models.MessageResponse "Invalid request"
// @Failure 401 {object} models.MessageResponse "Invalid credentials"
// @Router /api/v1/auth/signin [post]
func (h *SignInHandler) Handle(c contracts.Context) error {
	ctx := c.Request().Context()

	var req models.SignInRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.MessageResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if h.validator != nil {
		if err := h.validator.Struct(req); err != nil {
			return c.JSON(http.StatusBadRequest, &models.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}
	}

	// Authenticate user
	resp, err := h.authService.SignIn(ctx, &req)
	if err != nil {
		h.logger.Warnf("Failed to sign in: %v", err)
		return c.JSON(http.StatusUnauthorized, &models.MessageResponse{
			Success: false,
			Message: "Invalid credentials",
		})
	}

	return c.JSON(http.StatusOK, resp)
}
