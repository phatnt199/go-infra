package responses

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// TokenInfo represents token details matching node-infra token structure
type TokenInfo struct {
	Value     string    `json:"value"`
	Scheme    string    `json:"scheme"` // e.g., "bearer"
	ExpiresAt time.Time `json:"expiresAt"`
	Type      string    `json:"type,omitempty"` // e.g., "access", "refresh"
}

// AuthResponse DTO for authentication responses (signup/signin)
// Matches node-infra response structure with token and user information
type AuthResponse struct {
	UserID       uuid.UUID  `json:"userId"`
	Username     string     `json:"username,omitempty"`
	Firstname    string     `json:"firstname,omitempty"`
	Lastname     string     `json:"lastname,omitempty"`
	Status       string     `json:"status,omitempty"`   // user status
	UserType     string     `json:"userType,omitempty"` // default: "user"
	CreatedAt    time.Time  `json:"createdAt,omitempty"`
	AccessToken  *TokenInfo `json:"token,omitempty"`        // main token (for node-infra compat as "token")
	RefreshToken *TokenInfo `json:"refreshToken,omitempty"` // refresh token
}
