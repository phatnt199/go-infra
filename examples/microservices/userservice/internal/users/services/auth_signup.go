package services

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	"github.com/phatnt199/go-infra/pkg/crypto"
	uuid "github.com/satori/go.uuid"
)

// SignUpRequest contains the data needed to sign up a new user
type SignUpRequest struct {
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Locale    string `json:"locale"`
}

// SignUpResponse contains the response data after signup
type SignUpResponse struct {
	UserID       uuid.UUID `json:"userId"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
}

// SignUp creates a new user with credentials and profile
func (s *UserService) SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error) {
	// Check if username already exists
	existing, _ := s.repository.GetUserIdentifierBySchemeAndIdentifier(ctx, models.IdentifierSchemeUsername, req.Username)
	if existing != nil {
		return nil, errors.New("username already exists")
	}

	// Create user
	now := time.Now()
	user := &models.User{
		Status:      models.UserStatusActivated,
		UserType:    models.UserTypeSystem,
		ActivatedAt: &now,
	}

	user, err := s.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	// Create identifier
	identifier := &models.UserIdentifier{
		UserID:     user.ID,
		Scheme:     models.IdentifierSchemeUsername,
		Identifier: req.Username,
		Verified:   true,
		Details:    make(map[string]interface{}),
	}

	_, err = s.repository.CreateUserIdentifier(ctx, identifier)
	if err != nil {
		// Rollback user creation if identifier creation fails
		_ = s.repository.DeleteUser(ctx, user.ID)
		return nil, errors.Wrap(err, "failed to create user identifier")
	}

	// Hash password and create credential
	hashedPassword, err := s.hasher.HashPassword(req.Password, crypto.AlgorithmBcrypt)
	if err != nil {
		_ = s.repository.DeleteUser(ctx, user.ID)
		return nil, errors.Wrap(err, "failed to hash password")
	}

	credential := &models.UserCredential{
		UserID:     user.ID,
		Scheme:     models.CredentialSchemeBasic,
		Credential: hashedPassword,
		Details:    make(map[string]interface{}),
	}

	_, err = s.repository.CreateUserCredential(ctx, credential)
	if err != nil {
		_ = s.repository.DeleteUser(ctx, user.ID)
		return nil, errors.Wrap(err, "failed to create user credential")
	}

	// Create profile
	locale := req.Locale
	if locale == "" {
		locale = "en_US"
	}

	profile := &models.UserProfile{
		UserID:    user.ID,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Locale:    locale,
		Details:   make(map[string]interface{}),
	}

	_, err = s.repository.CreateUserProfile(ctx, profile)
	if err != nil {
		s.log.Warnf("failed to create user profile: %v", err)
		// Don't fail the signup if profile creation fails
	}

	// Generate JWT tokens
	claims := &crypto.Claims{
		UserID:   user.ID.String(),
		Username: req.Username,
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(claims)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate tokens")
	}

	return &SignUpResponse{
		UserID:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
