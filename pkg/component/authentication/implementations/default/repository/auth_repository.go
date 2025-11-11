package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/phatnt199/go-infra/pkg/component/authentication/implementations/default/models"
)

// IAuthRepository defines repository operations for authentication
type IAuthRepository interface {
	// User operations
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// User Identifier operations
	CreateUserIdentifier(ctx context.Context, identifier *models.UserIdentifier) (*models.UserIdentifier, error)
	GetUserIdentifierBySchemeAndValue(ctx context.Context, scheme, value string) (*models.UserIdentifier, error)
	GetUserIdentifiersByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserIdentifier, error)
	UpdateUserIdentifier(ctx context.Context, identifier *models.UserIdentifier) (*models.UserIdentifier, error)
	DeleteUserIdentifier(ctx context.Context, id uuid.UUID) error

	// User Credential operations
	CreateUserCredential(ctx context.Context, credential *models.UserCredential) (*models.UserCredential, error)
	GetUserCredentialByUserIDAndScheme(ctx context.Context, userID uuid.UUID, scheme string) (*models.UserCredential, error)
	GetUserCredentialsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserCredential, error)
	UpdateUserCredential(ctx context.Context, credential *models.UserCredential) (*models.UserCredential, error)
	DeleteUserCredential(ctx context.Context, id uuid.UUID) error

	// User Profile operations
	CreateUserProfile(ctx context.Context, profile *models.UserProfile) (*models.UserProfile, error)
	GetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error)
	UpdateUserProfile(ctx context.Context, profile *models.UserProfile) (*models.UserProfile, error)
	DeleteUserProfile(ctx context.Context, id uuid.UUID) error

	// Complex queries
	GetUserFullDetails(ctx context.Context, userID uuid.UUID) (*models.UserFullDetails, error)
}
