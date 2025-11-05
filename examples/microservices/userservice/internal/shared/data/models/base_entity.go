package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// BaseEntity contains common fields for all entities
type BaseEntity struct {
	ID         uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt  time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"createdAt"`
	ModifiedAt time.Time      `gorm:"column:modified_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"modifiedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}

// BeforeCreate sets the ID before creating a new record
func (b *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.NewV4()
	}
	now := time.Now()
	if b.CreatedAt.IsZero() {
		b.CreatedAt = now
	}
	if b.ModifiedAt.IsZero() {
		b.ModifiedAt = now
	}
	return nil
}

// BeforeUpdate sets the ModifiedAt before updating a record
func (b *BaseEntity) BeforeUpdate(tx *gorm.DB) error {
	b.ModifiedAt = time.Now()
	return nil
}
