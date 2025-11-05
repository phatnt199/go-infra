package responses

// MessageResponse DTO for simple message responses
type MessageResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
