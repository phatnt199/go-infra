package requests

// SignUpRequest DTO for user registration
// Identifier/credential and additional profile fields
type SignUpRequest struct {
	// Required fields
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`

	// Profile fields - can be omitted, will use defaults if not provided
	Firstname string `json:"firstname,omitempty" validate:"omitempty,max=100"`
	Lastname  string `json:"lastname,omitempty" validate:"omitempty,max=100"`
	Birthday  string `json:"birthday,omitempty" validate:"omitempty"`     // ISO format: YYYY-MM-DD
	Locale    string `json:"locale,omitempty" validate:"omitempty,len=5"` // e.g., en_US, vi_VN

	// For credential object if provided instead of plain password
	Credential *Credential `json:"credential,omitempty"`

	// Additional dynamic fields (JSON passthrough)
	AdditionalFields map[string]interface{} `json:"-"`
}

// GetPassword extracts password from either credential or direct password field
func (r *SignUpRequest) GetPassword() string {
	if r.Credential != nil && r.Credential.Value != "" {
		return r.Credential.Value
	}
	return r.Password
}

// GetFirstname returns firstname, or a default value if empty
func (r *SignUpRequest) GetFirstname() string {
	if r.Firstname != "" {
		return r.Firstname
	}
	return "User"
}

// GetLastname returns lastname, or a default value if empty
func (r *SignUpRequest) GetLastname() string {
	if r.Lastname != "" {
		return r.Lastname
	}
	return ""
}

// GetLocale returns locale, or a default value if empty
func (r *SignUpRequest) GetLocale() string {
	if r.Locale != "" {
		return r.Locale
	}
	return "en_US"
}
