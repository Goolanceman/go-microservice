package models

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
} 