package services

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/shared/utils"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/contracts"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/dtos/v1"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/mappers"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	"github.com/phatnt199/go-infra/pkg/logger"
	"gorm.io/gorm"
)

type authService struct {
	log      logger.Logger
	userRepo contracts.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(log logger.Logger, userRepo contracts.UserRepository) contracts.AuthService {
	return &authService{
		log:      log,
		userRepo: userRepo,
	}
}

// SignUp creates a new user account
func (s *authService) SignUp(ctx context.Context, req *dtosv1.SignUpRequest) (*dtosv1.AuthResponse, error) {
	// Check if username already exists
	existingIdentifier, err := s.userRepo.GetUserIdentifierBySchemeAndIdentifier(ctx, models.IdentifierSchemeUsername, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrap(err, "failed to check existing username")
	}
	if existingIdentifier != nil {
		return nil, errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash password")
	}

	// Create user
	now := time.Now()
	user := &models.User{
		Status:      models.UserStatusActivated,
		UserType:    models.UserTypeSystem,
		ActivatedAt: &now,
	}

	createdUser, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	// Create user identifier
	identifier := &models.UserIdentifier{
		UserID:     createdUser.ID,
		Scheme:     models.IdentifierSchemeUsername,
		Identifier: req.Username,
		Verified:   true,
		Details:    make(map[string]interface{}),
	}

	createdIdentifier, err := s.userRepo.CreateUserIdentifier(ctx, identifier)
	if err != nil {
		// Rollback: delete user
		_ = s.userRepo.DeleteUser(ctx, createdUser.ID)
		return nil, errors.Wrap(err, "failed to create user identifier")
	}

	// Create user credential
	credential := &models.UserCredential{
		UserID:     createdUser.ID,
		Scheme:     models.CredentialSchemeBasic,
		Credential: hashedPassword,
		Details:    make(map[string]interface{}),
	}

	_, err = s.userRepo.CreateUserCredential(ctx, credential)
	if err != nil {
		// Rollback: delete user and identifier
		_ = s.userRepo.DeleteUser(ctx, createdUser.ID)
		return nil, errors.Wrap(err, "failed to create user credential")
	}

	// Create user profile
	locale := req.Locale
	if locale == "" {
		locale = "en_US"
	}

	var birthday *time.Time
	if req.Birthday != "" {
		bd, err := mappers.ParseISODate(req.Birthday)
		if err != nil {
			return nil, errors.Wrap(err, "invalid birthday format")
		}
		birthday = bd
	}

	profile := &models.UserProfile{
		UserID:    createdUser.ID,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Birthday:  birthday,
		Locale:    locale,
		Details:   make(map[string]interface{}),
	}

	createdProfile, err := s.userRepo.CreateUserProfile(ctx, profile)
	if err != nil {
		// Rollback: delete user
		_ = s.userRepo.DeleteUser(ctx, createdUser.ID)
		return nil, errors.Wrap(err, "failed to create user profile")
	}

	// Generate token (for now, just return a simple message)
	// In production, you should generate a JWT token here
	token := "token_placeholder"

	s.log.Infof("User signed up successfully: %s", createdUser.ID)

	return mappers.ToAuthResponse(createdUser, createdIdentifier, createdProfile, token), nil
}

// SignIn authenticates a user
func (s *authService) SignIn(ctx context.Context, req *dtosv1.SignInRequest) (*dtosv1.AuthResponse, error) {
	// Get user by username
	identifier, err := s.userRepo.GetUserIdentifierBySchemeAndIdentifier(ctx, models.IdentifierSchemeUsername, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, errors.Wrap(err, "failed to get user identifier")
	}

	// Get user
	user, err := s.userRepo.GetUserByID(ctx, identifier.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	// Check if user is activated
	if user.Status != models.UserStatusActivated {
		return nil, errors.New("user account is not activated")
	}

	// Get credential
	credential, err := s.userRepo.GetUserCredentialByUserIDAndScheme(ctx, user.ID, models.CredentialSchemeBasic)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user credential")
	}

	// Verify password
	if err := utils.ComparePassword(credential.Credential, req.Password); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	_, err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		s.log.Warnf("Failed to update last login time: %v", err)
	}

	// Get profile
	profile, _ := s.userRepo.GetUserProfileByUserID(ctx, user.ID)

	// Generate token (for now, just return a simple message)
	// In production, you should generate a JWT token here
	token := "token_placeholder"

	s.log.Infof("User signed in successfully: %s", user.ID)

	return mappers.ToAuthResponse(user, identifier, profile, token), nil
}

// ChangePassword changes user password
func (s *authService) ChangePassword(ctx context.Context, req *dtosv1.ChangePasswordRequest) error {
	// Get current credential
	credential, err := s.userRepo.GetUserCredentialByUserIDAndScheme(ctx, req.UserID, models.CredentialSchemeBasic)
	if err != nil {
		return errors.Wrap(err, "failed to get user credential")
	}

	// Verify old password
	if err := utils.ComparePassword(credential.Credential, req.OldPassword); err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.Wrap(err, "failed to hash new password")
	}

	// Update credential
	credential.Credential = hashedPassword
	_, err = s.userRepo.UpdateUserCredential(ctx, credential)
	if err != nil {
		return errors.Wrap(err, "failed to update credential")
	}

	s.log.Infof("Password changed successfully for user: %s", req.UserID)

	return nil
}
