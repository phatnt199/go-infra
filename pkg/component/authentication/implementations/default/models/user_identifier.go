package models

import (
	"time"

	"github.com/google/uuid"
	customTypes "github.com/phatnt199/go-infra/pkg/core/customtypes"
	"gorm.io/gorm"
)

// IdentifierScheme represents different identifier types
type IdentifierScheme string

const (
	IdentifierSchemeUsername IdentifierScheme = "username"
	IdentifierSchemeEmail    IdentifierScheme = "email"
	IdentifierSchemePhone    IdentifierScheme = "phone"
)

// UserIdentifier represents how users can be identified (username, email, phone, etc.)
type UserIdentifier struct {
	ID         uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID         `gorm:"type:uuid;not null;index:idx_user_identifiers_user_id"`
	Scheme     string            `gorm:"type:varchar(50);not null;index:idx_user_identifiers_scheme_identifier"`
	Identifier string            `gorm:"type:varchar(255);not null;index:idx_user_identifiers_scheme_identifier"`
	Verified   bool              `gorm:"type:boolean;not null;default:false"`
	VerifiedAt *time.Time        `gorm:"type:timestamptz"`
	Details    customTypes.JSONB `gorm:"type:jsonb"`
	CreatedAt  time.Time         `gorm:"type:timestamptz;not null;default:current_timestamp"`
	UpdatedAt  time.Time         `gorm:"type:timestamptz;not null;default:current_timestamp"`
	DeletedAt  gorm.DeletedAt    `gorm:"index;type:timestamptz"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName overrides the table name
func (UserIdentifier) TableName() string {
	return "user_identifiers"
}

// BeforeCreate sets UUID before creating
func (ui *UserIdentifier) BeforeCreate(tx *gorm.DB) error {
	if ui.ID == uuid.Nil {
		ui.ID = uuid.New()
	}
	if ui.Details == nil {
		ui.Details = make(customTypes.JSONB)
	}
	return nil
}
