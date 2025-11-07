package services

import (
	"context"

	"emperror.dev/errors"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	"github.com/phatnt199/go-infra/pkg/crypto"
	uuid "github.com/satori/go.uuid"
)

// ChangePasswordRequest contains the data needed to change password
type ChangePasswordRequest struct {
	UserID      uuid.UUID `json:"userId" validate:"required"`
	OldPassword string    `json:"oldPassword" validate:"required"`
	NewPassword string    `json:"newPassword" validate:"required,min=8"`
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, req *ChangePasswordRequest) error {
	// Get user credential
	credential, err := s.repository.GetUserCredentialByUserIDAndScheme(ctx, req.UserID, models.CredentialSchemeBasic)
	if err != nil {
		return errors.Wrap(err, "failed to get user credential")
	}

	// Verify old password
	valid, err := s.hasher.ComparePassword(req.OldPassword, credential.Credential)
	if err != nil || !valid {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := s.hasher.HashPassword(req.NewPassword, crypto.AlgorithmBcrypt)
	if err != nil {
		return errors.Wrap(err, "failed to hash new password")
	}

	// Update credential
	credential.Credential = hashedPassword
	_, err = s.repository.UpdateUserCredential(ctx, credential)
	if err != nil {
		return errors.Wrap(err, "failed to update credential")
	}

	return nil
}
