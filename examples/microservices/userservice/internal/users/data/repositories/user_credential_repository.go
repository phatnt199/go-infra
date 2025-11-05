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

// UserCredential CRUD operations

func (p *postgresUserRepository) CreateUserCredential(ctx context.Context, credential *models.UserCredential) (*models.UserCredential, error) {
	dataModel := &datamodels.UserCredentialDataModel{
		UserID:     credential.UserID,
		Scheme:     credential.Scheme,
		Credential: credential.Credential,
		Details:    datamodels.JSONB(credential.Details),
	}

	if err := p.db.WithContext(ctx).Create(dataModel).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create user credential")
	}

	credential.ID = dataModel.ID
	credential.BaseEntity.CreatedAt = dataModel.CreatedAt
	credential.BaseEntity.ModifiedAt = dataModel.ModifiedAt
	return credential, nil
}

func (p *postgresUserRepository) GetUserCredentialByUserIDAndScheme(ctx context.Context, userID uuid.UUID, scheme string) (*models.UserCredential, error) {
	var dataModel datamodels.UserCredentialDataModel
	if err := p.db.WithContext(ctx).
		Where("user_id = ? AND scheme = ?", userID, scheme).
		First(&dataModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(err, "user credential not found")
		}
		return nil, errors.Wrap(err, "failed to get user credential")
	}

	return p.toUserCredentialModel(&dataModel), nil
}

func (p *postgresUserRepository) UpdateUserCredential(ctx context.Context, credential *models.UserCredential) (*models.UserCredential, error) {
	updates := map[string]interface{}{
		"credential": credential.Credential,
		"details":    datamodels.JSONB(credential.Details),
	}

	if err := p.db.WithContext(ctx).Model(&datamodels.UserCredentialDataModel{}).
		Where("id = ?", credential.ID).
		Updates(updates).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update user credential")
	}

	return p.GetUserCredentialByUserIDAndScheme(ctx, credential.UserID, credential.Scheme)
}

// Helper method to convert data model to domain model
func (p *postgresUserRepository) toUserCredentialModel(dm *datamodels.UserCredentialDataModel) *models.UserCredential {
	return &models.UserCredential{
		BaseEntity: basemodels.BaseEntity{
			ID:         dm.ID,
			CreatedAt:  dm.CreatedAt,
			ModifiedAt: dm.ModifiedAt,
			DeletedAt:  dm.DeletedAt,
		},
		UserID:     dm.UserID,
		Scheme:     dm.Scheme,
		Credential: dm.Credential,
		Details:    map[string]interface{}(dm.Details),
	}
}
