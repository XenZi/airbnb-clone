package utils

import (
	"fmt"
	"reservation-service/domain"
	"time"
)

const (
	StartDate         = "Start date must be before end date"
	EndDate           = "End date must be after start date"
	Username          = "Username field can't be empty"
	AccommodationName = "Accommodation name can't be empty"
)

var errorMessages = map[string]string{
	"StartDate":         StartDate,
	"EndDate":           EndDate,
	"Username":          Username,
	"AccommodationName": AccommodationName,
}

type Validator struct {
	Errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

type ValidationRule func(value string) bool

func (v *Validator) ValidateField(fieldName string, value string, rules ...ValidationRule) {
	for _, rule := range rules {
		if !rule(value) {
			if errorMessage, exists := errorMessages[fieldName]; exists {
				v.Errors[fieldName] = errorMessage
			} else {
				v.Errors[fieldName] = "Validation failed for field " + fieldName
			}
			return
		}
	}
	delete(v.Errors, fieldName)
}

func (v *Validator) GetErrors() map[string]string {
	return v.Errors
}

func DateNotAfter(date1, date2 string) ValidationRule {
	return func(value string) bool {
		date, err := time.Parse("2006-01-02", value)
		if err != nil {
			return false
		}

		endDate, err := time.Parse("2006-01-02", date2)
		if err != nil {
			return false
		}

		return !date.After(endDate)
	}
}

func DateNotBefore(date1, date2 string) ValidationRule {
	return func(value string) bool {
		date, err := time.Parse("2006-01-02", value)
		if err != nil {
			return false
		}

		startDate, err := time.Parse("2006-01-02", date1)
		if err != nil {
			return false
		}

		return !date.Before(startDate)
	}
}
func MinLength(minLength int) ValidationRule {
	return func(value string) bool {
		return len(value) >= minLength
	}
}
func (v *Validator) ValidateReservation(reservation *domain.Reservation) {
	v.ValidateField("username", reservation.Username, MinLength(2))
	v.ValidateField("accommodationName", reservation.AccommodationName, MinLength(2))
	v.ValidateField("startDate", reservation.StartDate, DateNotAfter(reservation.StartDate, reservation.EndDate))
	v.ValidateField("endDate", reservation.EndDate, DateNotBefore(reservation.EndDate, reservation.StartDate))
	foundErrors := v.GetErrors()
	if len(foundErrors) > 0 {
		for field, message := range foundErrors {
			fmt.Printf("%s: %s\n", field, message)
		}
	}
}
