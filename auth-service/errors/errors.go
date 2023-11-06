package errors

import "errors"

var (
	errUnauthorized    = errors.New("Unauthorized ")
	errDuplicateEntity = errors.New("Duplicated entity")
	errInternalServer  = errors.New("Internal server error")
)

func ErrInternalServerError() error {
	return errInternalServer
}
func ErrUnauthorized() error {
	return errUnauthorized
}

func ErrDuplicateEntity() error {
	return errDuplicateEntity
}
