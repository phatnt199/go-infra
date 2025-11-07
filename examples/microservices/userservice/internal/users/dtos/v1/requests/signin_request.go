package requests

// Credential represents authentication credential with scheme and value
type Credential struct {
	Scheme string `json:"scheme" validate:"required"`
	Value  string `json:"value" validate:"required"`
}

// Identifier represents user identifier with scheme and value
type Identifier struct {
	Scheme string `json:"scheme" validate:"required"`
	Value  string `json:"value" validate:"required"`
}

// SignInRequest DTO for user authentication
// Identifier/credential schemes.
// Schemes are required and must be explicitly specified.
//
// Example payload:
//
//	{
//		"identifier": {
//			"scheme": "username",
//			"value": "john"
//		},
//		"credential": {
//			"scheme": "basic",
//			"value": "pass123"
//		}
//	}
type SignInRequest struct {
	Identifier *Identifier `json:"identifier" validate:"required"`
	Credential *Credential `json:"credential" validate:"required"`
	ClientID   string      `json:"clientId,omitempty"` // optional client ID for multi-tenant scenarios
}

// GetUsername returns the username from the identifier
func (r *SignInRequest) GetUsername() string {
	if r.Identifier != nil && r.Identifier.Scheme == "username" {
		return r.Identifier.Value
	}
	return ""
}

// GetPassword returns the password from the credential
func (r *SignInRequest) GetPassword() string {
	if r.Credential != nil && r.Credential.Scheme == "basic" {
		return r.Credential.Value
	}
	return ""
}
