package validators

import (
	"fmt"
	"strings"
	"unicode"
)

// PasswordValidator provides custom password validation
type PasswordValidator struct {
	minLength      int
	requireUpper   bool
	requireLower   bool
	requireNumber  bool
	requireSpecial bool
	minEntropy     float64
}

// NewPasswordValidator creates a new password validator
func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{
		minLength:      8,
		requireUpper:   true,
		requireLower:   true,
		requireNumber:  true,
		requireSpecial: true,
		minEntropy:     50.0,
	}
}

// ValidatePassword validates a password against custom rules
func (v *PasswordValidator) ValidatePassword(password string) error {
	if len(password) < v.minLength {
		return fmt.Errorf("password must be at least %d characters long", v.minLength)
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if v.requireUpper && !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	if v.requireLower && !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	if v.requireNumber && !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}

	if v.requireSpecial && !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	// Check for common weak passwords
	if v.isCommonPassword(password) {
		return fmt.Errorf("password is too common, please choose a stronger password")
	}

	// Calculate password entropy
	entropy := v.calculateEntropy(password)
	if entropy < v.minEntropy {
		return fmt.Errorf("password is too weak (entropy: %.2f, required: %.2f)", entropy, v.minEntropy)
	}

	return nil
}

// isCommonPassword checks if password is in common passwords list
func (v *PasswordValidator) isCommonPassword(password string) bool {
	commonPasswords := []string{
		"password", "password123", "12345678", "qwerty", "abc123",
		"monkey", "1234567", "letmein", "trustno1", "dragon",
		"baseball", "iloveyou", "master", "sunshine", "ashley",
		"bailey", "passw0rd", "shadow", "123123", "654321",
	}

	lowerPassword := strings.ToLower(password)
	for _, common := range commonPasswords {
		if lowerPassword == common {
			return true
		}
	}

	return false
}

// calculateEntropy calculates password entropy
func (v *PasswordValidator) calculateEntropy(password string) float64 {
	var charsetSize int

	if v.hasLowercase(password) {
		charsetSize += 26
	}
	if v.hasUppercase(password) {
		charsetSize += 26
	}
	if v.hasNumber(password) {
		charsetSize += 10
	}
	if v.hasSpecial(password) {
		charsetSize += 32
	}

	if charsetSize == 0 {
		return 0
	}

	// Entropy = length * log2(charset_size)
	return float64(len(password)) * log2(float64(charsetSize))
}

// Helper functions
func (v *PasswordValidator) hasLowercase(s string) bool {
	for _, char := range s {
		if unicode.IsLower(char) {
			return true
		}
	}
	return false
}

func (v *PasswordValidator) hasUppercase(s string) bool {
	for _, char := range s {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

func (v *PasswordValidator) hasNumber(s string) bool {
	for _, char := range s {
		if unicode.IsNumber(char) {
			return true
		}
	}
	return false
}

func (v *PasswordValidator) hasSpecial(s string) bool {
	for _, char := range s {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			return true
		}
	}
	return false
}

// log2 calculates log base 2
func log2(x float64) float64 {
	if x <= 0 {
		return 0
	}
	// log2(x) = ln(x) / ln(2)
	return 1.44269504089 * logN(x) // ln(2) ≈ 0.693, 1/ln(2) ≈ 1.443
}

// logN calculates natural logarithm using Taylor series
func logN(x float64) float64 {
	if x <= 0 {
		return 0
	}
	if x == 1 {
		return 0
	}

	// Use built-in approximation for simplicity
	// In production, you might want to use math.Log
	sum := 0.0
	term := (x - 1) / x
	power := term
	for i := 1; i <= 100; i++ {
		sum += power / float64(i)
		power *= term
	}
	return sum
}
