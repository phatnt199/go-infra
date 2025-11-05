package dtosv1

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// SignUpRequest represents user registration request
type SignUpRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=6,max=100"`
	Firstname string `json:"firstname" validate:"required,min=1,max=255"`
	Lastname  string `json:"lastname" validate:"required,min=1,max=255"`
	Birthday  string `json:"birthday,omitempty"` // ISO format: YYYY-MM-DD
	Locale    string `json:"locale,omitempty"`   // default: en_US
}

// SignInRequest represents user login request
type SignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	UserID      uuid.UUID `json:"userId" validate:"required"`
	OldPassword string    `json:"oldPassword" validate:"required"`
	NewPassword string    `json:"newPassword" validate:"required,min=6,max=100"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	UserID    uuid.UUID `json:"userId"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Status    int       `json:"status"`
	UserType  string    `json:"userType"`
	CreatedAt string    `json:"createdAt"` // ISO format
	Token     string    `json:"token,omitempty"`
	Message   string    `json:"message,omitempty"`
}

// UserIdentifierResponse represents user identifier response
type UserIdentifierResponse struct {
	ID         uuid.UUID              `json:"id"`
	UserID     uuid.UUID              `json:"userId"`
	Scheme     string                 `json:"scheme"`
	Identifier string                 `json:"identifier"`
	Verified   bool                   `json:"verified"`
	Details    map[string]interface{} `json:"details,omitempty"`
	CreatedAt  string                 `json:"createdAt"`  // ISO format
	ModifiedAt string                 `json:"modifiedAt"` // ISO format
}

// UserProfileResponse represents user profile response
type UserProfileResponse struct {
	ID         uuid.UUID              `json:"id"`
	UserID     uuid.UUID              `json:"userId"`
	Firstname  string                 `json:"firstname"`
	Lastname   string                 `json:"lastname"`
	Birthday   string                 `json:"birthday,omitempty"` // ISO format: YYYY-MM-DD
	Locale     string                 `json:"locale"`
	Details    map[string]interface{} `json:"details,omitempty"`
	CreatedAt  string                 `json:"createdAt"`  // ISO format
	ModifiedAt string                 `json:"modifiedAt"` // ISO format
}

// UserResponse represents user basic information
type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Status      int       `json:"status"`
	UserType    string    `json:"userType"`
	ActivatedAt string    `json:"activatedAt,omitempty"` // ISO format
	LastLoginAt string    `json:"lastLoginAt,omitempty"` // ISO format
	ParentID    string    `json:"parentId,omitempty"`
	ValidFrom   string    `json:"validFrom,omitempty"` // ISO format
	ValidTo     string    `json:"validTo,omitempty"`   // ISO format
	CreatedAt   string    `json:"createdAt"`           // ISO format
	ModifiedAt  string    `json:"modifiedAt"`          // ISO format
}

// UserFullDetailsResponse represents complete user information (without password)
type UserFullDetailsResponse struct {
	User       *UserResponse           `json:"user"`
	Identifier *UserIdentifierResponse `json:"identifier,omitempty"`
	Profile    *UserProfileResponse    `json:"profile,omitempty"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	Firstname string `json:"firstname,omitempty" validate:"omitempty,min=1,max=255"`
	Lastname  string `json:"lastname,omitempty" validate:"omitempty,min=1,max=255"`
	Birthday  string `json:"birthday,omitempty"` // ISO format: YYYY-MM-DD
	Locale    string `json:"locale,omitempty"`
}

// UserListResponse represents paginated user list response
type UserListResponse struct {
	Items []*UserFullDetailsResponse `json:"items"`
	Total int64                      `json:"total"`
	Page  int                        `json:"page"`
	Size  int                        `json:"size"`
}

// FormatTimeToISO formats time.Time to ISO string
func FormatTimeToISO(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

// FormatDateToISO formats date to YYYY-MM-DD
func FormatDateToISO(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}
