package errors

import (
	"fmt"
)

type ReservationError struct {
	Status  int    `json: "status"`
	Message string `json: "message"`
}

func NewReservationError(status int, message string) *ReservationError {
	return &ReservationError{
		Status:  status,
		Message: message,
	}
}

func (e *ReservationError) Error() string {
	return fmt.Sprintf("ReservationError - Status: %d, Message: %s", e.Status, e.Message)
}
