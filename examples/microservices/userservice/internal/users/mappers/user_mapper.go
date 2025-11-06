package mappers

import (
	"time"

	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/contracts"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/data/datamodels"
	dtosv1 "github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/dtos/v1"
	"github.com/phatnt199/go-infra/examples/microservices/userservice/internal/users/models"
	uuid "github.com/satori/go.uuid"
)

// ToUserResponse converts User model to UserResponse DTO
func ToUserResponse(user *models.User) *dtosv1.UserResponse {
	if user == nil {
		return nil
	}

	resp := &dtosv1.UserResponse{
		ID:         user.ID,
		Status:     user.Status,
		UserType:   string(user.UserType),
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
		ModifiedAt: user.ModifiedAt.Format(time.RFC3339),
	}

	if user.ActivatedAt != nil {
		resp.ActivatedAt = user.ActivatedAt.Format(time.RFC3339)
	}
	if user.LastLoginAt != nil {
		resp.LastLoginAt = user.LastLoginAt.Format(time.RFC3339)
	}
	if user.ParentID != nil {
		resp.ParentID = user.ParentID.String()
	}
	if user.ValidFrom != nil {
		resp.ValidFrom = user.ValidFrom.Format(time.RFC3339)
	}
	if user.ValidTo != nil {
		resp.ValidTo = user.ValidTo.Format(time.RFC3339)
	}

	return resp
}

// ToUserIdentifierResponse converts UserIdentifier model to UserIdentifierResponse DTO
func ToUserIdentifierResponse(identifier *models.UserIdentifier) *dtosv1.UserIdentifierResponse {
	if identifier == nil {
		return nil
	}

	return &dtosv1.UserIdentifierResponse{
		ID:         identifier.ID,
		UserID:     identifier.UserID,
		Scheme:     identifier.Scheme,
		Identifier: identifier.Identifier,
		Verified:   identifier.Verified,
		Details:    identifier.Details,
		CreatedAt:  identifier.CreatedAt.Format(time.RFC3339),
		ModifiedAt: identifier.ModifiedAt.Format(time.RFC3339),
	}
}

// ToUserProfileResponse converts UserProfile model to UserProfileResponse DTO
func ToUserProfileResponse(profile *models.UserProfile) *dtosv1.UserProfileResponse {
	if profile == nil {
		return nil
	}

	resp := &dtosv1.UserProfileResponse{
		ID:         profile.ID,
		UserID:     profile.UserID,
		Firstname:  profile.Firstname,
		Lastname:   profile.Lastname,
		Locale:     profile.Locale,
		Details:    profile.Details,
		CreatedAt:  profile.CreatedAt.Format(time.RFC3339),
		ModifiedAt: profile.ModifiedAt.Format(time.RFC3339),
	}

	if profile.Birthday != nil {
		resp.Birthday = profile.Birthday.Format("2006-01-02")
	}

	return resp
}

// ToUserFullDetailsResponse converts UserFullDetails to UserFullDetailsResponse DTO
func ToUserFullDetailsResponse(details *contracts.UserFullDetails) *dtosv1.UserFullDetailsResponse {
	if details == nil {
		return nil
	}

	return &dtosv1.UserFullDetailsResponse{
		User:       ToUserResponse(details.User),
		Identifier: ToUserIdentifierResponse(details.Identifier),
		Profile:    ToUserProfileResponse(details.Profile),
	}
}

// ToAuthResponse converts user details to AuthResponse DTO
func ToAuthResponse(user *models.User, identifier *models.UserIdentifier, profile *models.UserProfile, token string) *dtosv1.AuthResponse {
	resp := &dtosv1.AuthResponse{
		UserID:    user.ID,
		Status:    user.Status,
		UserType:  string(user.UserType),
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		Token:     token,
	}

	if identifier != nil {
		resp.Username = identifier.Identifier
	}

	if profile != nil {
		resp.Firstname = profile.Firstname
		resp.Lastname = profile.Lastname
	}

	return resp
}

// ParseISODate parses ISO date string (YYYY-MM-DD) to time.Time
func ParseISODate(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// ParseUUID parses string to UUID
func ParseUUID(uuidStr string) (*uuid.UUID, error) {
	if uuidStr == "" {
		return nil, nil
	}

	id, err := uuid.FromString(uuidStr)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// ToUser converts UserDataModel to User domain model
func ToUser(dataModel *datamodels.UserDataModel) *models.User {
	if dataModel == nil {
		return nil
	}

	user := &models.User{
		Status:      dataModel.Status.String(),   // Convert enum to string
		UserType:    dataModel.UserType.String(), // Convert enum to string
		ActivatedAt: dataModel.ActivatedAt,
		LastLoginAt: dataModel.LastLoginAt,
		ParentID:    dataModel.ParentID,
		ValidFrom:   dataModel.ValidFrom,
		ValidTo:     dataModel.ValidTo,
	}

	// Set base entity fields
	user.ID = dataModel.ID
	user.CreatedAt = dataModel.CreatedAt
	user.ModifiedAt = dataModel.ModifiedAt

	return user
}

// ToUserDataModel converts User domain model to UserDataModel
func ToUserDataModel(user *models.User) *datamodels.UserDataModel {
	if user == nil {
		return nil
	}

	dataModel := &datamodels.UserDataModel{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		ModifiedAt:  user.ModifiedAt,
		ActivatedAt: user.ActivatedAt,
		LastLoginAt: user.LastLoginAt,
		ParentID:    user.ParentID,
		ValidFrom:   user.ValidFrom,
		ValidTo:     user.ValidTo,
	}

	// Convert domain types to data model enums
	// dataModel.Status.FromInt(int(user.Status))
	dataModel.UserType = datamodels.UserTypeEnum(user.UserType)

	return dataModel
}
