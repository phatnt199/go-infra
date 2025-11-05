package responses

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// UserResponse DTO for user information (without sensitive data)
type UserResponse struct {
	ID          uuid.UUID  `json:"id"`
	Status      int        `json:"status"`
	UserType    string     `json:"userType"`
	ActivatedAt *time.Time `json:"activatedAt,omitempty"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	ValidFrom   *time.Time `json:"validFrom,omitempty"`
	ValidTo     *time.Time `json:"validTo,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	ModifiedAt  time.Time  `json:"modifiedAt"`
}

// UserIdentifierResponse DTO for user identifier
type UserIdentifierResponse struct {
	ID         uuid.UUID `json:"id"`
	Scheme     string    `json:"scheme"`
	Identifier string    `json:"identifier"`
	Verified   bool      `json:"verified"`
}

// UserProfileResponse DTO for user profile
type UserProfileResponse struct {
	ID         uuid.UUID              `json:"id"`
	Firstname  string                 `json:"firstname"`
	Lastname   string                 `json:"lastname"`
	Birthday   *time.Time             `json:"birthday,omitempty"`
	Locale     string                 `json:"locale"`
	Details    map[string]interface{} `json:"details,omitempty"`
	CreatedAt  time.Time              `json:"createdAt"`
	ModifiedAt time.Time              `json:"modifiedAt"`
}

// UserFullDetailsResponse DTO for complete user information
type UserFullDetailsResponse struct {
	User       *UserResponse           `json:"user"`
	Identifier *UserIdentifierResponse `json:"identifier,omitempty"`
	Profile    *UserProfileResponse    `json:"profile,omitempty"`
}
