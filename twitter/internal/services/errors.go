package services

// Generic internal error used in the services for better error handling
type ServiceError struct {
	StatusCode int `json:"status_code"`
	Message string `json:"message"`
}

func (e ServiceError) Error() string {
	return e.Message
}

func (e ServiceError) Code() int {
	return e.StatusCode
}
