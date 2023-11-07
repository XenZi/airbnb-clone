package errors

import "errors"

var (
	errUnauthorized    = errors.New("Unauthorized ")
	errDuplicateEntity = errors.New("Duplicated entity")
	errInternalServer  = errors.New("Internal server error")
)

type ErrorStruct struct {
	message string
	status  int
}

func NewError(message string, status int) *ErrorStruct {
	return &ErrorStruct{
		message: message,
		status:  status,
	}
}

func (e *ErrorStruct) GetErrorMessage() string {
	return e.message
}

func (e *ErrorStruct) GetErrorStatus() int {
	return e.status
}

func ErrInternalServerError() error {
	return errInternalServer
}

func ErrUnauthorized() error {
	return errUnauthorized
}

func ErrDuplicateEntity() error {
	return errDuplicateEntity
}
