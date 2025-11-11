package models

import (
	"time"

	"github.com/google/uuid"
	customTypes "github.com/phatnt199/go-infra/pkg/core/customtypes"
	"gorm.io/gorm"
)

// UserProfile represents user profile information
type UserProfile struct {
	ID        uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID         `gorm:"type:uuid;not null;unique;index:idx_user_profiles_user_id"`
	Firstname string            `gorm:"type:varchar(100)"`
	Lastname  string            `gorm:"type:varchar(100)"`
	Email     string            `gorm:"type:varchar(255)"`
	Birthday  *time.Time        `gorm:"type:date"`
	Locale    string            `gorm:"type:varchar(10);default:'en_US'"`
	Metadata  customTypes.JSONB `gorm:"type:jsonb"`
	CreatedAt time.Time         `gorm:"type:timestamptz;not null;default:current_timestamp"`
	UpdatedAt time.Time         `gorm:"type:timestamptz;not null;default:current_timestamp"`
	DeletedAt gorm.DeletedAt    `gorm:"index;type:timestamptz"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name
func (UserProfile) TableName() string {
	return "user_profiles"
}

// BeforeCreate sets UUID before creating
func (up *UserProfile) BeforeCreate(tx *gorm.DB) error {
	if up.ID == uuid.Nil {
		up.ID = uuid.New()
	}
	if up.Metadata == nil {
		up.Metadata = make(customTypes.JSONB)
	}
	return nil
}
