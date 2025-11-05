package contracts

import (
	"context"

	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	"github.com/phatnt199/go-infra/pkg/utils"
	uuid "github.com/satori/go.uuid"
)

type UserRepository interface {
	// User operations
	GetAllUsers(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*models.User], error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// UserIdentifier operations
	CreateUserIdentifier(ctx context.Context, identifier *models.UserIdentifier) (*models.UserIdentifier, error)
	GetUserIdentifierByUserID(ctx context.Context, userID uuid.UUID) (*models.UserIdentifier, error)
	GetUserIdentifierBySchemeAndIdentifier(ctx context.Context, scheme, identifier string) (*models.UserIdentifier, error)
	UpdateUserIdentifier(ctx context.Context, identifier *models.UserIdentifier) (*models.UserIdentifier, error)

	// UserCredential operations
	CreateUserCredential(ctx context.Context, credential *models.UserCredential) (*models.UserCredential, error)
	GetUserCredentialByUserIDAndScheme(ctx context.Context, userID uuid.UUID, scheme string) (*models.UserCredential, error)
	UpdateUserCredential(ctx context.Context, credential *models.UserCredential) (*models.UserCredential, error)

	// UserProfile operations
	CreateUserProfile(ctx context.Context, profile *models.UserProfile) (*models.UserProfile, error)
	GetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error)
	UpdateUserProfile(ctx context.Context, profile *models.UserProfile) (*models.UserProfile, error)

	// Combined operations
	GetUserFullDetails(ctx context.Context, userID uuid.UUID) (*UserFullDetails, error)
}

// UserFullDetails contains all user-related information (without password)
type UserFullDetails struct {
	User       *models.User           `json:"user"`
	Identifier *models.UserIdentifier `json:"identifier,omitempty"`
	Profile    *models.UserProfile    `json:"profile,omitempty"`
}
