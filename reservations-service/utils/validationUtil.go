package utils

import (
	"fmt"
	"reservation-service/domain"
	"strings"
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
func DateNotInPast() ValidationRule {
	return func(value string) bool {
		today := time.Now().UTC()
		date, err := time.Parse("2006-01-02", value)
		if err != nil {
			return false
		}
		return !date.Before(today)
	}
}

func MinLength(minLength int) ValidationRule {
	return func(value string) bool {
		return len(value) >= minLength
	}
}

func MaxLength(maxLength int) ValidationRule {
	return func(value string) bool {
		return len(value) <= maxLength
	}
}
func WordBan(value string) bool {
	normalizedValue := strings.ReplaceAll(strings.ToUpper(value), " ", "")
	bannedWords := []string{"SELECT", "UPDATE", "DELETE", "FROM", "WHERE", "<", ">"}
	for _, word := range bannedWords {
		if strings.Contains(normalizedValue, word) {
			return false
		}
	}

	return true
}

func (v *Validator) ValidateReservation(reservation *domain.Reservation) {
	/*	v.ValidateField("username", reservation.Username, MinLength(2), MaxLength(15), WordBan)
		v.ValidateField("accommodationName", reservation.AccommodationName, MinLength(2), MaxLength(15), WordBan)
		v.ValidateField("startDate", reservation.StartDate, DateNotAfter(reservation.StartDate, reservation.EndDate), DateNotInPast())
		v.ValidateField("endDate", reservation.EndDate, DateNotBefore(reservation.EndDate, reservation.StartDate), DateNotInPast()) */
	foundErrors := v.GetErrors()
	if len(foundErrors) > 0 {
		for field, message := range foundErrors {
			fmt.Printf("%s: %s\n", field, message)
		}
	}
}
