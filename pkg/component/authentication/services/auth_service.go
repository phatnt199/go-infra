package services

import (
	"context"
	"fmt"
	"time"

	"github.com/phatnt199/go-infra/pkg/component/authentication/config"
	"github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"github.com/phatnt199/go-infra/pkg/component/authentication/models"
	"github.com/phatnt199/go-infra/pkg/logger"
)

// authService implements IAuthService
type authService struct {
	userProvider   contracts.IUserProvider
	tokenService   contracts.ITokenService
	passwordHasher contracts.IPasswordHasher
	config         *config.AuthConfig
	logger         logger.Logger
}

// NewAuthService creates a new auth service
func NewAuthService(
	userProvider contracts.IUserProvider,
	tokenService contracts.ITokenService,
	passwordHasher contracts.IPasswordHasher,
	config *config.AuthConfig,
	logger logger.Logger,
) contracts.IAuthService {
	return &authService{
		userProvider:   userProvider,
		tokenService:   tokenService,
		passwordHasher: passwordHasher,
		config:         config,
		logger:         logger,
	}
}

// SignUp creates a new user with credentials and profile
func (s *authService) SignUp(ctx context.Context, req *models.SignUpRequest) (*models.AuthResponse, error) {
	// Validate identifier scheme
	if !s.isValidIdentifierScheme(req.Identifier.Scheme) {
		return nil, fmt.Errorf("unsupported identifier scheme: %s", req.Identifier.Scheme)
	}

	// Validate credential scheme
	if !s.isValidCredentialScheme(req.Credential.Scheme) {
		return nil, fmt.Errorf("unsupported credential scheme: %s", req.Credential.Scheme)
	}

	// Validate password length
	if len(req.GetCredentialValue()) < s.config.MinPasswordLength {
		return nil, fmt.Errorf("password must be at least %d characters", s.config.MinPasswordLength)
	}

	// Check if user already exists
	existingUser, _ := s.userProvider.GetUserByIdentifier(ctx, req.Identifier.Scheme, req.GetIdentifierValue())
	if existingUser != nil {
		return nil, fmt.Errorf("user with this %s already exists", req.Identifier.Scheme)
	}

	// Hash password
	passwordHash, err := s.passwordHasher.HashPassword(req.GetCredentialValue())
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Parse birthday if provided
	var birthday *time.Time
	if req.Birthday != "" {
		t, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			birthday = &t
		}
	}

	// Create user
	createReq := &contracts.CreateUserRequest{
		Username:     req.GetIdentifierValue(),
		Email:        req.Email,
		PasswordHash: passwordHash,
		Firstname:    req.Firstname,
		Lastname:     req.Lastname,
		Birthday:     birthday,
		Locale:       req.Locale,
		Roles:        []string{"user"},
		UserType:     "user",
		ClientID:     req.ClientID,
	}

	user, err := s.userProvider.CreateUser(ctx, createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	accessToken, refreshToken, err := s.tokenService.GenerateTokenPair(user.UserID, user.Username, user.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Build response
	return &models.AuthResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Status:    user.Status,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		AccessToken: &models.TokenInfo{
			Value:     accessToken,
			Scheme:    "bearer",
			Type:      "access",
			ExpiresAt: time.Now().Add(s.tokenService.GetAccessTokenExpiry()),
		},
		RefreshToken: &models.TokenInfo{
			Value:     refreshToken,
			Scheme:    "bearer",
			Type:      "refresh",
			ExpiresAt: time.Now().Add(s.tokenService.GetRefreshTokenExpiry()),
		},
	}, nil
}

// SignIn authenticates a user and returns JWT tokens
func (s *authService) SignIn(ctx context.Context, req *models.SignInRequest) (*models.AuthResponse, error) {
	// Validate identifier scheme
	if !s.isValidIdentifierScheme(req.Identifier.Scheme) {
		return nil, fmt.Errorf("unsupported identifier scheme: %s", req.Identifier.Scheme)
	}

	// Validate credential scheme
	if !s.isValidCredentialScheme(req.Credential.Scheme) {
		return nil, fmt.Errorf("unsupported credential scheme: %s", req.Credential.Scheme)
	}

	// Get user by identifier
	user, err := s.userProvider.GetUserByIdentifier(ctx, req.Identifier.Scheme, req.GetIdentifierValue())
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check user status
	if user.Status != "active" {
		return nil, fmt.Errorf("user account is not active")
	}

	// Verify password
	valid, err := s.passwordHasher.ComparePassword(req.GetCredentialValue(), user.PasswordHash)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate tokens
	accessToken, refreshToken, err := s.tokenService.GenerateTokenPair(user.UserID, user.Username, user.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Build response
	return &models.AuthResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Status:    user.Status,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		AccessToken: &models.TokenInfo{
			Value:     accessToken,
			Scheme:    "bearer",
			Type:      "access",
			ExpiresAt: time.Now().Add(s.tokenService.GetAccessTokenExpiry()),
		},
		RefreshToken: &models.TokenInfo{
			Value:     refreshToken,
			Scheme:    "bearer",
			Type:      "refresh",
			ExpiresAt: time.Now().Add(s.tokenService.GetRefreshTokenExpiry()),
		},
	}, nil
}

// ChangePassword changes user password
func (s *authService) ChangePassword(ctx context.Context, req *models.ChangePasswordRequest) error {
	// Validate credential schemes
	if !s.isValidCredentialScheme(req.OldCredential.Scheme) {
		return fmt.Errorf("unsupported credential scheme: %s", req.OldCredential.Scheme)
	}
	if !s.isValidCredentialScheme(req.NewCredential.Scheme) {
		return fmt.Errorf("unsupported credential scheme: %s", req.NewCredential.Scheme)
	}

	// Validate new password length
	if len(req.GetNewPassword()) < s.config.MinPasswordLength {
		return fmt.Errorf("password must be at least %d characters", s.config.MinPasswordLength)
	}

	// Get user
	user, err := s.userProvider.GetUserByID(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	valid, err := s.passwordHasher.ComparePassword(req.GetOldPassword(), user.PasswordHash)
	if err != nil || !valid {
		return fmt.Errorf("invalid old password")
	}

	// Hash new password
	newPasswordHash, err := s.passwordHasher.HashPassword(req.GetNewPassword())
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	err = s.userProvider.UpdatePassword(ctx, req.UserID, newPasswordHash)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// GetProfile gets user profile information
func (s *authService) GetProfile(ctx context.Context, userID string) (*models.UserResponse, error) {
	user, err := s.userProvider.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &models.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Birthday:  user.Birthday,
		Locale:    user.Locale,
		Status:    user.Status,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		Metadata:  user.Metadata,
	}, nil
}

// UpdateProfile updates user profile information
func (s *authService) UpdateProfile(ctx context.Context, req *models.UpdateProfileRequest) (*models.UserResponse, error) {
	// Parse birthday if provided
	var birthday *time.Time
	if req.Birthday != "" {
		t, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			birthday = &t
		}
	}

	updateReq := &contracts.UpdateUserRequest{
		UserID:    req.UserID,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Email:     req.Email,
		Birthday:  birthday,
		Locale:    req.Locale,
		Metadata:  req.Metadata,
	}

	user, err := s.userProvider.UpdateProfile(ctx, updateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return &models.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Birthday:  user.Birthday,
		Locale:    user.Locale,
		Status:    user.Status,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		Metadata:  user.Metadata,
	}, nil
}

// isValidIdentifierScheme checks if identifier scheme is supported
func (s *authService) isValidIdentifierScheme(scheme string) bool {
	for _, supportedScheme := range s.config.SupportedIdentifierSchemes {
		if scheme == supportedScheme {
			return true
		}
	}
	return false
}

// isValidCredentialScheme checks if credential scheme is supported
func (s *authService) isValidCredentialScheme(scheme string) bool {
	for _, supportedScheme := range s.config.SupportedCredentialSchemes {
		if scheme == supportedScheme {
			return true
		}
	}
	return false
}
