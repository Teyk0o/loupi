package models

// ErrorResponse is the standard error format returned by the API.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// MessageResponse is a simple success message.
type MessageResponse struct {
	Message string `json:"message"`
}
