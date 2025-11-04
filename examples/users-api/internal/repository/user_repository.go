package repository

import (
	"context"

	"github.com/google/uuid"
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

// GetByID retrieves a user by ID from the database
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User

	r.logger.Infow("Fetching user by ID", logger.Fields{
		"id": id,
	})

	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Infow("User not found", logger.Fields{
				"id": id,
			})
			return nil, nil
		}
		r.logger.Errorw("Failed to fetch user from database", logger.Fields{
			"id":    id,
			"error": err,
		})
		return nil, err
	}

	r.logger.Infow("Successfully fetched user", logger.Fields{
		"id": id,
	})

	return &user, nil
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	r.logger.Infow("Creating new user", logger.Fields{
		"email": user.Email,
		"name":  user.Name,
	})

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		r.logger.Errorw("Failed to create user in database", logger.Fields{
			"email": user.Email,
			"error": err,
		})
		return err
	}

	r.logger.Infow("Successfully created user", logger.Fields{
		"id":    user.ID,
		"email": user.Email,
	})

	return nil
}

// Update updates an existing user in the database
func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, updates *domain.User) (*domain.User, error) {
	r.logger.Infow("Updating user", logger.Fields{
		"id": id,
	})

	// First fetch the existing user
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Infow("User not found for update", logger.Fields{
				"id": id,
			})
			return nil, nil
		}
		r.logger.Errorw("Failed to fetch user for update", logger.Fields{
			"id":    id,
			"error": err,
		})
		return nil, err
	}

	// Update the fields
	if updates.Name != "" {
		user.Name = updates.Name
	}
	if updates.Email != "" {
		user.Email = updates.Email
	}

	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		r.logger.Errorw("Failed to update user in database", logger.Fields{
			"id":    id,
			"error": err,
		})
		return nil, err
	}

	r.logger.Infow("Successfully updated user", logger.Fields{
		"id": id,
	})

	return &user, nil
}

// Delete deletes a user from the database (soft delete)
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.logger.Infow("Deleting user", logger.Fields{
		"id": id,
	})

	if err := r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id).Error; err != nil {
		r.logger.Errorw("Failed to delete user from database", logger.Fields{
			"id":    id,
			"error": err,
		})
		return err
	}

	r.logger.Infow("Successfully deleted user", logger.Fields{
		"id": id,
	})

	return nil
}
