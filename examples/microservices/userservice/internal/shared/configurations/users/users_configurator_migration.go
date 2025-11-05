package users

import (
	"context"

	"github.com/phatnt199/go-infra/pkg/migration/contracts"
)

func (ic *UsersServiceConfigurator) migrateUsers(
	runner contracts.PostgresMigrationRunner,
) error {
	return migrateGoose(runner)
}

func migrateGoose(runner contracts.PostgresMigrationRunner) error {
	err := runner.Up(context.Background(), 0)

	return err
}

// func migrateGorm(db *gorm.DB) error {
// 	err := db.AutoMigrate(&models.User{})

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
