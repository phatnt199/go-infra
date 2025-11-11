package services

import (
	"context"
	"fmt"
	"time"

	authConfig "github.com/phatnt199/go-infra/pkg/component/authentication/config"
	authContracts "github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	authModels "github.com/phatnt199/go-infra/pkg/component/authentication/models"
	authServices "github.com/phatnt199/go-infra/pkg/component/authentication/services"
	"github.com/phatnt199/go-infra/pkg/logger"

	"github.com/phatnt199/go-infra/examples/authentication-service/internal/custom-auth/models"
	"github.com/phatnt199/go-infra/examples/authentication-service/internal/custom-auth/validators"
	uuid "github.com/satori/go.uuid"
)

// CustomAuthService implements IAuthService with custom business logic
type CustomAuthService struct {
	userProvider        authContracts.IUserProvider
	tokenService        authContracts.ITokenService
	passwordHasher      authServices.IPasswordHasher
	config              *authConfig.Config
	logger              logger.Logger
	auditLogger         *AuditLogger
	loginAttemptManager *LoginAttemptManager
	passwordValidator   *validators.PasswordValidator
}

// NewCustomAuthService creates a new custom auth service
func NewCustomAuthService(
	userProvider authContracts.IUserProvider,
	tokenService authContracts.ITokenService,
	passwordHasher authServices.IPasswordHasher,
	cfg *authConfig.Config,
	logger logger.Logger,
	auditLogger *AuditLogger,
	loginAttemptManager *LoginAttemptManager,
) authContracts.IAuthService {
	return &CustomAuthService{
		userProvider:        userProvider,
		tokenService:        tokenService,
		passwordHasher:      passwordHasher,
		config:              cfg,
		logger:              logger,
		auditLogger:         auditLogger,
		loginAttemptManager: loginAttemptManager,
		passwordValidator:   validators.NewPasswordValidator(),
	}
}

// SignUp creates a new user with custom validation
func (s *CustomAuthService) SignUp(ctx context.Context, req *authModels.SignUpRequest) (*authModels.AuthResponse, error) {
	s.logger.Infof("Custom SignUp: user=%s", req.GetIdentifierValue())

	// Validate identifier scheme
	if !s.isValidIdentifierScheme(req.Identifier.Scheme) {
		return nil, fmt.Errorf("unsupported identifier scheme: %s", req.Identifier.Scheme)
	}

	// Validate credential scheme
	if !s.isValidCredentialScheme(req.Credential.Scheme) {
		return nil, fmt.Errorf("unsupported credential scheme: %s", req.Credential.Scheme)
	}

	// Custom password validation with enhanced rules
	if err := s.passwordValidator.ValidatePassword(req.GetCredentialValue()); err != nil {
		s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
			Username:  req.GetIdentifierValue(),
			Action:    models.AuditActionSignUp,
			Success:   false,
			Details:   fmt.Sprintf("Password validation failed: %v", err),
			CreatedAt: time.Now(),
		})
		return nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Check if user already exists
	existingUser, _ := s.userProvider.GetUserByIdentifier(ctx, req.Identifier.Scheme, req.GetIdentifierValue())
	if existingUser != nil {
		s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
			Username:  req.GetIdentifierValue(),
			Action:    models.AuditActionSignUp,
			Success:   false,
			Details:   "User already exists",
			CreatedAt: time.Now(),
		})
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
	createReq := &authContracts.CreateUserRequest{
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
		s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
			Username:  req.GetIdentifierValue(),
			Action:    models.AuditActionSignUp,
			Success:   false,
			Details:   fmt.Sprintf("Failed to create user: %v", err),
			CreatedAt: time.Now(),
		})
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	accessToken, refreshToken, err := s.tokenService.GenerateTokenPair(user.UserID, user.Username, user.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Log successful signup
	s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
		UserID:    user.UserID,
		Username:  user.Username,
		Action:    models.AuditActionSignUp,
		Success:   true,
		Details:   "User successfully registered",
		CreatedAt: time.Now(),
	})

	// Build response
	return &authModels.AuthResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Status:    user.Status,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		AccessToken: &authModels.TokenInfo{
			Value:     accessToken,
			Scheme:    "bearer",
			Type:      "access",
			ExpiresAt: time.Now().Add(s.tokenService.GetAccessTokenExpiry()),
		},
		RefreshToken: &authModels.TokenInfo{
			Value:     refreshToken,
			Scheme:    "bearer",
			Type:      "refresh",
			ExpiresAt: time.Now().Add(s.tokenService.GetRefreshTokenExpiry()),
		},
	}, nil
}

// SignIn authenticates with login attempt tracking
func (s *CustomAuthService) SignIn(ctx context.Context, req *authModels.SignInRequest) (*authModels.AuthResponse, error) {
	identifier := req.GetIdentifierValue()
	s.logger.Infof("Custom SignIn attempt: user=%s", identifier)

	// Validate identifier scheme
	if !s.isValidIdentifierScheme(req.Identifier.Scheme) {
		return nil, fmt.Errorf("unsupported identifier scheme: %s", req.Identifier.Scheme)
	}

	// Validate credential scheme
	if !s.isValidCredentialScheme(req.Credential.Scheme) {
		return nil, fmt.Errorf("unsupported credential scheme: %s", req.Credential.Scheme)
	}

	// Check if account is locked
	locked, err := s.loginAttemptManager.IsAccountLocked(ctx, identifier)
	if err != nil {
		s.logger.Errorf("Error checking account lock status: %v", err)
	}

	if locked {
		lockTime, _ := s.loginAttemptManager.GetLockTimeRemaining(ctx, identifier)
		s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
			Username:  identifier,
			Action:    models.AuditActionAccountLocked,
			Success:   false,
			Details:   fmt.Sprintf("Account locked, remaining time: %v", lockTime),
			CreatedAt: time.Now(),
		})
		return nil, fmt.Errorf("account is locked due to too many failed login attempts. Please try again in %v", lockTime.Round(time.Second))
	}

	// Get user by identifier
	user, err := s.userProvider.GetUserByIdentifier(ctx, req.Identifier.Scheme, identifier)
	if err != nil {
		// Record failed attempt
		s.loginAttemptManager.RecordLoginAttempt(ctx, identifier, false)

		remaining, _ := s.loginAttemptManager.GetRemainingAttempts(ctx, identifier)
		s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
			Username:  identifier,
			Action:    models.AuditActionFailedLogin,
			Success:   false,
			Details:   fmt.Sprintf("User not found. Remaining attempts: %d", remaining),
			CreatedAt: time.Now(),
		})

		if remaining > 0 {
			return nil, fmt.Errorf("invalid credentials. %d attempts remaining", remaining)
		}
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check user status
	if user.Status != "active" {
		s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
			UserID:    user.UserID,
			Username:  user.Username,
			Action:    models.AuditActionFailedLogin,
			Success:   false,
			Details:   fmt.Sprintf("Account not active. Status: %s", user.Status),
			CreatedAt: time.Now(),
		})
		return nil, fmt.Errorf("user account is not active")
	}

	// Verify password
	valid, err := s.passwordHasher.ComparePassword(req.GetCredentialValue(), user.PasswordHash)
	if err != nil || !valid {
		// Record failed attempt
		s.loginAttemptManager.RecordLoginAttempt(ctx, identifier, false)

		remaining, _ := s.loginAttemptManager.GetRemainingAttempts(ctx, identifier)
		s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
			UserID:    user.UserID,
			Username:  user.Username,
			Action:    models.AuditActionFailedLogin,
			Success:   false,
			Details:   fmt.Sprintf("Invalid password. Remaining attempts: %d", remaining),
			CreatedAt: time.Now(),
		})

		if remaining > 0 {
			return nil, fmt.Errorf("invalid credentials. %d attempts remaining", remaining)
		}
		return nil, fmt.Errorf("invalid credentials")
	}

	// Successful login - reset attempts
	s.loginAttemptManager.RecordLoginAttempt(ctx, identifier, true)

	// Generate tokens
	accessToken, refreshToken, err := s.tokenService.GenerateTokenPair(user.UserID, user.Username, user.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Log successful signin
	s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
		UserID:    user.UserID,
		Username:  user.Username,
		Action:    models.AuditActionSignIn,
		Success:   true,
		Details:   "User successfully signed in",
		CreatedAt: time.Now(),
	})

	// Build response
	return &authModels.AuthResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Status:    user.Status,
		UserType:  user.UserType,
		CreatedAt: user.CreatedAt,
		AccessToken: &authModels.TokenInfo{
			Value:     accessToken,
			Scheme:    "bearer",
			Type:      "access",
			ExpiresAt: time.Now().Add(s.tokenService.GetAccessTokenExpiry()),
		},
		RefreshToken: &authModels.TokenInfo{
			Value:     refreshToken,
			Scheme:    "bearer",
			Type:      "refresh",
			ExpiresAt: time.Now().Add(s.tokenService.GetRefreshTokenExpiry()),
		},
	}, nil
}

// ChangePassword changes user password with audit logging
func (s *CustomAuthService) ChangePassword(ctx context.Context, req *authModels.ChangePasswordRequest) error {
	s.logger.Infof("Custom ChangePassword: userID=%s", req.UserID)

	// Validate credential schemes
	if !s.isValidCredentialScheme(req.OldCredential.Scheme) {
		return fmt.Errorf("unsupported credential scheme: %s", req.OldCredential.Scheme)
	}
	if !s.isValidCredentialScheme(req.NewCredential.Scheme) {
		return fmt.Errorf("unsupported credential scheme: %s", req.NewCredential.Scheme)
	}

	// Custom password validation
	if err := s.passwordValidator.ValidatePassword(req.GetNewPassword()); err != nil {
		return fmt.Errorf("new password validation failed: %w", err)
	}

	// Get user
	user, err := s.userProvider.GetUserByID(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	valid, err := s.passwordHasher.ComparePassword(req.GetOldPassword(), user.PasswordHash)
	if err != nil || !valid {
		s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
			UserID:    user.UserID,
			Username:  user.Username,
			Action:    models.AuditActionChangePassword,
			Success:   false,
			Details:   "Invalid old password",
			CreatedAt: time.Now(),
		})
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
		s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
			UserID:    user.UserID,
			Username:  user.Username,
			Action:    models.AuditActionChangePassword,
			Success:   false,
			Details:   fmt.Sprintf("Failed to update password: %v", err),
			CreatedAt: time.Now(),
		})
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Log successful password change
	s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
		UserID:    user.UserID,
		Username:  user.Username,
		Action:    models.AuditActionChangePassword,
		Success:   true,
		Details:   "Password successfully changed",
		CreatedAt: time.Now(),
	})

	return nil
}

// GetProfile gets user profile information
func (s *CustomAuthService) GetProfile(ctx context.Context, userID string) (*authModels.UserResponse, error) {
	user, err := s.userProvider.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &authModels.UserResponse{
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

// UpdateProfile updates user profile information with audit logging
func (s *CustomAuthService) UpdateProfile(ctx context.Context, req *authModels.UpdateProfileRequest) (*authModels.UserResponse, error) {
	s.logger.Infof("Custom UpdateProfile: userID=%s", req.UserID)

	// Parse birthday if provided
	var birthday *time.Time
	if req.Birthday != "" {
		t, err := time.Parse("2006-01-02", req.Birthday)
		if err == nil {
			birthday = &t
		}
	}

	updateReq := &authContracts.UpdateUserRequest{
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
		s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
			UserID:    req.UserID,
			Action:    models.AuditActionUpdateProfile,
			Success:   false,
			Details:   fmt.Sprintf("Failed to update profile: %v", err),
			CreatedAt: time.Now(),
		})
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	// Log successful profile update
	s.auditLogger.LogAuthEvent(ctx, models.AuditLog{
		UserID:    user.UserID,
		Username:  user.Username,
		Action:    models.AuditActionUpdateProfile,
		Success:   true,
		Details:   "Profile successfully updated",
		CreatedAt: time.Now(),
	})

	return &authModels.UserResponse{
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

// GetAuditLogs retrieves audit logs for a user
func (s *CustomAuthService) GetAuditLogs(ctx context.Context, userID string, page, pageSize int) (*models.PaginatedAuditLogsResponse, error) {
	offset := (page - 1) * pageSize
	logs, total, err := s.auditLogger.GetUserAuditLogs(ctx, userID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	// Convert to response models
	responses := make([]models.AuditLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = models.AuditLogResponse{
			ID:        log.ID.String(),
			UserID:    log.UserID,
			Username:  log.Username,
			Action:    string(log.Action),
			Success:   log.Success,
			IPAddress: log.IPAddress,
			Details:   log.Details,
			CreatedAt: log.CreatedAt,
		}
	}

	return &models.PaginatedAuditLogsResponse{
		Data:       responses,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}, nil
}

// Helper methods

func (s *CustomAuthService) isValidIdentifierScheme(scheme string) bool {
	for _, supportedScheme := range s.config.Schemes.IdentifierSchemes {
		if scheme == supportedScheme {
			return true
		}
	}
	return false
}

func (s *CustomAuthService) isValidCredentialScheme(scheme string) bool {
	for _, supportedScheme := range s.config.Schemes.CredentialSchemes {
		if scheme == supportedScheme {
			return true
		}
	}
	return false
}

// generateSessionID generates a unique session ID
func generateSessionID() string {
	return uuid.NewV4().String()
}
