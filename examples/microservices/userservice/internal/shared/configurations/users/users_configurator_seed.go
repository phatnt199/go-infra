package users

import (
	"time"

	"emperror.dev/errors"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/data/datamodels"
	"github.com/phatnt199/go-infra/pkg/crypto"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (usc *UsersServiceConfigurator) seedUsers(db *gorm.DB) error {
	err := seedDataManually(db)
	if err != nil {
		return err
	}
	return nil
}

// seedDataManually seeds 5 users with different roles
// Passwords: All users have password "Password123!" (shown here for development)
func seedDataManually(db *gorm.DB) error {
	var count int64

	db.Model(&datamodels.UserDataModel{}).Count(&count)
	if count > 0 {
		return nil
	}

	// Hash the common password
	hasher := crypto.NewHasher(nil)
	hashedPassword, err := hasher.HashPassword("Password123!", crypto.AlgorithmBcrypt)
	if err != nil {
		return errors.Wrap(err, "Error hashing password")
	}

	now := time.Now()

	// Define 5 users with their details
	seedUsers := []struct {
		username  string
		firstname string
		lastname  string
		role      string // for display, not used in DB yet
	}{
		{username: "superadmin", firstname: "Super", lastname: "Admin", role: "Super Administrator"},
		{username: "admin", firstname: "Admin", lastname: "User", role: "Administrator"},
		{username: "maintainer", firstname: "System", lastname: "Maintainer", role: "Maintainer"},
		{username: "manager", firstname: "Project", lastname: "Manager", role: "Manager"},
		{username: "customer", firstname: "Regular", lastname: "Customer", role: "Customer"},
	}

	// Create users, identifiers, credentials, and profiles in transaction
	return db.Transaction(func(tx *gorm.DB) error {
		for _, su := range seedUsers {
			userID := uuid.NewV4()

			// Create user
			user := datamodels.UserDataModel{
				ID:          userID,
				Status:      datamodels.UserStatusActivated, // Use enum constant
				UserType:    datamodels.UserTypeSystem,      // Use enum constant
				ActivatedAt: &now,
				CreatedAt:   now,
				ModifiedAt:  now,
			}

			if err := tx.Create(&user).Error; err != nil {
				return errors.Wrapf(err, "Error creating user %s", su.username)
			}

			// Create identifier
			identifier := datamodels.UserIdentifierDataModel{
				ID:         uuid.NewV4(),
				UserID:     userID,
				Scheme:     "username",
				Identifier: su.username,
				Verified:   true,
				Details:    datamodels.JSONB{},
				CreatedAt:  now,
				ModifiedAt: now,
			}

			if err := tx.Create(&identifier).Error; err != nil {
				return errors.Wrapf(err, "Error creating identifier for %s", su.username)
			}

			// Create credential
			credential := datamodels.UserCredentialDataModel{
				ID:         uuid.NewV4(),
				UserID:     userID,
				Scheme:     "basic",
				Credential: hashedPassword,
				Details:    datamodels.JSONB{},
				CreatedAt:  now,
				ModifiedAt: now,
			}

			if err := tx.Create(&credential).Error; err != nil {
				return errors.Wrapf(err, "Error creating credential for %s", su.username)
			}

			// Create profile
			profile := datamodels.UserProfileDataModel{
				ID:         uuid.NewV4(),
				UserID:     userID,
				Firstname:  su.firstname,
				Lastname:   su.lastname,
				Locale:     "en_US",
				Details:    datamodels.JSONB{"role": su.role},
				CreatedAt:  now,
				ModifiedAt: now,
			}

			if err := tx.Create(&profile).Error; err != nil {
				return errors.Wrapf(err, "Error creating profile for %s", su.username)
			}
		}

		return nil
	})
}
