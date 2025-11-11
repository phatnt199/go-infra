package models

import "time"

// CustomAuthResponse extends the standard auth response
type CustomAuthResponse struct {
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Firstname string    `json:"firstname,omitempty"`
	Lastname  string    `json:"lastname,omitempty"`
	Status    string    `json:"status"`
	UserType  string    `json:"userType"`
	CreatedAt time.Time `json:"createdAt"`

	AccessToken  *TokenInfo `json:"token"`
	RefreshToken *TokenInfo `json:"refreshToken,omitempty"`

	// Custom fields
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	LoginCount  int        `json:"loginCount"`
	RequiresMFA bool       `json:"requiresMfa"`
	SessionID   string     `json:"sessionId"`
}

// TokenInfo represents token information
type TokenInfo struct {
	Value     string    `json:"value"`
	Scheme    string    `json:"scheme"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// AuditLogResponse represents an audit log response
type AuditLogResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Username  string    `json:"username"`
	Action    string    `json:"action"`
	Success   bool      `json:"success"`
	IPAddress string    `json:"ipAddress"`
	Details   string    `json:"details,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// PaginatedAuditLogsResponse represents paginated audit logs
type PaginatedAuditLogsResponse struct {
	Data       []AuditLogResponse `json:"data"`
	TotalCount int64              `json:"totalCount"`
	Page       int                `json:"page"`
	PageSize   int                `json:"pageSize"`
}
