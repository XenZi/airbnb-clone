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

type ValidationRule func(value []domain.DateRangeWithPrice) bool

func (v *Validator) ValidateField(fieldName string, value []domain.DateRangeWithPrice, rules ...ValidationRule) {
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

/*
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
*/
func DateNotSame() ValidationRule {
	return func(existingDates []domain.DateRangeWithPrice) bool {
		validDates := make([][]time.Time, 0)
		allDates := make([][]time.Time, 0)
		for n, dateRangeWithPrice := range existingDates {
			datePair := make([]time.Time, 0)
			firstDate, err := getTimeFromString(dateRangeWithPrice.DateRange[0])
			if err != nil {
				return false
			}
			lastDate, erro := getTimeFromString(dateRangeWithPrice.DateRange[len(dateRangeWithPrice.DateRange)-1])
			if erro != nil {
				return false
			}
			datePair = append(datePair, firstDate)
			datePair = append(datePair, lastDate)
			if n == 0 {
				if checkPair(datePair) {
					validDates = append(validDates, datePair)
				} else {
					return false
				}
			} else {
				if checkPair(datePair) {
					allDates = append(allDates, datePair)
				} else {
					return false
				}
			}
		}
		for _, pair := range allDates {

			for _, validPair := range validDates {
				if checkLeft(validPair, pair) || checkRigth(validPair, pair) {
					continue
				} else {
					return false
				}

			}

			validDates = append(validDates, pair)

		}
		return true
	}
}

func checkLeft(validPair, newPair []time.Time) bool {
	if newPair[1].Before(validPair[0]) {
		return true
	}
	return false
}
func checkRigth(validPair, newPair []time.Time) bool {
	if newPair[0].After(validPair[1]) {
		return true
	}
	return false
}
func checkOver(validPair, newPair []time.Time) bool {
	if newPair[0].Before(validPair[0]) && newPair[1].After(validPair[1]) {
		return false
	}
	return true

}
func checkUnder(validPair, newPair []time.Time) bool {
	if newPair[0].After(validPair[0]) && newPair[1].Before(validPair[1]) {
		return false
	}
	return true

}
func checkPair(pairs []time.Time) bool {
	if pairs[0].Before(pairs[1]) {
		return true
	}
	return false
}

/*
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
			v.ValidateField("endDate", reservation.EndDate, DateNotBefore(reservation.EndDate, reservation.StartDate), DateNotInPast())
		foundErrors := v.GetErrors()
		if len(foundErrors) > 0 {
			for field, message := range foundErrors {
				fmt.Printf("%s: %s\n", field, message)
			}
		}
	}
*/
func (v *Validator) ValidateAvailability(reservation *domain.FreeReservation) {
	v.ValidateField("dateRange", reservation.DateRange, DateNotSame())
	foundErrors := v.GetErrors()
	if len(foundErrors) > 0 {
		for field, message := range foundErrors {
			fmt.Printf("%s: %s\n", field, message)
		}
	}
}

func getTimeFromString(str string) (time.Time, error) {
	layout := "2006-01-02"
	tm, err := time.Parse(layout, str)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %v", err)
	}
	return tm, nil
}
