package requests

// SignUpRequest DTO for user registration
type SignUpRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Password  string `json:"password" validate:"required,min=8"`
	Firstname string `json:"firstname" validate:"max=100"`
	Lastname  string `json:"lastname" validate:"max=100"`
	Locale    string `json:"locale" validate:"omitempty,len=5"` // e.g., en_US, vi_VN
}
