package services

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	"github.com/phatnt199/go-infra/pkg/crypto"
)

// SignInRequest contains the data needed to sign in
type SignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// SignInResponse contains the response data after signin
type SignInResponse struct {
	UserID       string `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	// Additional fields for handlers
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
}

// SignIn authenticates a user and returns JWT tokens
func (s *UserService) SignIn(ctx context.Context, req *SignInRequest) (*SignInResponse, error) {
	// Get user identifier
	identifier, err := s.repository.GetUserIdentifierBySchemeAndIdentifier(ctx, models.IdentifierSchemeUsername, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Get user
	user, err := s.repository.GetUserByID(ctx, identifier.UserID)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Check if user is activated
	if user.Status != models.UserStatusActivated {
		return nil, errors.New("user account is not activated")
	}

	// Get user credential
	credential, err := s.repository.GetUserCredentialByUserIDAndScheme(ctx, user.ID, models.CredentialSchemeBasic)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Verify password
	valid, err := s.hasher.ComparePassword(req.Password, credential.Credential)
	if err != nil || !valid {
		return nil, errors.New("invalid username or password")
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	_, _ = s.repository.UpdateUser(ctx, user)

	// Generate JWT tokens
	claims := &crypto.Claims{
		UserID:   user.ID.String(),
		Username: req.Username,
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(claims)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate tokens")
	}

	return &SignInResponse{
		UserID:       user.ID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		// include createdAt and username so handlers can build full AuthResponse
		Username:  req.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}
