package services

import (
	"context"

	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/contracts"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	"github.com/phatnt199/go-infra/pkg/crypto"
	"github.com/phatnt199/go-infra/pkg/logger"
	uuid "github.com/satori/go.uuid"
)

type UserService struct {
	log        logger.Logger
	repository contracts.UserRepository
	hasher     *crypto.Hasher
	jwtManager *crypto.JWTManager
}

func NewUserService(
	log logger.Logger,
	repository contracts.UserRepository,
	jwtConfig *crypto.JWTConfig,
) (*UserService, error) {
	jwtManager, err := crypto.NewJWTManager(jwtConfig)
	if err != nil {
		return nil, err
	}

	return &UserService{
		log:        log,
		repository: repository,
		hasher:     crypto.NewHasher(nil), // Use default hasher config
		jwtManager: jwtManager,
	}, nil
}

// GetUserByID gets a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.repository.GetUserByID(ctx, id)
}

// GetUserFullDetails gets full user details without password
func (s *UserService) GetUserFullDetails(ctx context.Context, userID uuid.UUID) (*contracts.UserFullDetails, error) {
	return s.repository.GetUserFullDetails(ctx, userID)
}
