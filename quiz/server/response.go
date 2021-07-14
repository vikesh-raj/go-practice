package server

// BasicResponse contains status and error message
type BasicResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message,omitempty"`
}
