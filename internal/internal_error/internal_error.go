package internal_error

type InternalError struct {
	Message string
	Err     string
}

func NewNotFoundError(message string, err string) *InternalError {
	return &InternalError{
		Message: message,
		Err:     "not_found",
	}
}

func NewInternalServerError(message string, err string) *InternalError {
	return &InternalError{
		Message: message,
		Err:     "internal_server_error",
	}
}
