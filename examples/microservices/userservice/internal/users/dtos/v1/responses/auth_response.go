package responses

import uuid "github.com/satori/go.uuid"

// AuthResponse DTO for authentication responses (signup/signin)
type AuthResponse struct {
	UserID       uuid.UUID `json:"userId"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}
