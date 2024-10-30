package utils

type ServerError struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

func NewServerError(errorMessage, reason string) *ServerError {
	return &ServerError{
		Error:  errorMessage,
		Reason: reason,
	}
}
