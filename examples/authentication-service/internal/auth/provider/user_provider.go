package provider

import (
	"context"
	"fmt"

	"github.com/phatnt199/go-infra/examples/authentication-service/internal/auth/models"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"github.com/phatnt199/go-infra/pkg/infra/postgres/gorm/contracts"
	uuid "github.com/satori/go.uuid"
)

// UserProvider implements IUserProvider interface
type UserProvider struct {
	dbContext contracts.GormDBContext
}

// NewUserProvider creates a new user provider
func NewUserProvider(dbContext contracts.GormDBContext) authContracts.IUserProvider {
	return &UserProvider{dbContext: dbContext}
}

// GetUserByIdentifier finds a user by identifier (username, email, phone, etc.)
func (p *UserProvider) GetUserByIdentifier(ctx context.Context, scheme, value string) (*authContracts.UserAuthInfo, error) {
	var user models.User
	var err error

	switch scheme {
	case "username":
		err = p.dbContext.DB().WithContext(ctx).Where("username = ?", value).First(&user).Error
	case "email":
		err = p.dbContext.DB().WithContext(ctx).Where("email = ?", value).First(&user).Error
	case "phone":
		// Assuming phone is stored in metadata or a separate column
		err = p.dbContext.DB().WithContext(ctx).Where("metadata->>'phone' = ?", value).First(&user).Error
	default:
		return nil, fmt.Errorf("unsupported identifier scheme: %s", scheme)
	}

	if err != nil {
		return nil, err
	}

	return p.toUserAuthInfo(&user), nil
}

// CreateUser creates a new user
func (p *UserProvider) CreateUser(ctx context.Context, req *authContracts.CreateUserRequest) (*authContracts.UserAuthInfo, error) {
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.PasswordHash,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Birthday:  req.Birthday,
		Locale:    req.Locale,
		UserType:  req.UserType,
		Status:    "active",
		Metadata:  req.Metadata,
	}

	if user.Locale == "" {
		user.Locale = "en_US"
	}

	if user.UserType == "" {
		user.UserType = "user"
	}

	err := p.dbContext.DB().WithContext(ctx).Create(user).Error
	if err != nil {
		return nil, err
	}

	return p.toUserAuthInfo(user), nil
}

// UpdatePassword updates user password
func (p *UserProvider) UpdatePassword(ctx context.Context, userID, newPasswordHash string) error {
	uid, err := uuid.FromString(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	return p.dbContext.DB().WithContext(ctx).Model(&models.User{}).
		Where("id = ?", uid).
		Update("password", newPasswordHash).Error
}

// GetUserByID gets user info by ID
func (p *UserProvider) GetUserByID(ctx context.Context, userID string) (*authContracts.UserAuthInfo, error) {
	uid, err := uuid.FromString(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	var user models.User
	err = p.dbContext.DB().WithContext(ctx).Where("id = ?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}

	return p.toUserAuthInfo(&user), nil
}

// UpdateProfile updates user profile information
func (p *UserProvider) UpdateProfile(ctx context.Context, req *authContracts.UpdateUserRequest) (*authContracts.UserAuthInfo, error) {
	uid, err := uuid.FromString(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Firstname != "" {
		updates["firstname"] = req.Firstname
	}
	if req.Lastname != "" {
		updates["lastname"] = req.Lastname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Birthday != nil {
		updates["birthday"] = req.Birthday
	}
	if req.Locale != "" {
		updates["locale"] = req.Locale
	}
	if req.Metadata != nil {
		updates["metadata"] = req.Metadata
	}

	err = p.dbContext.DB().WithContext(ctx).Model(&models.User{}).Where("id = ?", uid).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return p.GetUserByID(ctx, req.UserID)
}

// toUserAuthInfo converts User model to UserAuthInfo
func (p *UserProvider) toUserAuthInfo(user *models.User) *authContracts.UserAuthInfo {
	return &authContracts.UserAuthInfo{
		UserID:       user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.Password,
		Firstname:    user.Firstname,
		Lastname:     user.Lastname,
		Birthday:     user.Birthday,
		Locale:       user.Locale,
		Roles:        []string{"user"}, // Default role
		Status:       user.Status,
		UserType:     user.UserType,
		CreatedAt:    user.CreatedAt,
		Metadata:     user.Metadata,
	}
}
