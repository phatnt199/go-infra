package services

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// IPasswordHasher defines password hashing operations
type IPasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(password, hash string) (bool, error)
}

// BcryptHasher implements IPasswordHasher using bcrypt
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new bcrypt hasher
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &BcryptHasher{
		cost: cost,
	}
}

// HashPassword hashes a password using bcrypt
func (h *BcryptHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// ComparePassword compares a password with a hash
func (h *BcryptHasher) ComparePassword(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, fmt.Errorf("failed to compare password: %w", err)
	}
	return true, nil
}
