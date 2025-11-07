package requests

import "time"

// UpdateProfileRequest DTO for updating user profile
type UpdateProfileRequest struct {
	Firstname string     `json:"firstname" validate:"max=100"`
	Lastname  string     `json:"lastname" validate:"max=100"`
	Birthday  *time.Time `json:"birthday"`
	Locale    string     `json:"locale" validate:"omitempty,len=5"`
}
