package contracts

import (
	"context"

	dtosv1 "github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/dtos/v1"
	uuid "github.com/satori/go.uuid"
)

// UserService handles user business operations
type UserService interface {
	GetUserFullDetails(ctx context.Context, userID uuid.UUID) (*dtosv1.UserFullDetailsResponse, error)
	UpdateUserProfile(ctx context.Context, userID uuid.UUID, req *dtosv1.UpdateProfileRequest) (*dtosv1.UserProfileResponse, error)
	DisableUser(ctx context.Context, userID uuid.UUID) error
	ActivateUser(ctx context.Context, userID uuid.UUID) error
	GetAllUsers(ctx context.Context, page, size int) (*dtosv1.UserListResponse, error)
}
