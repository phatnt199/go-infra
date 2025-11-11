package defaultimpl

import (
	"github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/implementations/default/repository"
	"github.com/phatnt199/go-infra/pkg/component/authentication/implementations/default/service"
	"github.com/phatnt199/go-infra/pkg/crypto"
	"github.com/phatnt199/go-infra/pkg/logger"
	"gorm.io/gorm"
)

// NewDefaultImplementation creates a complete default authentication implementation
// This provides a ready-to-use auth service following microservice pattern
func NewDefaultImplementation(
	db *gorm.DB,
	jwtConfig *crypto.JWTConfig,
	log logger.Logger,
) (contracts.IAuthService, error) {
	// Create repository
	repo := repository.NewGormAuthRepository(db)

	// Create service
	authService, err := service.NewDefaultAuthService(log, repo, jwtConfig)
	if err != nil {
		return nil, err
	}

	return authService, nil
}
