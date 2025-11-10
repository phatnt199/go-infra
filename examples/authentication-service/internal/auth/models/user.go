package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string    `gorm:"uniqueIndex;not null;size:100"`
	Email     string    `gorm:"uniqueIndex;size:255"`
	Password  string    `gorm:"not null;size:255"`
	Firstname string    `gorm:"size:100"`
	Lastname  string    `gorm:"size:100"`
	Birthday  *time.Time
	Locale    string                 `gorm:"size:10;default:'en_US'"`
	Status    string                 `gorm:"size:20;default:'active';index"`
	UserType  string                 `gorm:"size:20;default:'user'"`
	Metadata  map[string]interface{} `gorm:"type:jsonb"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}

// BeforeCreate hook to generate UUID
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.NewV4()
	}
	return nil
}
