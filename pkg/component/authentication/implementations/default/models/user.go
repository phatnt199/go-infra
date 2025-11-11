package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserStatus represents user account status
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBanned    UserStatus = "banned"
)

// User represents the core user entity
type User struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserType    string         `gorm:"type:varchar(50);not null;default:'user'"`
	Status      string         `gorm:"type:varchar(50);not null;default:'active'"`
	LastLoginAt *time.Time     `gorm:"type:timestamptz"`
	ActivatedAt *time.Time     `gorm:"type:timestamptz"`
	CreatedAt   time.Time      `gorm:"type:timestamptz;not null;default:current_timestamp"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz;not null;default:current_timestamp"`
	DeletedAt   gorm.DeletedAt `gorm:"index;type:timestamptz"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}

// BeforeCreate sets UUID before creating
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
