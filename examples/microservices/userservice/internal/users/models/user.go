package models

import (
	"time"

	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/shared/data/models"
	uuid "github.com/satori/go.uuid"
)

// User represents a user in the system
type User struct {
	models.BaseEntity
	Status      string     `json:"status"`
	UserType    string     `json:"userType"`
	ActivatedAt *time.Time `json:"activatedAt,omitempty"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	ParentID    *uuid.UUID `json:"parentId,omitempty"`
	ValidFrom   *time.Time `json:"validFrom,omitempty"`
	ValidTo     *time.Time `json:"validTo,omitempty"`
}

// UserIdentifier represents user identifier (username, phone, etc.)
type UserIdentifier struct {
	models.BaseEntity
	UserID     uuid.UUID              `json:"userId"`
	Scheme     string                 `json:"scheme"` // "username", "phoneNumber"
	Identifier string                 `json:"identifier"`
	Verified   bool                   `json:"verified"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// IdentifierScheme constants
const (
	IdentifierSchemeUsername    = "username"
	IdentifierSchemePhoneNumber = "phoneNumber"
)

// UserCredential represents user credentials (password, etc.)
type UserCredential struct {
	models.BaseEntity
	UserID     uuid.UUID              `json:"userId"`
	Scheme     string                 `json:"scheme"` // "basic", "master-basic"
	Credential string                 `json:"credential"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

// CredentialScheme constants
const (
	CredentialSchemeBasic       = "basic"
	CredentialSchemeMasterBasic = "master-basic"
)

// UserProfile represents user profile information
type UserProfile struct {
	models.BaseEntity
	UserID    uuid.UUID              `json:"userId"`
	Firstname string                 `json:"firstname"`
	Lastname  string                 `json:"lastname"`
	Birthday  *time.Time             `json:"birthday,omitempty"`
	Locale    string                 `json:"locale"`
	Details   map[string]interface{} `json:"details,omitempty"`
}
