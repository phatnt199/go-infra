package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/phatnt199/go-infra/pkg/component/authentication/implementations/default/models"
	"gorm.io/gorm"
)

// GormAuthRepository implements IAuthRepository using GORM
type GormAuthRepository struct {
	db *gorm.DB
}

// NewGormAuthRepository creates a new GORM-based auth repository
func NewGormAuthRepository(db *gorm.DB) IAuthRepository {
	return &GormAuthRepository{db: db}
}

// User operations

func (r *GormAuthRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormAuthRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormAuthRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormAuthRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}

// User Identifier operations

func (r *GormAuthRepository) CreateUserIdentifier(ctx context.Context, identifier *models.UserIdentifier) (*models.UserIdentifier, error) {
	if err := r.db.WithContext(ctx).Create(identifier).Error; err != nil {
		return nil, err
	}
	return identifier, nil
}

func (r *GormAuthRepository) GetUserIdentifierBySchemeAndValue(ctx context.Context, scheme, value string) (*models.UserIdentifier, error) {
	var identifier models.UserIdentifier
	if err := r.db.WithContext(ctx).
		Where("scheme = ? AND identifier = ?", scheme, value).
		First(&identifier).Error; err != nil {
		return nil, err
	}
	return &identifier, nil
}

func (r *GormAuthRepository) GetUserIdentifiersByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserIdentifier, error) {
	var identifiers []*models.UserIdentifier
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&identifiers).Error; err != nil {
		return nil, err
	}
	return identifiers, nil
}

func (r *GormAuthRepository) UpdateUserIdentifier(ctx context.Context, identifier *models.UserIdentifier) (*models.UserIdentifier, error) {
	if err := r.db.WithContext(ctx).Save(identifier).Error; err != nil {
		return nil, err
	}
	return identifier, nil
}

func (r *GormAuthRepository) DeleteUserIdentifier(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.UserIdentifier{}, "id = ?", id).Error
}

// User Credential operations

func (r *GormAuthRepository) CreateUserCredential(ctx context.Context, credential *models.UserCredential) (*models.UserCredential, error) {
	if err := r.db.WithContext(ctx).Create(credential).Error; err != nil {
		return nil, err
	}
	return credential, nil
}

func (r *GormAuthRepository) GetUserCredentialByUserIDAndScheme(ctx context.Context, userID uuid.UUID, scheme string) (*models.UserCredential, error) {
	var credential models.UserCredential
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND scheme = ?", userID, scheme).
		First(&credential).Error; err != nil {
		return nil, err
	}
	return &credential, nil
}

func (r *GormAuthRepository) GetUserCredentialsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.UserCredential, error) {
	var credentials []*models.UserCredential
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&credentials).Error; err != nil {
		return nil, err
	}
	return credentials, nil
}

func (r *GormAuthRepository) UpdateUserCredential(ctx context.Context, credential *models.UserCredential) (*models.UserCredential, error) {
	if err := r.db.WithContext(ctx).Save(credential).Error; err != nil {
		return nil, err
	}
	return credential, nil
}

func (r *GormAuthRepository) DeleteUserCredential(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.UserCredential{}, "id = ?", id).Error
}

// User Profile operations

func (r *GormAuthRepository) CreateUserProfile(ctx context.Context, profile *models.UserProfile) (*models.UserProfile, error) {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *GormAuthRepository) GetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	var profile models.UserProfile
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *GormAuthRepository) UpdateUserProfile(ctx context.Context, profile *models.UserProfile) (*models.UserProfile, error) {
	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *GormAuthRepository) DeleteUserProfile(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.UserProfile{}, "id = ?", id).Error
}

// Complex queries

func (r *GormAuthRepository) GetUserFullDetails(ctx context.Context, userID uuid.UUID) (*models.UserFullDetails, error) {
	details := &models.UserFullDetails{}

	// Get user
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	details.User = user

	// Get primary identifier (first one found, or could filter by scheme)
	identifiers, err := r.GetUserIdentifiersByUserID(ctx, userID)
	if err == nil && len(identifiers) > 0 {
		details.Identifier = identifiers[0]
	}

	// Get profile
	profile, err := r.GetUserProfileByUserID(ctx, userID)
	if err == nil {
		details.Profile = profile
	}

	// Get basic credential
	credential, err := r.GetUserCredentialByUserIDAndScheme(ctx, userID, string(models.CredentialSchemeBasic))
	if err == nil {
		details.Credential = credential
	}

	return details, nil
}
