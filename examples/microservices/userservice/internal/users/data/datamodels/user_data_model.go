package datamodels

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// JSONB is a custom type for JSONB fields
type JSONB map[string]interface{}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return json.Marshal(map[string]interface{}{})
	}
	return json.Marshal(j)
}

// UserTypeEnum represents user type as an enum
type UserTypeEnum string

const (
	UserTypeSystem UserTypeEnum = "SYSTEM"
	UserTypeLinked UserTypeEnum = "LINKED"
)

// Scan implements the sql.Scanner interface for UserTypeEnum
func (u *UserTypeEnum) Scan(value interface{}) error {
	if value == nil {
		*u = UserTypeSystem // default value
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into UserTypeEnum", value)
	}
	*u = UserTypeEnum(str)
	return nil
}

// Value implements the driver.Valuer interface for UserTypeEnum
func (u UserTypeEnum) Value() (driver.Value, error) {
	return string(u), nil
}

// String returns the string representation
func (u UserTypeEnum) String() string {
	return string(u)
}

// UserStatusEnum represents user status as an enum with string storage
type UserStatusEnum string

const (
	UserStatusActivated   UserStatusEnum = "100_ACTIVATE"
	UserStatusDeactivated UserStatusEnum = "101_DEACTIVATE"
)

// ToInt converts the enum to its integer value
func (u UserStatusEnum) ToInt() int {
	switch u {
	case UserStatusActivated:
		return 100
	case UserStatusDeactivated:
		return 101
	default:
		return 100 // default
	}
}

// FromInt creates a UserStatusEnum from an integer
func (u *UserStatusEnum) FromInt(value int) {
	switch value {
	case 100:
		*u = UserStatusActivated
	case 101:
		*u = UserStatusDeactivated
	default:
		*u = UserStatusActivated // default
	}
}

// Scan implements the sql.Scanner interface for UserStatusEnum
func (u *UserStatusEnum) Scan(value interface{}) error {
	if value == nil {
		*u = UserStatusActivated // default value
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into UserStatusEnum", value)
	}
	*u = UserStatusEnum(str)
	return nil
}

// Value implements the driver.Valuer interface for UserStatusEnum
func (u UserStatusEnum) Value() (driver.Value, error) {
	return string(u), nil
}

// String returns the string representation
func (u UserStatusEnum) String() string {
	return string(u)
}

// UserDataModel represents the users table
type UserDataModel struct {
	ID          uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	ModifiedAt  time.Time      `gorm:"column:modified_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Status      UserStatusEnum `gorm:"column:status;type:varchar(20);not null;default:'100_ACTIVATE'"`
	UserType    UserTypeEnum   `gorm:"column:user_type;type:varchar(10);not null;default:'SYSTEM'"`
	ActivatedAt *time.Time     `gorm:"column:activated_at;type:timestamp"`
	LastLoginAt *time.Time     `gorm:"column:last_login_at;type:timestamp"`
	ParentID    *uuid.UUID     `gorm:"column:parent_id;type:uuid"`
	ValidFrom   *time.Time     `gorm:"column:valid_from;type:timestamp"`
	ValidTo     *time.Time     `gorm:"column:valid_to;type:timestamp"`
}

func (u *UserDataModel) TableName() string {
	return "users"
}

// UserIdentifierDataModel represents the user_identifiers table
type UserIdentifierDataModel struct {
	ID         uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID      `gorm:"column:user_id;type:uuid;not null"`
	CreatedAt  time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	ModifiedAt time.Time      `gorm:"column:modified_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Scheme     string         `gorm:"column:scheme;type:varchar(50);not null;default:'username'"`
	Identifier string         `gorm:"column:identifier;type:varchar(255);not null"`
	Verified   bool           `gorm:"column:verified;type:boolean;not null;default:true"`
	Details    JSONB          `gorm:"column:details;type:jsonb;default:'{}'"`
}

func (u *UserIdentifierDataModel) TableName() string {
	return "user_identifiers"
}

// UserCredentialDataModel represents the user_credentials table
type UserCredentialDataModel struct {
	ID         uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID      `gorm:"column:user_id;type:uuid;not null"`
	CreatedAt  time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	ModifiedAt time.Time      `gorm:"column:modified_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Scheme     string         `gorm:"column:scheme;type:varchar(50);not null;default:'basic'"`
	Credential string         `gorm:"column:credential;type:varchar(255);not null"`
	Details    JSONB          `gorm:"column:details;type:jsonb;default:'{}'"`
}

func (u *UserCredentialDataModel) TableName() string {
	return "user_credentials"
}

// UserProfileDataModel represents the user_profiles table
type UserProfileDataModel struct {
	ID         uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID      `gorm:"column:user_id;type:uuid;not null"`
	CreatedAt  time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	ModifiedAt time.Time      `gorm:"column:modified_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Firstname  string         `gorm:"column:firstname;type:varchar(255)"`
	Lastname   string         `gorm:"column:lastname;type:varchar(255)"`
	Birthday   *time.Time     `gorm:"column:birthday;type:date"`
	Locale     string         `gorm:"column:locale;type:varchar(10);default:'en_US'"`
	Details    JSONB          `gorm:"column:details;type:jsonb;default:'{}'"`
}

func (u *UserProfileDataModel) TableName() string {
	return "user_profiles"
}
