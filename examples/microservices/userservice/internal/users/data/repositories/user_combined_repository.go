package repositories

import (
	"context"

	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/contracts"
	uuid "github.com/satori/go.uuid"
)

// Combined operations that fetch related data

func (p *postgresUserRepository) GetUserFullDetails(ctx context.Context, userID uuid.UUID) (*contracts.UserFullDetails, error) {
	user, err := p.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get related data (ignore errors if not found)
	identifier, _ := p.GetUserIdentifierByUserID(ctx, userID)
	profile, _ := p.GetUserProfileByUserID(ctx, userID)

	return &contracts.UserFullDetails{
		User:       user,
		Identifier: identifier,
		Profile:    profile,
	}, nil
}
