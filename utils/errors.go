package utils

type ErrorType string

const (
	UserError   ErrorType = "USER_ERROR"
	SystemError ErrorType = "SYSTEM_ERROR"
)

type customError struct {
	Type    ErrorType
	Message string
}

type CustomError interface {
	Error() string
	IsUserError() bool
	IsSystemError() bool
}

func (e *customError) Error() string {
	return e.Message
}

// NewCustomUserError creates a new user error with a message
func NewCustomUserError(message string) *customError {
	return &customError{
		Type:    UserError,
		Message: message,
	}
}

// NewCustomSystemError creates a new system error with a message
func NewCustomSystemError(message string) *customError {
	return &customError{
		Type:    SystemError,
		Message: message,
	}
}

// IsUserError method checks if the error is of type USER_ERROR
func (e *customError) IsUserError() bool {
	return e.Type == UserError
}

// IsSystemError method checks if the error is of type SYSTEM_ERROR
func (e *customError) IsSystemError() bool {
	return e.Type == SystemError
}
