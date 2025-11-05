package services

import (
	"context"
	"time"

	"emperror.dev/errors"
	uuid "github.com/satori/go.uuid"
)

// UpdateProfileRequest contains the data needed to update a user profile
type UpdateProfileRequest struct {
	UserID    uuid.UUID  `json:"userId" validate:"required"`
	Firstname string     `json:"firstname"`
	Lastname  string     `json:"lastname"`
	Birthday  *time.Time `json:"birthday"`
	Locale    string     `json:"locale"`
}

// UpdateProfile updates a user's profile
func (s *UserService) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) error {
	// Get existing profile
	profile, err := s.repository.GetUserProfileByUserID(ctx, req.UserID)
	if err != nil || profile == nil {
		return errors.Wrap(err, "failed to get user profile")
	}

	// Update fields
	if req.Firstname != "" {
		profile.Firstname = req.Firstname
	}
	if req.Lastname != "" {
		profile.Lastname = req.Lastname
	}
	if req.Birthday != nil {
		profile.Birthday = req.Birthday
	}
	if req.Locale != "" {
		profile.Locale = req.Locale
	}

	// Save profile
	_, err = s.repository.UpdateUserProfile(ctx, profile)
	if err != nil {
		return errors.Wrap(err, "failed to update profile")
	}

	return nil
}
