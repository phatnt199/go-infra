package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/models"
	"github.com/phatnt199/go-infra/pkg/logger"
)

// ChangePasswordHandler handles password change
type ChangePasswordHandler struct {
	authService authContracts.IAuthService
	validator   *validator.Validate
	logger      logger.Logger
}

// NewChangePasswordHandler creates a new change password handler
func NewChangePasswordHandler(authService authContracts.IAuthService, validator *validator.Validate, logger logger.Logger) *ChangePasswordHandler {
	return &ChangePasswordHandler{
		authService: authService,
		validator:   validator,
		logger:      logger,
	}
}

// Handle processes change password request
// @Summary Change password
// @Description Change user password using credential pattern
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Change password request"
// @Success 200 {object} models.MessageResponse "Success"
// @Failure 400 {object} models.MessageResponse "Invalid request"
// @Failure 401 {object} models.MessageResponse "Unauthorized"
// @Router /api/v1/auth/change-password [post]
func (h *ChangePasswordHandler) Handle(c contracts.Context) error {
	ctx := c.Request().Context()

	// Get user ID from context (set by middleware)
	userID := c.Get("userID")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, &models.MessageResponse{
			Success: false,
			Message: "Unauthorized",
		})
	}

	var req models.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.MessageResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// Set user ID from context
	req.UserID = userID.(string)

	if h.validator != nil {
		if err := h.validator.Struct(req); err != nil {
			return c.JSON(http.StatusBadRequest, &models.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}
	}

	// Change password
	err := h.authService.ChangePassword(ctx, &req)
	if err != nil {
		h.logger.Warnf("Failed to change password: %v", err)
		return c.JSON(http.StatusBadRequest, &models.MessageResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, &models.MessageResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}
