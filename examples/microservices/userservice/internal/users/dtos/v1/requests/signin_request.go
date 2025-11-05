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
// Supports both node-infra style (identifier/credential with schemes) and simple username/password.
// Example payloads:
//
//	{"identifier": {"scheme": "username", "value": "john"}, "credential": {"scheme": "basic", "value": "pass123"}}
//	{"username": "john", "password": "pass123"}
type SignInRequest struct {
	// Node-infra style fields (preferred)
	Identifier *Identifier `json:"identifier,omitempty" validate:"omitempty,required_without=Username"`
	Credential *Credential `json:"credential,omitempty" validate:"omitempty,required_without=Password"`
	ClientID   string      `json:"clientId,omitempty"` // optional client ID for multi-tenant scenarios

	// Simple legacy style (for backwards compatibility)
	Username string `json:"username,omitempty" validate:"omitempty,required_without=Identifier,min=1"`
	Password string `json:"password,omitempty" validate:"omitempty,required_without=Credential,min=1"`
}

// GetUsername extracts username from either identifier or direct username field
func (r *SignInRequest) GetUsername() string {
	if r.Identifier != nil && r.Identifier.Scheme == "username" {
		return r.Identifier.Value
	}
	return r.Username
}

// GetPassword extracts password from either credential or direct password field
func (r *SignInRequest) GetPassword() string {
	if r.Credential != nil && r.Credential.Scheme == "basic" {
		return r.Credential.Value
	}
	return r.Password
}
