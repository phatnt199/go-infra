package repositories

import (
	"context"

	"emperror.dev/errors"
	basemodels "github.com/phatnt199/go-infra/examples/microservices/userservice/internal/shared/data/models"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/data/datamodels"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// UserProfile CRUD operations

func (p *postgresUserRepository) CreateUserProfile(ctx context.Context, profile *models.UserProfile) (*models.UserProfile, error) {
	dataModel := &datamodels.UserProfileDataModel{
		UserID:    profile.UserID,
		Firstname: profile.Firstname,
		Lastname:  profile.Lastname,
		Birthday:  profile.Birthday,
		Locale:    profile.Locale,
		Details:   datamodels.JSONB(profile.Details),
	}

	if err := p.db.WithContext(ctx).Create(dataModel).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create user profile")
	}

	profile.ID = dataModel.ID
	profile.BaseEntity.CreatedAt = dataModel.CreatedAt
	profile.BaseEntity.ModifiedAt = dataModel.ModifiedAt
	return profile, nil
}

func (p *postgresUserRepository) GetUserProfileByUserID(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	var dataModel datamodels.UserProfileDataModel
	if err := p.db.WithContext(ctx).Where("user_id = ?", userID).First(&dataModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil without error if not found
		}
		return nil, errors.Wrap(err, "failed to get user profile")
	}

	return p.toUserProfileModel(&dataModel), nil
}

func (p *postgresUserRepository) UpdateUserProfile(ctx context.Context, profile *models.UserProfile) (*models.UserProfile, error) {
	updates := map[string]interface{}{
		"firstname": profile.Firstname,
		"lastname":  profile.Lastname,
		"birthday":  profile.Birthday,
		"locale":    profile.Locale,
		"details":   datamodels.JSONB(profile.Details),
	}

	if err := p.db.WithContext(ctx).Model(&datamodels.UserProfileDataModel{}).
		Where("id = ?", profile.ID).
		Updates(updates).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update user profile")
	}

	return p.GetUserProfileByUserID(ctx, profile.UserID)
}

// Helper method to convert data model to domain model
func (p *postgresUserRepository) toUserProfileModel(dm *datamodels.UserProfileDataModel) *models.UserProfile {
	return &models.UserProfile{
		BaseEntity: basemodels.BaseEntity{
			ID:         dm.ID,
			CreatedAt:  dm.CreatedAt,
			ModifiedAt: dm.ModifiedAt,
			DeletedAt:  dm.DeletedAt,
		},
		UserID:    dm.UserID,
		Firstname: dm.Firstname,
		Lastname:  dm.Lastname,
		Birthday:  dm.Birthday,
		Locale:    dm.Locale,
		Details:   map[string]interface{}(dm.Details),
	}
}
