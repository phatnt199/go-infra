package requests

import uuid "github.com/satori/go.uuid"

// ChangePasswordRequest DTO for changing password
// Supports both node-infra style (oldCredential/newCredential with schemes) and simple oldPassword/newPassword
type ChangePasswordRequest struct {
	// Node-infra style fields
	OldCredential *Credential `json:"oldCredential,omitempty"`
	NewCredential *Credential `json:"newCredential,omitempty"`

	// Simple style (for backwards compatibility)
	OldPassword string `json:"oldPassword,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`

	// User ID (populated by handler/service, not from request body)
	UserID uuid.UUID `json:"-"`
}

// GetOldPassword extracts old password from either credential or direct field
func (r *ChangePasswordRequest) GetOldPassword() string {
	if r.OldCredential != nil && r.OldCredential.Value != "" {
		return r.OldCredential.Value
	}
	return r.OldPassword
}

// GetNewPassword extracts new password from either credential or direct field
func (r *ChangePasswordRequest) GetNewPassword() string {
	if r.NewCredential != nil && r.NewCredential.Value != "" {
		return r.NewCredential.Value
	}
	return r.NewPassword
}
