package handlers

import (
	"net/http"

	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/models"
	"github.com/phatnt199/go-infra/pkg/logger"
)

// ProfileHandler handles profile operations
type ProfileHandler struct {
	authService authContracts.IAuthService
	logger      logger.Logger
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(authService authContracts.IAuthService, logger logger.Logger) *ProfileHandler {
	return &ProfileHandler{
		authService: authService,
		logger:      logger,
	}
}

// GetProfile gets user profile
// @Summary Get profile
// @Description Get current user profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse "Success"
// @Failure 401 {object} models.MessageResponse "Unauthorized"
// @Router /api/v1/auth/profile [get]
func (h *ProfileHandler) GetProfile(c contracts.Context) error {
	ctx := c.Request().Context()

	// Get user ID from context (set by middleware)
	userID := c.Get("userID")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, &models.MessageResponse{
			Success: false,
			Message: "Unauthorized",
		})
	}

	profile, err := h.authService.GetProfile(ctx, userID.(string))
	if err != nil {
		h.logger.Warnf("Failed to get profile: %v", err)
		return c.JSON(http.StatusNotFound, &models.MessageResponse{
			Success: false,
			Message: "Profile not found",
		})
	}

	return c.JSON(http.StatusOK, profile)
}
