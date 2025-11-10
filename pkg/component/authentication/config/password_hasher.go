package config

import (
	"fmt"

	"github.com/phatnt199/go-infra/pkg/component/authentication/contracts"
	"golang.org/x/crypto/bcrypt"
)

// passwordHasher implements IPasswordHasher
type passwordHasher struct {
	algorithm string
	cost      int
}

// NewPasswordHasher creates a new password hasher
func NewPasswordHasher(algorithm string) contracts.IPasswordHasher {
	cost := bcrypt.DefaultCost
	if algorithm == "bcrypt-strong" {
		cost = 12
	}

	return &passwordHasher{
		algorithm: algorithm,
		cost:      cost,
	}
}

// HashPassword hashes a password using bcrypt
func (h *passwordHasher) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// ComparePassword compares a password with a hash
func (h *passwordHasher) ComparePassword(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, fmt.Errorf("failed to compare password: %w", err)
	}
	return true, nil
}
