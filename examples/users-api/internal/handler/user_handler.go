package handler

import (
	"github.com/google/uuid"
	"github.com/phatnt199/go-infra/examples/users-api/internal/domain"
	"github.com/phatnt199/go-infra/examples/users-api/internal/repository"
	"github.com/phatnt199/go-infra/pkg/adapter/http/contracts"
	"github.com/phatnt199/go-infra/pkg/logger"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	repo   *repository.UserRepository
	logger logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(repo *repository.UserRepository, logger logger.Logger) *UserHandler {
	return &UserHandler{
		repo:   repo,
		logger: logger,
	}
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} UsersListResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c contracts.Context) error {
	h.logger.Info("Received request to get all users")

	users, err := h.repo.GetAll(c.Request().Context())
	if err != nil {
		h.logger.Errorw("Failed to get users", logger.Fields{
			"error": err,
		})
		return c.JSON(500, ErrorResponse{
			Error:   "Failed to fetch users",
			Message: err.Error(),
		})
	}

	h.logger.Infow("Successfully retrieved users", logger.Fields{
		"count": len(users),
	})

	// Convert domain users to response
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
		})
	}

	return c.JSON(200, UsersListResponse{
		Data:  userResponses,
		Count: len(userResponses),
	})
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Get a specific user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} SuccessResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c contracts.Context) error {
	idParam := c.Param("id")
	h.logger.Infow("Received request to get user", logger.Fields{
		"id": idParam,
	})

	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Errorw("Invalid user ID format", logger.Fields{
			"id":    idParam,
			"error": err,
		})
		return c.JSON(400, ErrorResponse{
			Error:   "Invalid user ID format",
			Message: err.Error(),
		})
	}

	user, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		h.logger.Errorw("Failed to get user", logger.Fields{
			"id":    id,
			"error": err,
		})
		return c.JSON(500, ErrorResponse{
			Error:   "Failed to fetch user",
			Message: err.Error(),
		})
	}

	if user == nil {
		h.logger.Infow("User not found", logger.Fields{
			"id": id,
		})
		return c.JSON(404, ErrorResponse{
			Error:   "User not found",
			Message: "No user found with the provided ID",
		})
	}

	userResponse := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}

	return c.JSON(200, SuccessResponse{
		Data: userResponse,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided data
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "Create user request"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c contracts.Context) error {
	h.logger.Info("Received request to create a new user")

	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Errorw("Failed to parse request body", logger.Fields{
			"error": err,
		})
		return c.JSON(400, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	user := &domain.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := h.repo.Create(c.Request().Context(), user); err != nil {
		h.logger.Errorw("Failed to create user", logger.Fields{
			"email": req.Email,
			"error": err,
		})
		return c.JSON(500, ErrorResponse{
			Error:   "Failed to create user",
			Message: err.Error(),
		})
	}

	h.logger.Infow("Successfully created user", logger.Fields{
		"id":    user.ID,
		"email": user.Email,
	})

	userResponse := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}

	return c.JSON(201, SuccessResponse{
		Data:    userResponse,
		Message: "User created successfully",
	})
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update an existing user with the provided data
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body UpdateUserRequest true "Update user request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c contracts.Context) error {
	idParam := c.Param("id")
	h.logger.Infow("Received request to update user", logger.Fields{
		"id": idParam,
	})

	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Errorw("Invalid user ID format", logger.Fields{
			"id":    idParam,
			"error": err,
		})
		return c.JSON(400, ErrorResponse{
			Error:   "Invalid user ID format",
			Message: err.Error(),
		})
	}

	var req UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Errorw("Failed to parse request body", logger.Fields{
			"error": err,
		})
		return c.JSON(400, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	updates := &domain.User{
		Name:  req.Name,
		Email: req.Email,
	}

	user, err := h.repo.Update(c.Request().Context(), id, updates)
	if err != nil {
		h.logger.Errorw("Failed to update user", logger.Fields{
			"id":    id,
			"error": err,
		})
		return c.JSON(500, ErrorResponse{
			Error:   "Failed to update user",
			Message: err.Error(),
		})
	}

	if user == nil {
		h.logger.Infow("User not found for update", logger.Fields{
			"id": id,
		})
		return c.JSON(404, ErrorResponse{
			Error:   "User not found",
			Message: "No user found with the provided ID",
		})
	}

	h.logger.Infow("Successfully updated user", logger.Fields{
		"id": id,
	})

	userResponse := UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}

	return c.JSON(200, SuccessResponse{
		Data:    userResponse,
		Message: "User updated successfully",
	})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c contracts.Context) error {
	idParam := c.Param("id")
	h.logger.Infow("Received request to delete user", logger.Fields{
		"id": idParam,
	})

	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Errorw("Invalid user ID format", logger.Fields{
			"id":    idParam,
			"error": err,
		})
		return c.JSON(400, ErrorResponse{
			Error:   "Invalid user ID format",
			Message: err.Error(),
		})
	}

	if err := h.repo.Delete(c.Request().Context(), id); err != nil {
		h.logger.Errorw("Failed to delete user", logger.Fields{
			"id":    id,
			"error": err,
		})
		return c.JSON(500, ErrorResponse{
			Error:   "Failed to delete user",
			Message: err.Error(),
		})
	}

	h.logger.Infow("Successfully deleted user", logger.Fields{
		"id": id,
	})

	return c.JSON(200, SuccessResponse{
		Message: "User deleted successfully",
	})
}
