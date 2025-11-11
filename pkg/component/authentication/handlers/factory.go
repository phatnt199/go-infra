package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/models"
	"github.com/phatnt199/go-infra/pkg/logger"
)

// HandlerFactory provides factory methods to create auth handlers
// Users call these functions to get handler functions they can register in their routes
type HandlerFactory struct {
	authService authContracts.IAuthService
	validator   *validator.Validate
	logger      logger.Logger
}

// NewHandlerFactory creates a new handler factory
func NewHandlerFactory(authService authContracts.IAuthService, validator *validator.Validate, logger logger.Logger) *HandlerFactory {
	return &HandlerFactory{
		authService: authService,
		validator:   validator,
		logger:      logger,
	}
}

// SignIn returns a handler function for sign-in
// Users can register this in their own routes
// Example: router.POST("/auth/signin", factory.SignIn())
// @Summary Sign in user
// @Description Authenticate with identifier/credential pattern
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.SignInRequest true "Sign in request"
// @Success 200 {object} models.AuthResponse "Success"
// @Failure 400 {object} models.MessageResponse "Invalid request"
// @Failure 401 {object} models.MessageResponse "Invalid credentials"
// @Router /auth/signin [post]
func (f *HandlerFactory) SignIn() contracts.HandlerFunc {
	return func(c contracts.Context) error {
		ctx := c.Request().Context()

		var req models.SignInRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, &models.MessageResponse{
				Success: false,
				Message: "Invalid request body",
			})
		}

		if f.validator != nil {
			if err := f.validator.Struct(req); err != nil {
				return c.JSON(http.StatusBadRequest, &models.MessageResponse{
					Success: false,
					Message: err.Error(),
				})
			}
		}

		// Authenticate user
		resp, err := f.authService.SignIn(ctx, &req)
		if err != nil {
			f.logger.Warnf("Failed to sign in: %v", err)
			return c.JSON(http.StatusUnauthorized, &models.MessageResponse{
				Success: false,
				Message: "Invalid credentials",
			})
		}

		return c.JSON(http.StatusOK, resp)
	}
}

// SignUp returns a handler function for sign-up
// @Summary Sign up new user
// @Description Register new user with identifier/credential pattern
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.SignUpRequest true "Sign up request"
// @Success 201 {object} models.AuthResponse "Success"
// @Failure 400 {object} models.MessageResponse "Invalid request"
// @Failure 409 {object} models.MessageResponse "User already exists"
// @Router /auth/signup [post]
func (f *HandlerFactory) SignUp() contracts.HandlerFunc {
	return func(c contracts.Context) error {
		ctx := c.Request().Context()

		var req models.SignUpRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, &models.MessageResponse{
				Success: false,
				Message: "Invalid request body",
			})
		}

		if f.validator != nil {
			if err := f.validator.Struct(req); err != nil {
				return c.JSON(http.StatusBadRequest, &models.MessageResponse{
					Success: false,
					Message: err.Error(),
				})
			}
		}

		// Register user
		resp, err := f.authService.SignUp(ctx, &req)
		if err != nil {
			f.logger.Warnf("Failed to sign up: %v", err)
			return c.JSON(http.StatusBadRequest, &models.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, resp)
	}
}

// ChangePassword returns a handler function for change password (requires authentication)
// @Summary Change user password
// @Description Change password using old/new credential pattern
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Change password request"
// @Success 200 {object} models.MessageResponse "Success"
// @Failure 400 {object} models.MessageResponse "Invalid request"
// @Failure 401 {object} models.MessageResponse "Unauthorized"
// @Router /auth/change-password [post]
func (f *HandlerFactory) ChangePassword() contracts.HandlerFunc {
	return func(c contracts.Context) error {
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

		if f.validator != nil {
			if err := f.validator.Struct(req); err != nil {
				return c.JSON(http.StatusBadRequest, &models.MessageResponse{
					Success: false,
					Message: err.Error(),
				})
			}
		}

		// Change password
		err := f.authService.ChangePassword(ctx, &req)
		if err != nil {
			f.logger.Warnf("Failed to change password: %v", err)
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
}

// GetProfile returns a handler function for getting user profile (requires authentication)
// @Summary Get user profile
// @Description Get authenticated user profile information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserResponse "Success"
// @Failure 401 {object} models.MessageResponse "Unauthorized"
// @Failure 404 {object} models.MessageResponse "User not found"
// @Router /auth/profile [get]
func (f *HandlerFactory) GetProfile() contracts.HandlerFunc {
	return func(c contracts.Context) error {
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
		resp, err := f.authService.GetProfile(ctx, userID)
		if err != nil {
			f.logger.Warnf("Failed to get profile: %v", err)
			return c.JSON(http.StatusNotFound, &models.MessageResponse{
				Success: false,
				Message: "User not found",
			})
		}

		return c.JSON(http.StatusOK, resp)
	}
}

// UpdateProfile returns a handler function for updating user profile (requires authentication)
// @Summary Update user profile
// @Description Update authenticated user profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} models.UserResponse "Success"
// @Failure 400 {object} models.MessageResponse "Invalid request"
// @Failure 401 {object} models.MessageResponse "Unauthorized"
// @Router /auth/profile [put]
func (f *HandlerFactory) UpdateProfile() contracts.HandlerFunc {
	return func(c contracts.Context) error {
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

		if f.validator != nil {
			if err := f.validator.Struct(req); err != nil {
				return c.JSON(http.StatusBadRequest, &models.MessageResponse{
					Success: false,
					Message: err.Error(),
				})
			}
		}

		// Update profile
		resp, err := f.authService.UpdateProfile(ctx, &req)
		if err != nil {
			f.logger.Warnf("Failed to update profile: %v", err)
			return c.JSON(http.StatusBadRequest, &models.MessageResponse{
				Success: false,
				Message: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, resp)
	}
}
