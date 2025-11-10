package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/models"
	"github.com/phatnt199/go-infra/pkg/logger"
)

// SignUpHandler handles user sign up
type SignUpHandler struct {
	authService authContracts.IAuthService
	validator   *validator.Validate
	logger      logger.Logger
}

// NewSignUpHandler creates a new sign up handler
func NewSignUpHandler(authService authContracts.IAuthService, validator *validator.Validate, logger logger.Logger) *SignUpHandler {
	return &SignUpHandler{
		authService: authService,
		validator:   validator,
		logger:      logger,
	}
}

// Handle processes sign up request
// @Summary Sign up user
// @Description Register a new user with identifier/credential pattern
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.SignUpRequest true "Sign up request"
// @Success 201 {object} models.AuthResponse "Success"
// @Failure 400 {object} models.MessageResponse "Invalid request"
// @Failure 409 {object} models.MessageResponse "User already exists"
// @Router /api/v1/auth/signup [post]
func (h *SignUpHandler) Handle(c contracts.Context) error {
	ctx := c.Request().Context()
	fmt.Println("HERE.............")

	var req models.SignUpRequest
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

	// Create user
	resp, err := h.authService.SignUp(ctx, &req)
	if err != nil {
		h.logger.Warnf("Failed to sign up: %v", err)
		return c.JSON(http.StatusBadRequest, &models.MessageResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, resp)
}
