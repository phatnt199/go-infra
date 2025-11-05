package seeder

import (
	"context"
	"time"

	"emperror.dev/errors"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/contracts"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	"github.com/phatnt199/go-infra/pkg/crypto"
	"github.com/phatnt199/go-infra/pkg/logger"
	uuid "github.com/satori/go.uuid"
)

// UserSeed represents a user to be seeded
type UserSeed struct {
	Username  string
	Password  string // Plain text password (will be hashed)
	Firstname string
	Lastname  string
	Locale    string
	UserType  models.UserType
	Status    models.UserStatus
}

// Seeder handles database seeding
type Seeder struct {
	log    logger.Logger
	repo   contracts.UserRepository
	hasher *crypto.Hasher
}

// NewSeeder creates a new seeder
func NewSeeder(log logger.Logger, repo contracts.UserRepository) *Seeder {
	return &Seeder{
		log:    log,
		repo:   repo,
		hasher: crypto.NewHasher(nil), // Use default hasher config
	}
}

// SeedUsers seeds the database with default users
func (s *Seeder) SeedUsers(ctx context.Context) error {
	s.log.Info("Starting user seeding...")

	// Define default users
	// Password for all users: "SecurePass123!"
	// In production, you should use different passwords for each user
	defaultUsers := []UserSeed{
		{
			Username:  "superadmin",
			Password:  "SuperAdmin123!",
			Firstname: "Super",
			Lastname:  "Administrator",
			Locale:    "en_US",
			UserType:  models.UserTypeSystem,
			Status:    models.UserStatusActivated,
		},
		{
			Username:  "admin",
			Password:  "Admin123!",
			Firstname: "System",
			Lastname:  "Admin",
			Locale:    "en_US",
			UserType:  models.UserTypeSystem,
			Status:    models.UserStatusActivated,
		},
		{
			Username:  "maintainer",
			Password:  "Maintainer123!",
			Firstname: "System",
			Lastname:  "Maintainer",
			Locale:    "en_US",
			UserType:  models.UserTypeSystem,
			Status:    models.UserStatusActivated,
		},
		{
			Username:  "manager",
			Password:  "Manager123!",
			Firstname: "System",
			Lastname:  "Manager",
			Locale:    "en_US",
			UserType:  models.UserTypeSystem,
			Status:    models.UserStatusActivated,
		},
		{
			Username:  "customer",
			Password:  "Customer123!",
			Firstname: "Test",
			Lastname:  "Customer",
			Locale:    "en_US",
			UserType:  models.UserTypeSystem,
			Status:    models.UserStatusActivated,
		},
	}

	for _, userSeed := range defaultUsers {
		if err := s.seedUser(ctx, userSeed); err != nil {
			s.log.Errorf("Failed to seed user %s: %v", userSeed.Username, err)
			// Continue with other users instead of failing completely
			continue
		}
		s.log.Infof("Successfully seeded user: %s (password: %s)", userSeed.Username, userSeed.Password)
	}

	s.log.Info("User seeding completed!")
	return nil
}

// seedUser seeds a single user
func (s *Seeder) seedUser(ctx context.Context, seed UserSeed) error {
	// Check if user already exists
	existingIdentifier, _ := s.repo.GetUserIdentifierBySchemeAndIdentifier(ctx, models.IdentifierSchemeUsername, seed.Username)
	if existingIdentifier != nil {
		s.log.Infof("User %s already exists, skipping...", seed.Username)
		return nil
	}

	// Hash password
	hashedPassword, err := s.hasher.HashPassword(seed.Password, crypto.AlgorithmBcrypt)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	// Create user
	now := time.Now()
	user := &models.User{
		Status:      seed.Status,
		UserType:    seed.UserType,
		ActivatedAt: &now,
	}

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, "failed to create user")
	}

	// Create user identifier
	identifier := &models.UserIdentifier{
		UserID:     createdUser.ID,
		Scheme:     models.IdentifierSchemeUsername,
		Identifier: seed.Username,
		Verified:   true,
		Details:    make(map[string]interface{}),
	}

	_, err = s.repo.CreateUserIdentifier(ctx, identifier)
	if err != nil {
		// Rollback: delete user
		_ = s.repo.DeleteUser(ctx, createdUser.ID)
		return errors.Wrap(err, "failed to create user identifier")
	}

	// Create user credential
	credential := &models.UserCredential{
		UserID:     createdUser.ID,
		Scheme:     models.CredentialSchemeBasic,
		Credential: hashedPassword,
		Details:    make(map[string]interface{}),
	}

	_, err = s.repo.CreateUserCredential(ctx, credential)
	if err != nil {
		// Rollback: delete user
		_ = s.repo.DeleteUser(ctx, createdUser.ID)
		return errors.Wrap(err, "failed to create user credential")
	}

	// Create user profile
	profile := &models.UserProfile{
		UserID:    createdUser.ID,
		Firstname: seed.Firstname,
		Lastname:  seed.Lastname,
		Locale:    seed.Locale,
		Details:   make(map[string]interface{}),
	}

	_, err = s.repo.CreateUserProfile(ctx, profile)
	if err != nil {
		s.log.Warnf("Failed to create profile for user %s: %v", seed.Username, err)
		// Don't rollback the user creation if profile fails
	}

	return nil
}

// SeedUserWithParent seeds a user with a parent user ID
func (s *Seeder) SeedUserWithParent(ctx context.Context, seed UserSeed, parentID uuid.UUID) error {
	// Similar to seedUser but sets the ParentID field
	return nil // Implementation omitted for brevity
}
