package services

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/data/datamodels"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/dtos/v1/requests"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	"github.com/phatnt199/go-infra/pkg/crypto"
	"gorm.io/gorm"
)

// SignInResponse contains the response data after signin
type SignInResponse struct {
	UserID             string        `json:"userId"`
	AccessToken        string        `json:"accessToken"`
	RefreshToken       string        `json:"refreshToken"`
	AccessTokenExpiry  time.Duration `json:"-"` // Not serialized, used internally
	RefreshTokenExpiry time.Duration `json:"-"` // Not serialized, used internally
	Username           string        `json:"username,omitempty"`
	CreatedAt          time.Time     `json:"createdAt,omitempty"`
	UserType           string        `json:"userType,omitempty"`
	UserStatus         string        `json:"userStatus,omitempty"`
}

// SignIn authenticates a user and returns JWT tokens
func (s *UserService) SignIn(ctx context.Context, req *requests.SignInRequest) (*SignInResponse, error) {
	// Determine schemes to use (with defaults)
	identifierScheme := req.Identifier.Scheme
	if identifierScheme == "" {
		identifierScheme = string(models.IdentifierSchemeUsername)
	}

	credentialScheme := req.Credential.Scheme
	if credentialScheme == "" {
		credentialScheme = string(models.CredentialSchemeBasic)
	}

	// Get user identifier
	identifier, err := s.repository.GetUserIdentifierBySchemeAndIdentifier(ctx, identifierScheme, req.GetUsername())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, errors.Wrap(err, "failed to get user identifier")
	}

	// Get user
	user, err := s.repository.GetUserByID(ctx, identifier.UserID)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Check if user is activated
	if user.Status != string(datamodels.UserStatusActivated) {
		return nil, errors.New("user account is not activated")
	}

	// Get user credential
	credential, err := s.repository.GetUserCredentialByUserIDAndScheme(ctx, user.ID, credentialScheme)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Verify password
	valid, err := s.hasher.ComparePassword(req.GetPassword(), credential.Credential)
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
		Username: req.GetUsername(),
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(claims)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate tokens")
	}

	return &SignInResponse{
		UserID:             user.ID.String(),
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		AccessTokenExpiry:  s.jwtManager.GetAccessTokenExpiry(),
		RefreshTokenExpiry: s.jwtManager.GetRefreshTokenExpiry(),
		Username:           req.GetUsername(),
		CreatedAt:          user.CreatedAt,
		UserType:           user.UserType,
		UserStatus:         user.Status,
	}, nil
}
