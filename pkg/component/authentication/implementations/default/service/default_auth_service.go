package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/phatnt199/go-infra/pkg/component/authentication/implementations/default/models"
	"github.com/phatnt199/go-infra/pkg/component/authentication/implementations/default/repository"
	authModels "github.com/phatnt199/go-infra/pkg/component/authentication/models"
	"github.com/phatnt199/go-infra/pkg/crypto"
	"github.com/phatnt199/go-infra/pkg/logger"
	"gorm.io/gorm"
)

// DefaultAuthService provides a default implementation of authentication service
// following the microservice pattern with User/Identifier/Credential/Profile separation
type DefaultAuthService struct {
	log        logger.Logger
	repository repository.IAuthRepository
	hasher     *crypto.Hasher
	jwtManager *crypto.JWTManager
}

// NewDefaultAuthService creates a new default auth service
func NewDefaultAuthService(
	log logger.Logger,
	repository repository.IAuthRepository,
	jwtConfig *crypto.JWTConfig,
) (*DefaultAuthService, error) {
	jwtManager, err := crypto.NewJWTManager(jwtConfig)
	if err != nil {
		return nil, err
	}

	return &DefaultAuthService{
		log:        log,
		repository: repository,
		hasher:     crypto.NewHasher(nil), // Use default hasher config
		jwtManager: jwtManager,
	}, nil
}

// SignIn authenticates a user and returns JWT tokens
func (s *DefaultAuthService) SignIn(ctx context.Context, req *authModels.SignInRequest) (*authModels.AuthResponse, error) {
	s.log.Infof("DefaultAuth SignIn: identifier=%s", req.GetIdentifierValue())

	// Get identifier scheme (default to username)
	identifierScheme := req.Identifier.Scheme
	if identifierScheme == "" {
		identifierScheme = string(models.IdentifierSchemeUsername)
	}

	// Get credential scheme (default to basic)
	credentialScheme := req.Credential.Scheme
	if credentialScheme == "" {
		credentialScheme = string(models.CredentialSchemeBasic)
	}

	// Get user identifier
	identifier, err := s.repository.GetUserIdentifierBySchemeAndValue(ctx, identifierScheme, req.GetIdentifierValue())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("invalid credentials")
		}
		return nil, fmt.Errorf("failed to get user identifier: %w", err)
	}

	// Get user
	user, err := s.repository.GetUserByID(ctx, identifier.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check user status
	if user.Status != string(models.UserStatusActive) {
		return nil, fmt.Errorf("user account is not active")
	}

	// Get user credential
	credential, err := s.repository.GetUserCredentialByUserIDAndScheme(ctx, user.ID, credentialScheme)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	valid, err := s.hasher.ComparePassword(req.GetCredentialValue(), credential.Credential)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	_, _ = s.repository.UpdateUser(ctx, user)

	// Get profile for additional info
	profile, _ := s.repository.GetUserProfileByUserID(ctx, user.ID)

	// Generate JWT tokens
	claims := &crypto.Claims{
		UserID:   user.ID.String(),
		Username: identifier.Identifier,
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Build response
	response := &authModels.AuthResponse{
		UserID:    user.ID.String(),
		Username:  identifier.Identifier,
		Status:    user.Status,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		AccessToken: &authModels.TokenInfo{
			Value:     accessToken,
			Scheme:    "bearer",
			Type:      "access",
			ExpiresAt: time.Now().Add(s.jwtManager.GetAccessTokenExpiry()),
		},
		RefreshToken: &authModels.TokenInfo{
			Value:     refreshToken,
			Scheme:    "bearer",
			Type:      "refresh",
			ExpiresAt: time.Now().Add(s.jwtManager.GetRefreshTokenExpiry()),
		},
	}

	if profile != nil {
		response.Email = profile.Email
		response.Firstname = profile.Firstname
		response.Lastname = profile.Lastname
	}

	return response, nil
}

// SignUp creates a new user
func (s *DefaultAuthService) SignUp(ctx context.Context, req *authModels.SignUpRequest) (*authModels.AuthResponse, error) {
	s.log.Infof("DefaultAuth SignUp: identifier=%s", req.GetIdentifierValue())

	// Get schemes
	identifierScheme := req.Identifier.Scheme
	if identifierScheme == "" {
		identifierScheme = string(models.IdentifierSchemeUsername)
	}

	credentialScheme := req.Credential.Scheme
	if credentialScheme == "" {
		credentialScheme = string(models.CredentialSchemeBasic)
	}

	// Check if user already exists
	existingIdentifier, err := s.repository.GetUserIdentifierBySchemeAndValue(ctx, identifierScheme, req.GetIdentifierValue())
	if err == nil && existingIdentifier != nil {
		return nil, fmt.Errorf("user with this %s already exists", identifierScheme)
	}

	// Hash password using Bcrypt (can be changed to Argon2 for better security)
	passwordHash, err := s.hasher.HashPassword(req.GetCredentialValue(), crypto.AlgorithmBcrypt)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	now := time.Now()
	user := &models.User{
		UserType:    "user",
		Status:      string(models.UserStatusActive),
		ActivatedAt: &now,
	}
	user, err = s.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create identifier
	identifier := &models.UserIdentifier{
		UserID:     user.ID,
		Scheme:     identifierScheme,
		Identifier: req.GetIdentifierValue(),
		Verified:   false,
	}
	_, err = s.repository.CreateUserIdentifier(ctx, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to create identifier: %w", err)
	}

	// Create credential
	credential := &models.UserCredential{
		UserID:     user.ID,
		Scheme:     credentialScheme,
		Credential: passwordHash,
	}
	_, err = s.repository.CreateUserCredential(ctx, credential)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	// Create profile
	var birthday *time.Time
	if req.Birthday != "" {
		t, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			birthday = &t
		}
	}

	profile := &models.UserProfile{
		UserID:    user.ID,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Birthday:  birthday,
		Locale:    req.Locale,
	}
	if profile.Locale == "" {
		profile.Locale = "en_US"
	}
	_, err = s.repository.CreateUserProfile(ctx, profile)
	if err != nil {
		s.log.Warnf("Failed to create profile: %v", err)
		// Don't fail signup if profile creation fails
	}

	// Generate tokens
	claims := &crypto.Claims{
		UserID:   user.ID.String(),
		Username: identifier.Identifier,
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(claims)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &authModels.AuthResponse{
		UserID:    user.ID.String(),
		Username:  identifier.Identifier,
		Email:     profile.Email,
		Firstname: profile.Firstname,
		Lastname:  profile.Lastname,
		Status:    user.Status,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		AccessToken: &authModels.TokenInfo{
			Value:     accessToken,
			Scheme:    "bearer",
			Type:      "access",
			ExpiresAt: time.Now().Add(s.jwtManager.GetAccessTokenExpiry()),
		},
		RefreshToken: &authModels.TokenInfo{
			Value:     refreshToken,
			Scheme:    "bearer",
			Type:      "refresh",
			ExpiresAt: time.Now().Add(s.jwtManager.GetRefreshTokenExpiry()),
		},
	}, nil
}

// ChangePassword changes user password
func (s *DefaultAuthService) ChangePassword(ctx context.Context, req *authModels.ChangePasswordRequest) error {
	s.log.Infof("DefaultAuth ChangePassword: userID=%s", req.UserID)

	// Parse user ID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	// Get credential scheme
	credentialScheme := req.OldCredential.Scheme
	if credentialScheme == "" {
		credentialScheme = string(models.CredentialSchemeBasic)
	}

	// Get user credential
	credential, err := s.repository.GetUserCredentialByUserIDAndScheme(ctx, userID, credentialScheme)
	if err != nil {
		return fmt.Errorf("failed to get credential")
	}

	// Verify old password
	valid, err := s.hasher.ComparePassword(req.GetOldPassword(), credential.Credential)
	if err != nil || !valid {
		return fmt.Errorf("invalid old password")
	}

	// Hash new password using Bcrypt (can be changed to Argon2 for better security)
	newPasswordHash, err := s.hasher.HashPassword(req.GetNewPassword(), crypto.AlgorithmBcrypt)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update credential
	credential.Credential = newPasswordHash
	_, err = s.repository.UpdateUserCredential(ctx, credential)
	if err != nil {
		return fmt.Errorf("failed to update credential: %w", err)
	}

	return nil
}

// GetProfile gets user profile
func (s *DefaultAuthService) GetProfile(ctx context.Context, userID string) (*authModels.UserResponse, error) {
	// Parse user ID
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Get full details
	details, err := s.repository.GetUserFullDetails(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user details: %w", err)
	}

	response := &authModels.UserResponse{
		UserID:    details.User.ID.String(),
		Status:    details.User.Status,
		UserType:  details.User.UserType,
		CreatedAt: details.User.CreatedAt,
	}

	if details.Identifier != nil {
		response.Username = details.Identifier.Identifier
	}

	if details.Profile != nil {
		response.Email = details.Profile.Email
		response.Firstname = details.Profile.Firstname
		response.Lastname = details.Profile.Lastname
		response.Birthday = details.Profile.Birthday
		response.Locale = details.Profile.Locale
		response.Metadata = details.Profile.Metadata
	}

	return response, nil
}

// UpdateProfile updates user profile
func (s *DefaultAuthService) UpdateProfile(ctx context.Context, req *authModels.UpdateProfileRequest) (*authModels.UserResponse, error) {
	s.log.Infof("DefaultAuth UpdateProfile: userID=%s", req.UserID)

	// Parse user ID
	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}

	// Get existing profile
	profile, err := s.repository.GetUserProfileByUserID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Update fields
	if req.Firstname != "" {
		profile.Firstname = req.Firstname
	}
	if req.Lastname != "" {
		profile.Lastname = req.Lastname
	}
	if req.Email != "" {
		profile.Email = req.Email
	}
	if req.Birthday != "" {
		t, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			profile.Birthday = &t
		}
	}
	if req.Locale != "" {
		profile.Locale = req.Locale
	}
	if req.Metadata != nil {
		profile.Metadata = req.Metadata
	}

	// Save profile
	_, err = s.repository.UpdateUserProfile(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Return updated profile
	return s.GetProfile(ctx, req.UserID)
}
