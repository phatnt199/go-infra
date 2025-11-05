package datamodels

import (
	"github.com/goccy/go-json"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type UserDataModel struct {
	ID       uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Name     string    `gorm:"column:name"`
	Email    string    `gorm:"column:email"`
	Password string    `gorm:"column:password"`

	gorm.DeletedAt
}

func (u *UserDataModel) TableName() string {
	return "users"
}

func (u *UserDataModel) String() string {
	j, _ := json.Marshal(u)

	return string(j)
}
