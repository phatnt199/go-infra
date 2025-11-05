package repositories

import (
	"context"

	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/contracts"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	"github.com/phatnt199/go-infra/pkg/core/data"
	"github.com/phatnt199/go-infra/pkg/infra/postgres/gorm/repository"
	"github.com/phatnt199/go-infra/pkg/logger"
	"github.com/phatnt199/go-infra/pkg/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type postgresUserRepository struct {
	log                   logger.Logger
	gormGenericRepository data.GenericRepository[*models.User]
}

// CreateUser implements contracts.UserRepository.
func (p *postgresUserRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	err := p.gormGenericRepository.Add(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser implements contracts.UserRepository.
func (p *postgresUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return p.gormGenericRepository.Delete(ctx, id)
}

// GetAllUsers implements contracts.UserRepository.
func (p *postgresUserRepository) GetAllUsers(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*models.User], error) {
	result, err := p.gormGenericRepository.GetAll(ctx, listQuery)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetUserByID implements contracts.UserRepository.
func (p *postgresUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	panic("unimplemented")
}

// UpdateUser implements contracts.UserRepository.
func (p *postgresUserRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	panic("unimplemented")
}

func NewPostgresUserRepository(
	log logger.Logger,
	db *gorm.DB) contracts.UserRepository {
	gormRepository := repository.NewGenericGormRepository[*models.User](db)

	return &postgresUserRepository{
		log:                   log,
		gormGenericRepository: gormRepository,
	}
}
