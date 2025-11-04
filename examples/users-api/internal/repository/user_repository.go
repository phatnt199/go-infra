package repository

import (
	"context"

	"github.com/phatnt199/go-infra/examples/users-api/internal/domain"
	"github.com/phatnt199/go-infra/pkg/logger"

	"gorm.io/gorm"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB, logger logger.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// GetAll retrieves all users from the database
func (r *UserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	var users []domain.User

	r.logger.Info("Fetching all users from database")

	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		r.logger.Errorw("Failed to fetch users from database", logger.Fields{
			"error": err,
		})
		return nil, err
	}

	r.logger.Infow("Successfully fetched users", logger.Fields{
		"count": len(users),
	})

	return users, nil
}
