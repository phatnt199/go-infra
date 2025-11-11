package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents authentication component configuration
// This follows the go-infra config pattern (loaded from environment variables)
type Config struct {
	JWT      JWTConfig      `json:"jwt"`
	Password PasswordConfig `json:"password"`
	Schemes  SchemesConfig  `json:"schemes"`
}

// JWTConfig contains JWT token settings
type JWTConfig struct {
	Secret        string        `json:"-"` // Never log secrets
	Issuer        string        `json:"issuer"`
	Audience      string        `json:"audience"`
	Algorithm     string        `json:"algorithm"` // HS256, RS256
	AccessExpiry  time.Duration `json:"access_expiry"`
	RefreshExpiry time.Duration `json:"refresh_expiry"`
}

// PasswordConfig contains password hashing and policy settings
type PasswordConfig struct {
	HashAlgorithm  string `json:"hash_algorithm"` // bcrypt, argon2
	BcryptCost     int    `json:"bcrypt_cost"`
	MinLength      int    `json:"min_length"`
	RequireUpper   bool   `json:"require_upper"`
	RequireLower   bool   `json:"require_lower"`
	RequireNumber  bool   `json:"require_number"`
	RequireSpecial bool   `json:"require_special"`
}

// SchemesConfig contains supported authentication schemes
type SchemesConfig struct {
	IdentifierSchemes []string `json:"identifier_schemes"` // username, email, phone
	CredentialSchemes []string `json:"credential_schemes"` // basic, token
}

// LoadFromEnv loads authentication configuration from environment variables
// Following the go-infra config pattern
func LoadFromEnv() *Config {
	return &Config{
		JWT: JWTConfig{
			Secret:        getEnv("AUTH_JWT_SECRET", ""),
			Issuer:        getEnv("AUTH_JWT_ISSUER", "go-infra-auth"),
			Audience:      getEnv("AUTH_JWT_AUDIENCE", "go-infra-api"),
			Algorithm:     getEnv("AUTH_JWT_ALGORITHM", "HS256"),
			AccessExpiry:  getEnvAsDuration("AUTH_JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshExpiry: getEnvAsDuration("AUTH_JWT_REFRESH_EXPIRY", 7*24*time.Hour),
		},
		Password: PasswordConfig{
			HashAlgorithm:  getEnv("AUTH_PASSWORD_HASH_ALGORITHM", "bcrypt"),
			BcryptCost:     getEnvAsInt("AUTH_PASSWORD_BCRYPT_COST", 12),
			MinLength:      getEnvAsInt("AUTH_PASSWORD_MIN_LENGTH", 8),
			RequireUpper:   getEnvAsBool("AUTH_PASSWORD_REQUIRE_UPPER", false),
			RequireLower:   getEnvAsBool("AUTH_PASSWORD_REQUIRE_LOWER", false),
			RequireNumber:  getEnvAsBool("AUTH_PASSWORD_REQUIRE_NUMBER", false),
			RequireSpecial: getEnvAsBool("AUTH_PASSWORD_REQUIRE_SPECIAL", false),
		},
		Schemes: SchemesConfig{
			IdentifierSchemes: getEnvAsSlice("AUTH_SUPPORTED_IDENTIFIER_SCHEMES", []string{"username", "email", "phone"}),
			CredentialSchemes: getEnvAsSlice("AUTH_SUPPORTED_CREDENTIAL_SCHEMES", []string{"basic"}),
		},
	}
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		JWT: JWTConfig{
			Secret:        "change-this-secret",
			Issuer:        "go-infra-auth",
			Audience:      "go-infra-api",
			Algorithm:     "HS256",
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 7 * 24 * time.Hour,
		},
		Password: PasswordConfig{
			HashAlgorithm:  "bcrypt",
			BcryptCost:     12,
			MinLength:      8,
			RequireUpper:   false,
			RequireLower:   false,
			RequireNumber:  false,
			RequireSpecial: false,
		},
		Schemes: SchemesConfig{
			IdentifierSchemes: []string{"username", "email", "phone"},
			CredentialSchemes: []string{"basic"},
		},
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Add validation logic here if needed
	return nil
}

// Helper functions for environment variable parsing

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}
