package repositories

import (
	"context"

	"emperror.dev/errors"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/contracts"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/data/datamodels"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	"github.com/phatnt199/go-infra/pkg/logger"
	"github.com/phatnt199/go-infra/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type postgresUserRepository struct {
	log logger.Logger
	db  *gorm.DB
}

// NewPostgresUserRepository creates a new postgres user repository
func NewPostgresUserRepository(
	log logger.Logger,
	db *gorm.DB) contracts.UserRepository {
	return &postgresUserRepository{
		log: log,
		db:  db,
	}
}

// User CRUD operations

func (p *postgresUserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	dataModel := &datamodels.UserDataModel{
		ID:          user.ID,
		ActivatedAt: user.ActivatedAt,
		LastLoginAt: user.LastLoginAt,
		ParentID:    user.ParentID,
		ValidFrom:   user.ValidFrom,
		ValidTo:     user.ValidTo,
	}

	// Convert domain enums to data model enums
	// dataModel.Status.FromInt(int(user.Status))
	dataModel.UserType = datamodels.UserTypeEnum(user.UserType)

	if err := p.db.WithContext(ctx).Create(dataModel).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	user.ID = dataModel.ID
	user.CreatedAt = dataModel.CreatedAt
	user.BaseEntity.ModifiedAt = dataModel.ModifiedAt
	return user, nil
}

func (p *postgresUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var dataModel datamodels.UserDataModel
	if err := p.db.WithContext(ctx).Where("id = ?", id).First(&dataModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrap(err, "user not found")
		}
		return nil, errors.Wrap(err, "failed to get user")
	}

	return p.toUserModel(&dataModel), nil
}

func (p *postgresUserRepository) GetAllUsers(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*models.User], error) {
	var dataModels []datamodels.UserDataModel
	var total int64

	query := p.db.WithContext(ctx).Model(&datamodels.UserDataModel{})

	if err := query.Count(&total).Error; err != nil {
		return nil, errors.Wrap(err, "failed to count users")
	}

	if listQuery != nil {
		if listQuery.GetLimit() > 0 {
			query = query.Limit(listQuery.GetLimit())
		}
		if listQuery.GetOffset() > 0 {
			query = query.Offset(listQuery.GetOffset())
		}
	}

	if err := query.Find(&dataModels).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get users")
	}

	users := make([]*models.User, len(dataModels))
	for i, dm := range dataModels {
		users[i] = p.toUserModel(&dm)
	}

	return utils.NewListResult(users, listQuery.GetSize(), listQuery.GetPage(), total), nil
}

func (p *postgresUserRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	updates := map[string]interface{}{
		"status":        user.Status,
		"user_type":     user.UserType,
		"activated_at":  user.ActivatedAt,
		"last_login_at": user.LastLoginAt,
		"parent_id":     user.ParentID,
		"valid_from":    user.ValidFrom,
		"valid_to":      user.ValidTo,
	}

	if err := p.db.WithContext(ctx).Model(&datamodels.UserDataModel{}).
		Where("id = ?", user.ID).
		Updates(updates).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update user")
	}

	return p.GetUserByID(ctx, user.ID)
}

func (p *postgresUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := p.db.WithContext(ctx).Delete(&datamodels.UserDataModel{}, "id = ?", id).Error; err != nil {
		return errors.Wrap(err, "failed to delete user")
	}
	return nil
}

// Helper method to convert data model to domain model
func (p *postgresUserRepository) toUserModel(dm *datamodels.UserDataModel) *models.User {
	user := &models.User{
		Status:      dm.Status.String(),
		UserType:    dm.UserType.String(),
		ActivatedAt: dm.ActivatedAt,
		LastLoginAt: dm.LastLoginAt,
		ParentID:    dm.ParentID,
		ValidFrom:   dm.ValidFrom,
		ValidTo:     dm.ValidTo,
	}

	// Set base entity fields
	user.ID = dm.ID
	user.CreatedAt = dm.CreatedAt
	user.ModifiedAt = dm.ModifiedAt

	return user
}
