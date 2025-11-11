package handlers

import (
	"net/http"

	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/models"
	"github.com/phatnt199/go-infra/pkg/logger"
)

// ProtectedHandlers provides protected endpoint examples
type ProtectedHandlers struct {
	logger logger.Logger
}

// NewProtectedHandlers creates new protected handlers
func NewProtectedHandlers(logger logger.Logger) *ProtectedHandlers {
	return &ProtectedHandlers{
		logger: logger,
	}
}

// GetUserDashboard returns user-specific dashboard data
// @Summary Get user dashboard
// @Description Get dashboard data for authenticated user (requires authentication)
// @Tags protected
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Dashboard data"
// @Failure 401 {object} models.MessageResponse "Unauthorized"
// @Router /api/v1/protected/dashboard [get]
func (h *ProtectedHandlers) GetUserDashboard(c contracts.Context) error {
	// Get user info from context (set by auth middleware)
	userID, _ := c.Get("userID").(string)
	username, _ := c.Get("username").(string)
	roles, _ := c.Get("roles").([]string)

	// Return dashboard data
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Welcome to your dashboard",
		"data": map[string]interface{}{
			"userId":   userID,
			"username": username,
			"roles":    roles,
			"stats": map[string]interface{}{
				"totalItems":    42,
				"activeItems":   23,
				"completedTask": 15,
			},
		},
	})
}

// GetUserSettings returns user settings
// @Summary Get user settings
// @Description Get settings for authenticated user (requires authentication)
// @Tags protected
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User settings"
// @Failure 401 {object} models.MessageResponse "Unauthorized"
// @Router /api/v1/protected/settings [get]
func (h *ProtectedHandlers) GetUserSettings(c contracts.Context) error {
	userID, _ := c.Get("userID").(string)
	username, _ := c.Get("username").(string)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"userId":   userID,
			"username": username,
			"settings": map[string]interface{}{
				"theme":         "dark",
				"language":      "en",
				"notifications": true,
				"privacy":       "public",
			},
		},
	})
}

// AdminOnlyEndpoint demonstrates role-based access (admin only)
// @Summary Admin only endpoint
// @Description Endpoint accessible only to admin users (requires authentication with admin role)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Admin data"
// @Failure 401 {object} models.MessageResponse "Unauthorized"
// @Failure 403 {object} models.MessageResponse "Forbidden - Admin role required"
// @Router /api/v1/admin/dashboard [get]
func (h *ProtectedHandlers) AdminOnlyEndpoint(c contracts.Context) error {
	userID, _ := c.Get("userID").(string)
	username, _ := c.Get("username").(string)
	roles, _ := c.Get("roles").([]string)

	// Check if user has admin role
	hasAdminRole := false
	for _, role := range roles {
		if role == "admin" {
			hasAdminRole = true
			break
		}
	}

	if !hasAdminRole {
		return c.JSON(http.StatusForbidden, &models.MessageResponse{
			Success: false,
			Message: "Admin role required",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Welcome to admin dashboard",
		"data": map[string]interface{}{
			"userId":   userID,
			"username": username,
			"roles":    roles,
			"adminStats": map[string]interface{}{
				"totalUsers":     150,
				"activeUsers":    120,
				"pendingReports": 5,
			},
		},
	})
}

// UpdateUserSettings updates user settings
// @Summary Update user settings
// @Description Update settings for authenticated user (requires authentication)
// @Tags protected
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param settings body map[string]interface{} true "Settings to update"
// @Success 200 {object} models.MessageResponse "Settings updated"
// @Failure 400 {object} models.MessageResponse "Invalid request"
// @Failure 401 {object} models.MessageResponse "Unauthorized"
// @Router /api/v1/protected/settings [put]
func (h *ProtectedHandlers) UpdateUserSettings(c contracts.Context) error {
	userID, _ := c.Get("userID").(string)

	var settings map[string]interface{}
	if err := c.Bind(&settings); err != nil {
		return c.JSON(http.StatusBadRequest, &models.MessageResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	h.logger.Infof("User %s updating settings: %+v", userID, settings)

	return c.JSON(http.StatusOK, &models.MessageResponse{
		Success: true,
		Message: "Settings updated successfully",
	})
}
