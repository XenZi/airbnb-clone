package errors

import "errors"

var (
	errUnathorized error = errors.New("Unathorized")
)

func ErrUnathorized() error {
	return errUnathorized
}
