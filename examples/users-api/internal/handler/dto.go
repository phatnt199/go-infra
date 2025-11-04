package handler

import (
	"github.com/google/uuid"
)

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=255"`
	Email string `json:"email" binding:"required,email"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty,min=1,max=255"`
	Email string `json:"email" binding:"omitempty,email"`
}

// UserResponse represents the response body for a user
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

// UsersListResponse represents the response body for a list of users
type UsersListResponse struct {
	Data  []UserResponse `json:"data"`
	Count int            `json:"count"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}
