package users

import (
	"emperror.dev/errors"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/data/datamodels"
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

func seedDataManually(db *gorm.DB) error {
	var count int64

	db.Model(&datamodels.UserDataModel{}).Count(&count)
	if count > 0 {
		return nil
	}

	users := []datamodels.UserDataModel{
		{
			ID:    uuid.NewV4(),
			Name:  "John Doe",
			Email: "john.doe@example.com",
		},
		{
			ID:    uuid.NewV4(),
			Name:  "Jane Smith",
			Email: "jane.smith@example.com",
		},
	}

	err := db.CreateInBatches(users, len(users)).Error
	if err != nil {
		return errors.Wrap(err, "Error seeding users data")
	}

	return nil
}
