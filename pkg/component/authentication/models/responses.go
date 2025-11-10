package models

import "time"

// TokenInfo represents token details
type TokenInfo struct {
	Value     string    `json:"value" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Scheme    string    `json:"scheme" example:"bearer"`
	ExpiresAt time.Time `json:"expiresAt" example:"2024-12-31T23:59:59Z"`
	Type      string    `json:"type,omitempty" example:"access"`
}

// AuthResponse DTO for authentication responses (signup/signin)
// @Description Authentication response with user info and tokens
type AuthResponse struct {
	UserID       string     `json:"userId" example:"123e4567-e89b-12d3-a456-426614174000"`
	Username     string     `json:"username,omitempty" example:"john"`
	Email        string     `json:"email,omitempty" example:"john@example.com"`
	Firstname    string     `json:"firstname,omitempty" example:"John"`
	Lastname     string     `json:"lastname,omitempty" example:"Doe"`
	Status       string     `json:"status,omitempty" example:"active"`
	UserType     string     `json:"userType,omitempty" example:"user"`
	CreatedAt    time.Time  `json:"createdAt,omitempty" example:"2024-01-01T00:00:00Z"`
	AccessToken  *TokenInfo `json:"token,omitempty"`
	RefreshToken *TokenInfo `json:"refreshToken,omitempty"`
}

// UserResponse DTO for user information
// @Description User profile information
type UserResponse struct {
	UserID    string                 `json:"userId" example:"123e4567-e89b-12d3-a456-426614174000"`
	Username  string                 `json:"username,omitempty" example:"john"`
	Email     string                 `json:"email,omitempty" example:"john@example.com"`
	Firstname string                 `json:"firstname,omitempty" example:"John"`
	Lastname  string                 `json:"lastname,omitempty" example:"Doe"`
	Birthday  *time.Time             `json:"birthday,omitempty" example:"1990-01-01T00:00:00Z"`
	Locale    string                 `json:"locale,omitempty" example:"en_US"`
	Status    string                 `json:"status,omitempty" example:"active"`
	UserType  string                 `json:"userType,omitempty" example:"user"`
	CreatedAt time.Time              `json:"createdAt,omitempty" example:"2024-01-01T00:00:00Z"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// MessageResponse represents a simple message response
// @Description Standard message response
type MessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Operation completed successfully"`
}
