package models

import (
	"time"

	"github.com/google/uuid"
	customTypes "github.com/phatnt199/go-infra/pkg/core/customtypes"
	"gorm.io/gorm"
)

// CredentialScheme represents different credential types
type CredentialScheme string

const (
	CredentialSchemeBasic  CredentialScheme = "basic"  // Password-based
	CredentialSchemeOAuth  CredentialScheme = "oauth"  // OAuth tokens
	CredentialSchemeAPIKey CredentialScheme = "apikey" // API keys
)

// UserCredential represents user authentication credentials
type UserCredential struct {
	ID         uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID         `gorm:"type:uuid;not null;index:idx_user_credentials_user_id_scheme"`
	Scheme     string            `gorm:"type:varchar(50);not null;index:idx_user_credentials_user_id_scheme"`
	Credential string            `gorm:"type:text;not null"` // Hashed password or token
	ExpiresAt  *time.Time        `gorm:"type:timestamptz"`
	Details    customTypes.JSONB `gorm:"type:jsonb"`
	CreatedAt  time.Time         `gorm:"type:timestamptz;not null;default:current_timestamp"`
	UpdatedAt  time.Time         `gorm:"type:timestamptz;not null;default:current_timestamp"`
	DeletedAt  gorm.DeletedAt    `gorm:"index;type:timestamptz"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name
func (UserCredential) TableName() string {
	return "user_credentials"
}

// BeforeCreate sets UUID before creating
func (uc *UserCredential) BeforeCreate(tx *gorm.DB) error {
	if uc.ID == uuid.Nil {
		uc.ID = uuid.New()
	}
	if uc.Details == nil {
		uc.Details = make(customTypes.JSONB)
	}
	return nil
}
