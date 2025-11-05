package services

import (
	"context"

	"emperror.dev/errors"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	uuid "github.com/satori/go.uuid"
)

// DisableUser disables a user account
func (s *UserService) DisableUser(ctx context.Context, userID uuid.UUID) error {
	// Get user
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "failed to get user")
	}

	// Update status
	user.Status = models.UserStatusDeactivated

	// Save user
	_, err = s.repository.UpdateUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, "failed to disable user")
	}

	return nil
}

// EnableUser enables a user account
func (s *UserService) EnableUser(ctx context.Context, userID uuid.UUID) error {
	// Get user
	user, err := s.repository.GetUserByID(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "failed to get user")
	}

	// Update status
	user.Status = models.UserStatusActivated

	// Save user
	_, err = s.repository.UpdateUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, "failed to enable user")
	}

	return nil
}
