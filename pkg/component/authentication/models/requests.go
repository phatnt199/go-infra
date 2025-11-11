package models

// Credential represents authentication credential with scheme and value
type Credential struct {
	Scheme string `json:"scheme" validate:"required" example:"basic"`
	Value  string `json:"value" validate:"required" example:"pass123"`
}

// Identifier represents user identifier with scheme and value
type Identifier struct {
	Scheme string `json:"scheme" validate:"required" example:"username"`
	Value  string `json:"value" validate:"required" example:"john"`
}

// SignInRequest DTO for user authentication
// @Description Authentication request using identifier/credential pattern
// @Description Schemes must be explicitly specified
type SignInRequest struct {
	ClientID   string      `json:"clientId,omitempty" example:"web-app"`
	Identifier *Identifier `json:"identifier" validate:"required"`
	Credential *Credential `json:"credential" validate:"required"`
}

// GetIdentifierValue returns the identifier value (username, email, etc.)
func (r *SignInRequest) GetIdentifierValue() string {
	if r.Identifier != nil {
		return r.Identifier.Value
	}
	return ""
}

// GetCredentialValue returns the credential value (password, token, etc.)
func (r *SignInRequest) GetCredentialValue() string {
	if r.Credential != nil {
		return r.Credential.Value
	}
	return ""
}

// SignUpRequest DTO for user registration
// @Description User registration using identifier/credential pattern
type SignUpRequest struct {
	ClientID   string      `json:"clientId,omitempty" example:"web-app"`
	Identifier *Identifier `json:"identifier" validate:"required"`
	Credential *Credential `json:"credential" validate:"required"`

	// Profile fields
	Firstname string `json:"firstname,omitempty" validate:"omitempty,max=100" example:"John"`
	Lastname  string `json:"lastname,omitempty" validate:"omitempty,max=100" example:"Doe"`
	Email     string `json:"email,omitempty" validate:"omitempty,email" example:"john@example.com"`
	Birthday  string `json:"birthday,omitempty" validate:"omitempty" example:"1990-01-01"`
	Locale    string `json:"locale,omitempty" validate:"omitempty,len=5" example:"en_US"`
}

// GetIdentifierValue returns the identifier value
func (r *SignUpRequest) GetIdentifierValue() string {
	if r.Identifier != nil {
		return r.Identifier.Value
	}
	return ""
}

// GetCredentialValue returns the credential value
func (r *SignUpRequest) GetCredentialValue() string {
	if r.Credential != nil {
		return r.Credential.Value
	}
	return ""
}

// ChangePasswordRequest DTO for changing password
// @Description Change password using credential pattern
type ChangePasswordRequest struct {
	ClientID      string      `json:"clientId,omitempty" example:"web-app"`
	OldCredential *Credential `json:"oldCredential" validate:"required"`
	NewCredential *Credential `json:"newCredential" validate:"required"`

	// User ID (populated by middleware, not from request body)
	UserID string `json:"-"`
}

// GetOldPassword extracts old password value
func (r *ChangePasswordRequest) GetOldPassword() string {
	if r.OldCredential != nil {
		return r.OldCredential.Value
	}
	return ""
}

// GetNewPassword extracts new password value
func (r *ChangePasswordRequest) GetNewPassword() string {
	if r.NewCredential != nil {
		return r.NewCredential.Value
	}
	return ""
}

// UpdateProfileRequest DTO for updating user profile
// @Description Update user profile information
type UpdateProfileRequest struct {
	Firstname string                 `json:"firstname,omitempty" validate:"omitempty,max=100" example:"John"`
	Lastname  string                 `json:"lastname,omitempty" validate:"omitempty,max=100" example:"Doe"`
	Email     string                 `json:"email,omitempty" validate:"omitempty,email" example:"john@example.com"`
	Birthday  string                 `json:"birthday,omitempty" validate:"omitempty" example:"1990-01-01"`
	Locale    string                 `json:"locale,omitempty" validate:"omitempty,len=5" example:"en_US"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`

	// User ID (populated by middleware, not from request body)
	UserID string `json:"-"`
}
