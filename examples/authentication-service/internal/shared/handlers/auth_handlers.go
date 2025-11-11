package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/models"
	"github.com/phatnt199/go-infra/pkg/logger"
)

// AuthHandlers provides authentication HTTP handlers
type AuthHandlers struct {
	authService authContracts.IAuthService
	validator   *validator.Validate
	logger      logger.Logger
}

// NewAuthHandlers creates new authentication handlers
func NewAuthHandlers(authService authContracts.IAuthService, validator *validator.Validate, logger logger.Logger) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
		validator:   validator,
		logger:      logger,
	}
}

// SignIn handles user sign-in
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
func (h *AuthHandlers) SignIn(c contracts.Context) error {
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

// SignUp handles user registration
// @Summary Sign up new user
// @Description Register new user with identifier/credential pattern
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.SignUpRequest true "Sign up request"
// @Success 201 {object} models.AuthResponse "Success"
// @Failure 400 {object} models.MessageResponse "Invalid request"
// @Failure 409 {object} models.MessageResponse "User already exists"
// @Router /api/v1/auth/signup [post]
func (h *AuthHandlers) SignUp(c contracts.Context) error {
	ctx := c.Request().Context()

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

	// Register user
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

// ChangePassword handles password change
// @Summary Change user password
// @Description Change password using old/new credential pattern (requires authentication)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Change password request"
// @Success 200 {object} models.MessageResponse "Success"
// @Failure 400 {object} models.MessageResponse "Invalid request"
// @Failure 401 {object} models.MessageResponse "Unauthorized"
// @Router /api/v1/auth/change-password [put]
func (h *AuthHandlers) ChangePassword(c contracts.Context) error {
	ctx := c.Request().Context()

	// Get user ID from context (set by auth middleware)
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, &models.MessageResponse{
			Success: false,
			Message: "User ID not found in context",
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
	req.UserID = userID

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

// GetProfile handles getting user profile
// @Summary Get user profile
// @Description Get authenticated user profile information (requires authentication)
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse "Success"
// @Failure 401 {object} models.MessageResponse "Unauthorized"
// @Failure 404 {object} models.MessageResponse "User not found"
// @Router /api/v1/auth/profile [get]
func (h *AuthHandlers) GetProfile(c contracts.Context) error {
	ctx := c.Request().Context()

	// Get user ID from context (set by auth middleware)
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, &models.MessageResponse{
			Success: false,
			Message: "User ID not found in context",
		})
	}

	// Get profile
	resp, err := h.authService.GetProfile(ctx, userID)
	if err != nil {
		h.logger.Warnf("Failed to get profile: %v", err)
		return c.JSON(http.StatusNotFound, &models.MessageResponse{
			Success: false,
			Message: "User not found",
		})
	}

	return c.JSON(http.StatusOK, resp)
}

// UpdateProfile handles updating user profile
// @Summary Update user profile
// @Description Update authenticated user profile information (requires authentication)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} models.UserResponse "Success"
// @Failure 400 {object} models.MessageResponse "Invalid request"
// @Failure 401 {object} models.MessageResponse "Unauthorized"
// @Router /api/v1/auth/profile [put]
func (h *AuthHandlers) UpdateProfile(c contracts.Context) error {
	ctx := c.Request().Context()

	// Get user ID from context (set by auth middleware)
	userID, ok := c.Get("userID").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, &models.MessageResponse{
			Success: false,
			Message: "User ID not found in context",
		})
	}

	var req models.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &models.MessageResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// Set user ID from context
	req.UserID = userID

	if h.validator != nil {
		if err := h.validator.Struct(req); err != nil {
			return c.JSON(http.StatusBadRequest, &models.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}
	}

	// Update profile
	resp, err := h.authService.UpdateProfile(ctx, &req)
	if err != nil {
		h.logger.Warnf("Failed to update profile: %v", err)
		return c.JSON(http.StatusBadRequest, &models.MessageResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}
