package models

// CustomSignInRequest extends the standard sign-in request with additional fields
type CustomSignInRequest struct {
	ClientID   string      `json:"clientId,omitempty" example:"web-app"`
	Identifier *Identifier `json:"identifier" validate:"required"`
	Credential *Credential `json:"credential" validate:"required"`

	// Custom fields
	DeviceID  string `json:"deviceId,omitempty" example:"device-123"`
	IPAddress string `json:"-"` // Set from request context
	UserAgent string `json:"-"` // Set from request context
}

// Identifier represents user identifier with scheme and value
type Identifier struct {
	Scheme string `json:"scheme" validate:"required" example:"username"`
	Value  string `json:"value" validate:"required" example:"john"`
}

// Credential represents authentication credential with scheme and value
type Credential struct {
	Scheme string `json:"scheme" validate:"required" example:"basic"`
	Value  string `json:"value" validate:"required" example:"pass123"`
}

// GetIdentifierValue returns the identifier value
func (r *CustomSignInRequest) GetIdentifierValue() string {
	if r.Identifier != nil {
		return r.Identifier.Value
	}
	return ""
}

// GetCredentialValue returns the credential value
func (r *CustomSignInRequest) GetCredentialValue() string {
	if r.Credential != nil {
		return r.Credential.Value
	}
	return ""
}

// CustomSignUpRequest extends the standard sign-up request
type CustomSignUpRequest struct {
	ClientID   string      `json:"clientId,omitempty" example:"web-app"`
	Identifier *Identifier `json:"identifier" validate:"required"`
	Credential *Credential `json:"credential" validate:"required"`

	// Profile fields
	Firstname string `json:"firstname,omitempty" validate:"omitempty,max=100" example:"John"`
	Lastname  string `json:"lastname,omitempty" validate:"omitempty,max=100" example:"Doe"`
	Email     string `json:"email,omitempty" validate:"omitempty,email" example:"john@example.com"`
	Birthday  string `json:"birthday,omitempty" validate:"omitempty" example:"1990-01-01"`
	Locale    string `json:"locale,omitempty" validate:"omitempty,len=5" example:"en_US"`

	// Custom fields
	AcceptTerms  bool   `json:"acceptTerms" validate:"required,eq=true" example:"true"`
	ReferralCode string `json:"referralCode,omitempty" example:"REF123"`
	IPAddress    string `json:"-"`
	UserAgent    string `json:"-"`
}

// GetIdentifierValue returns the identifier value
func (r *CustomSignUpRequest) GetIdentifierValue() string {
	if r.Identifier != nil {
		return r.Identifier.Value
	}
	return ""
}

// GetCredentialValue returns the credential value
func (r *CustomSignUpRequest) GetCredentialValue() string {
	if r.Credential != nil {
		return r.Credential.Value
	}
	return ""
}
