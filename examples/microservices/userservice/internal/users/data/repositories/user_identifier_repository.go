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

// UserIdentifier CRUD operations

func (p *postgresUserRepository) CreateUserIdentifier(ctx context.Context, identifier *models.UserIdentifier) (*models.UserIdentifier, error) {
	dataModel := &datamodels.UserIdentifierDataModel{
		UserID:     identifier.UserID,
		Scheme:     identifier.Scheme,
		Identifier: identifier.Identifier,
		Verified:   identifier.Verified,
		Details:    datamodels.JSONB(identifier.Details),
	}

	if err := p.db.WithContext(ctx).Create(dataModel).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create user identifier")
	}

	identifier.ID = dataModel.ID
	identifier.BaseEntity.CreatedAt = dataModel.CreatedAt
	identifier.BaseEntity.ModifiedAt = dataModel.ModifiedAt
	return identifier, nil
}

func (p *postgresUserRepository) GetUserIdentifierByUserID(ctx context.Context, userID uuid.UUID) (*models.UserIdentifier, error) {
	var dataModel datamodels.UserIdentifierDataModel
	if err := p.db.WithContext(ctx).Where("user_id = ?", userID).First(&dataModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil without error if not found
		}
		return nil, errors.Wrap(err, "failed to get user identifier")
	}

	return p.toUserIdentifierModel(&dataModel), nil
}

func (p *postgresUserRepository) GetUserIdentifierBySchemeAndIdentifier(ctx context.Context, scheme, identifier string) (*models.UserIdentifier, error) {
	var dataModel datamodels.UserIdentifierDataModel
	if err := p.db.WithContext(ctx).
		Where("scheme = ? AND identifier = ?", scheme, identifier).
		First(&dataModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(err, "user identifier not found")
		}
		return nil, errors.Wrap(err, "failed to get user identifier")
	}

	return p.toUserIdentifierModel(&dataModel), nil
}

func (p *postgresUserRepository) UpdateUserIdentifier(ctx context.Context, identifier *models.UserIdentifier) (*models.UserIdentifier, error) {
	updates := map[string]interface{}{
		"scheme":     identifier.Scheme,
		"identifier": identifier.Identifier,
		"verified":   identifier.Verified,
		"details":    datamodels.JSONB(identifier.Details),
	}

	if err := p.db.WithContext(ctx).Model(&datamodels.UserIdentifierDataModel{}).
		Where("id = ?", identifier.ID).
		Updates(updates).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update user identifier")
	}

	return p.GetUserIdentifierByUserID(ctx, identifier.UserID)
}

// Helper method to convert data model to domain model
func (p *postgresUserRepository) toUserIdentifierModel(dm *datamodels.UserIdentifierDataModel) *models.UserIdentifier {
	return &models.UserIdentifier{
		BaseEntity: basemodels.BaseEntity{
			ID:         dm.ID,
			CreatedAt:  dm.CreatedAt,
			ModifiedAt: dm.ModifiedAt,
			DeletedAt:  dm.DeletedAt,
		},
		UserID:     dm.UserID,
		Scheme:     dm.Scheme,
		Identifier: dm.Identifier,
		Verified:   dm.Verified,
		Details:    map[string]interface{}(dm.Details),
	}
}
