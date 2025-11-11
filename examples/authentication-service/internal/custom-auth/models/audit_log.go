package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// AuditLogAction represents the type of action being audited
type AuditLogAction string

const (
	AuditActionSignIn         AuditLogAction = "signin"
	AuditActionSignUp         AuditLogAction = "signup"
	AuditActionSignOut        AuditLogAction = "signout"
	AuditActionChangePassword AuditLogAction = "change_password"
	AuditActionUpdateProfile  AuditLogAction = "update_profile"
	AuditActionFailedLogin    AuditLogAction = "failed_login"
	AuditActionAccountLocked  AuditLogAction = "account_locked"
)

// AuditLog represents an authentication audit log entry
type AuditLog struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    string         `json:"userId" gorm:"type:varchar(255);index"`
	Username  string         `json:"username" gorm:"type:varchar(255);index"`
	Action    AuditLogAction `json:"action" gorm:"type:varchar(50);index"`
	Success   bool           `json:"success" gorm:"type:boolean;index"`
	IPAddress string         `json:"ipAddress" gorm:"type:varchar(45)"`
	UserAgent string         `json:"userAgent" gorm:"type:text"`
	Details   string         `json:"details" gorm:"type:text"`
	CreatedAt time.Time      `json:"createdAt" gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for GORM
func (AuditLog) TableName() string {
	return "audit_logs"
}

// LoginAttempt tracks login attempts for rate limiting
type LoginAttempt struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Identifier  string     `gorm:"type:varchar(255);index;not null"`
	Attempts    int        `gorm:"type:int;default:0"`
	LockedUntil *time.Time `gorm:"type:timestamptz"`
	LastAttempt time.Time  `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	CreatedAt   time.Time  `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `gorm:"type:timestamptz;default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for GORM
func (LoginAttempt) TableName() string {
	return "login_attempts"
}
