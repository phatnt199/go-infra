package services

import (
	"context"
	"fmt"
	"time"

	"github.com/phatnt199/go-infra/examples/authentication-service/internal/custom-auth/models"
	"github.com/phatnt199/go-infra/pkg/infra/postgres/gorm/contracts"
	"github.com/phatnt199/go-infra/pkg/logger"
)

// AuditLogger handles authentication audit logging
type AuditLogger struct {
	dbContext contracts.GormDBContext
	logger    logger.Logger
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(dbContext contracts.GormDBContext, logger logger.Logger) *AuditLogger {
	return &AuditLogger{
		dbContext: dbContext,
		logger:    logger,
	}
}

// LogAuthEvent logs an authentication event
func (a *AuditLogger) LogAuthEvent(ctx context.Context, event models.AuditLog) error {
	if err := a.dbContext.DB().WithContext(ctx).Create(&event).Error; err != nil {
		a.logger.Errorf("Failed to create audit log: %v", err)
		return err
	}
	return nil
}

// GetUserAuditLogs retrieves audit logs for a specific user
func (a *AuditLogger) GetUserAuditLogs(ctx context.Context, userID string, limit, offset int) ([]models.AuditLog, int64, error) {
	var logs []models.AuditLog
	var total int64

	// Count total
	if err := a.dbContext.DB().WithContext(ctx).Model(&models.AuditLog{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated logs
	if err := a.dbContext.DB().WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// LoginAttemptManager handles login attempt tracking
type LoginAttemptManager struct {
	dbContext    contracts.GormDBContext
	logger       logger.Logger
	maxAttempts  int
	lockDuration time.Duration
}

// NewLoginAttemptManager creates a new login attempt manager
func NewLoginAttemptManager(
	dbContext contracts.GormDBContext,
	logger logger.Logger,
	maxAttempts int,
	lockDuration time.Duration,
) *LoginAttemptManager {
	return &LoginAttemptManager{
		dbContext:    dbContext,
		logger:       logger,
		maxAttempts:  maxAttempts,
		lockDuration: lockDuration,
	}
}

// RecordLoginAttempt records a login attempt
func (m *LoginAttemptManager) RecordLoginAttempt(ctx context.Context, identifier string, success bool) error {
	var attempt models.LoginAttempt

	// Find or create login attempt record
	err := m.dbContext.DB().WithContext(ctx).
		Where("identifier = ?", identifier).
		First(&attempt).Error

	if err != nil {
		// Create new record
		attempt = models.LoginAttempt{
			Identifier:  identifier,
			Attempts:    0,
			LastAttempt: time.Now(),
		}
	}

	if success {
		// Reset attempts on successful login
		attempt.Attempts = 0
		attempt.LockedUntil = nil
	} else {
		// Increment failed attempts
		attempt.Attempts++
		attempt.LastAttempt = time.Now()

		// Lock account if max attempts reached
		if attempt.Attempts >= m.maxAttempts {
			lockUntil := time.Now().Add(m.lockDuration)
			attempt.LockedUntil = &lockUntil
		}
	}

	// Save or update
	if err := m.dbContext.DB().WithContext(ctx).Save(&attempt).Error; err != nil {
		m.logger.Errorf("Failed to save login attempt: %v", err)
		return err
	}

	return nil
}

// IsAccountLocked checks if an account is locked
func (m *LoginAttemptManager) IsAccountLocked(ctx context.Context, identifier string) (bool, error) {
	var attempt models.LoginAttempt

	err := m.dbContext.DB().WithContext(ctx).
		Where("identifier = ?", identifier).
		First(&attempt).Error

	if err != nil {
		// No record found means not locked
		return false, nil
	}

	// Check if locked
	if attempt.LockedUntil != nil && time.Now().Before(*attempt.LockedUntil) {
		return true, nil
	}

	// If lock has expired, reset it
	if attempt.LockedUntil != nil && time.Now().After(*attempt.LockedUntil) {
		attempt.LockedUntil = nil
		attempt.Attempts = 0
		m.dbContext.DB().WithContext(ctx).Save(&attempt)
	}

	return false, nil
}

// GetRemainingAttempts returns the number of remaining login attempts
func (m *LoginAttemptManager) GetRemainingAttempts(ctx context.Context, identifier string) (int, error) {
	var attempt models.LoginAttempt

	err := m.dbContext.DB().WithContext(ctx).
		Where("identifier = ?", identifier).
		First(&attempt).Error

	if err != nil {
		// No record means all attempts available
		return m.maxAttempts, nil
	}

	remaining := m.maxAttempts - attempt.Attempts
	if remaining < 0 {
		return 0, nil
	}

	return remaining, nil
}

// GetLockTimeRemaining returns the remaining lock time
func (m *LoginAttemptManager) GetLockTimeRemaining(ctx context.Context, identifier string) (time.Duration, error) {
	var attempt models.LoginAttempt

	err := m.dbContext.DB().WithContext(ctx).
		Where("identifier = ?", identifier).
		First(&attempt).Error

	if err != nil || attempt.LockedUntil == nil {
		return 0, nil
	}

	if time.Now().After(*attempt.LockedUntil) {
		return 0, nil
	}

	return time.Until(*attempt.LockedUntil), nil
}

// UnlockAccount manually unlocks an account (admin function)
func (m *LoginAttemptManager) UnlockAccount(ctx context.Context, identifier string) error {
	return m.dbContext.DB().WithContext(ctx).
		Model(&models.LoginAttempt{}).
		Where("identifier = ?", identifier).
		Updates(map[string]interface{}{
			"attempts":     0,
			"locked_until": nil,
		}).Error
}

// CleanupOldAttempts removes old login attempt records
func (m *LoginAttemptManager) CleanupOldAttempts(ctx context.Context, olderThan time.Duration) error {
	cutoffTime := time.Now().Add(-olderThan)

	result := m.dbContext.DB().WithContext(ctx).
		Where("last_attempt < ? AND (locked_until IS NULL OR locked_until < ?)", cutoffTime, time.Now()).
		Delete(&models.LoginAttempt{})

	if result.Error != nil {
		return fmt.Errorf("failed to cleanup old attempts: %w", result.Error)
	}

	m.logger.Infof("Cleaned up %d old login attempt records", result.RowsAffected)
	return nil
}
