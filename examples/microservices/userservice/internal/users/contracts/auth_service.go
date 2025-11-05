package contracts

import (
	"context"

	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/dtos/v1"
)

// AuthService handles authentication operations
type AuthService interface {
	SignUp(ctx context.Context, req *dtosv1.SignUpRequest) (*dtosv1.AuthResponse, error)
	SignIn(ctx context.Context, req *dtosv1.SignInRequest) (*dtosv1.AuthResponse, error)
	ChangePassword(ctx context.Context, req *dtosv1.ChangePasswordRequest) error
}
