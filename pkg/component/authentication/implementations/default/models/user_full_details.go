package models

// UserFullDetails contains complete user information including relationships
type UserFullDetails struct {
	User       *User
	Identifier *UserIdentifier
	Profile    *UserProfile
	Credential *UserCredential
}
