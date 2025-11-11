package services

import (
	"github.com/phatnt199/go-infra/pkg/crypto"
	"golang.org/x/crypto/bcrypt"
)

// IPasswordHasher defines password hashing operations
type IPasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(password, hash string) (bool, error)
}

// BcryptHasher implements IPasswordHasher using pkg/crypto.Hasher
// This is a thin wrapper that provides backward compatibility while
// leveraging the more comprehensive crypto package implementation.
type BcryptHasher struct {
	hasher    *crypto.Hasher
	algorithm crypto.HashAlgorithm
}

// NewBcryptHasher creates a new bcrypt hasher
// The cost parameter is used to configure bcrypt hashing strength.
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}

	config := &crypto.HashConfig{
		BcryptCost: cost,
	}

	return &BcryptHasher{
		hasher:    crypto.NewHasher(config),
		algorithm: crypto.AlgorithmBcrypt,
	}
}

// NewArgon2Hasher creates a new Argon2 hasher
// This provides a more secure alternative to bcrypt for new implementations.
func NewArgon2Hasher() *BcryptHasher {
	return &BcryptHasher{
		hasher:    crypto.NewHasher(crypto.DefaultHashConfig()),
		algorithm: crypto.AlgorithmArgon2,
	}
}

// HashPassword hashes a password using the configured algorithm
func (h *BcryptHasher) HashPassword(password string) (string, error) {
	return h.hasher.HashPassword(password, h.algorithm)
}

// ComparePassword compares a password with a hash
// This automatically detects the algorithm used in the hash
func (h *BcryptHasher) ComparePassword(password, hash string) (bool, error) {
	return h.hasher.ComparePassword(password, hash)
}
