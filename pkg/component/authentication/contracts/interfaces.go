package contracts

import (
	"context"
	"time"

	"github.com/phatnt199/go-infra/pkg/component/authentication/models"
)

// UserAuthInfo contains minimum info needed for authentication
type UserAuthInfo struct {
	UserID       string
	Username     string
	Email        string
	PasswordHash string
	Firstname    string
	Lastname     string
	Birthday     *time.Time
	Locale       string
	Roles        []string
	Status       string // active, inactive, banned, etc.
	UserType     string
	CreatedAt    time.Time
	Metadata     map[string]interface{} // Custom fields
}

// CreateUserRequest represents user creation request
type CreateUserRequest struct {
	Username     string
	Email        string
	PasswordHash string
	Firstname    string
	Lastname     string
	Birthday     *time.Time
	Locale       string
	Roles        []string
	UserType     string
	ClientID     string
	Metadata     map[string]interface{}
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	UserID    string
	Firstname string
	Lastname  string
	Email     string
	Birthday  *time.Time
	Locale    string
	Metadata  map[string]interface{}
}

// IUserProvider defines how to fetch and manage user data
// Users implement this to connect auth to their user model
type IUserProvider interface {
	// GetUserByIdentifier finds a user by identifier (username, email, phone, etc.)
	GetUserByIdentifier(ctx context.Context, scheme, value string) (*UserAuthInfo, error)

	// CreateUser creates a new user (for sign up)
	CreateUser(ctx context.Context, req *CreateUserRequest) (*UserAuthInfo, error)

	// UpdatePassword updates user password
	UpdatePassword(ctx context.Context, userID, newPasswordHash string) error

	// GetUserByID gets user info by ID (for token refresh)
	GetUserByID(ctx context.Context, userID string) (*UserAuthInfo, error)

	// UpdateProfile updates user profile information
	UpdateProfile(ctx context.Context, req *UpdateUserRequest) (*UserAuthInfo, error)
}

// IAuthService defines authentication operations
type IAuthService interface {
	// SignUp creates a new user with credentials and profile
	SignUp(ctx context.Context, req *models.SignUpRequest) (*models.AuthResponse, error)

	// SignIn authenticates a user and returns JWT tokens
	SignIn(ctx context.Context, req *models.SignInRequest) (*models.AuthResponse, error)

	// ChangePassword changes user password
	ChangePassword(ctx context.Context, req *models.ChangePasswordRequest) error

	// GetProfile gets user profile information
	GetProfile(ctx context.Context, userID string) (*models.UserResponse, error)

	// UpdateProfile updates user profile information
	UpdateProfile(ctx context.Context, req *models.UpdateProfileRequest) (*models.UserResponse, error)
}

// ITokenService defines token management operations
type ITokenService interface {
	// GenerateAccessToken generates an access token for a user
	GenerateAccessToken(userID, username string, roles []string) (string, error)

	// GenerateRefreshToken generates a refresh token for a user
	GenerateRefreshToken(userID, username string) (string, error)

	// GenerateTokenPair generates both access and refresh tokens
	GenerateTokenPair(userID, username string, roles []string) (accessToken, refreshToken string, err error)

	// ValidateAccessToken validates an access token and returns claims
	ValidateAccessToken(token string) (userID, username string, roles []string, err error)

	// ValidateRefreshToken validates a refresh token
	ValidateRefreshToken(token string) (userID, username string, err error)

	// GetAccessTokenExpiry returns access token expiry duration
	GetAccessTokenExpiry() time.Duration

	// GetRefreshTokenExpiry returns refresh token expiry duration
	GetRefreshTokenExpiry() time.Duration
}

// IPasswordHasher defines password hashing operations
type IPasswordHasher interface {
	// HashPassword hashes a password
	HashPassword(password string) (string, error)

	// ComparePassword compares a password with a hash
	ComparePassword(password, hash string) (bool, error)
}
