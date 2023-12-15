package utils

import (
	"accommodations-service/domain"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
)

const (
	Name             = "Name can only contain letters and numbers, not special characters!"
	Location         = "Location can only contain letters!"
	Conveniences     = "Conveniences can only contain letters!"
	MinNumOfVisitors = "You need to input a number that is above 0, and needs to be lower that maximum value!"
	MaxNumOfVisitors = "You need to input a number lower than 100, and needs to be higher than minimum value!"
	StartDate        = "Start date is not a date format"
	EndDate          = "End date is not a date format"
)

var errorMessages = map[string]string{
	"Name":             Name,
	"Location":         Location,
	"Conveniences":     Conveniences,
	"MinNumOfVisitors": MinNumOfVisitors,
	"MaxNumOfVisitors": MaxNumOfVisitors,
	"StartDate":        StartDate,
	"EndDate":          EndDate,
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
func (v *Validator) ClearErrors() {
	v.Errors = make(map[string]string)
}

func IsName(value string) bool {
	nameRegex := `^[a-zA-Z0-9, ]*$`
	isValid, _ := regexp.MatchString(nameRegex, value)
	return isValid
}

func IsLocationOrConvenience(value string) bool {
	locationRegex := `^[a-zA-Z, ]+$`
	isValid, _ := regexp.MatchString(locationRegex, value)
	return isValid
}

func IsNumber(value string) bool {
	numberRegex := `^(?:[0-4]?[0-9]?[0-9]|500)$`
	isValid, _ := regexp.MatchString(numberRegex, value)
	return isValid
}

func MinLength(minLength int) ValidationRule {
	return func(value string) bool {
		return len(value) >= minLength
	}
}
func IsAddress(value string) bool {
	addressRegex := `^[a-zA-Z0-9 ,]+$`
	isValid, _ := regexp.MatchString(addressRegex, value)
	return isValid
}

func IsDateYYYYMMDD(input string) bool {
	dateRegex := regexp.MustCompile(`^(19|20)\d{2}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`)
	return dateRegex.MatchString(input)
}

func (v *Validator) ValidateAccommodation(accommodation *domain.Accommodation) {
	v.ValidateField("Name", accommodation.Name, MinLength(2), IsName)
	v.ValidateField("Address", accommodation.Address, MinLength(2), IsAddress)
	v.ValidateField("City", accommodation.City, MinLength(2), IsLocationOrConvenience)

	v.ValidateField("MinNumOfVisitors", strconv.Itoa(accommodation.MinNumOfVisitors), IsNumber)
	v.ValidateField("MaxNumOfVisitors", strconv.Itoa(accommodation.MaxNumOfVisitors), IsNumber)
	if accommodation.MinNumOfVisitors > accommodation.MaxNumOfVisitors {
		v.Errors["MaxNumOfVisitors"] = "Minimum number can not exceed maximum!"
	}

	foundErrors := v.GetErrors()

	if len(foundErrors) > 0 {
		for field, message := range foundErrors {
			fmt.Printf("%s: %s\n", field, message)
		}
	}
}

func (v *Validator) ValidateAvailabilities(availabilities *domain.CreateAccommodation) {
	layout := "2006-01-02" // Date layout format

	for _, availability := range availabilities.AvailableAccommodationDates {
		// Validate StartDate
		startDateParsed, err := time.Parse(layout, availability.StartDate)
		log.Println("Parsovan start date", startDateParsed)
		if err != nil {
			v.Errors["StartDate"] = "StartDate is not formatted correctly"
		}

		// Validate EndDate
		endDateParsed, err := time.Parse(layout, availability.EndDate)
		if err != nil {
			v.Errors["EndDate"] = "EndDate is not formatted correctly"
		}
		log.Println("Parsovan end date", endDateParsed)

		// Check date logic errors
		currentDate := time.Now()
		//log.Println(startDateParsed.After(endDateParsed))
		//log.Println(startDateParsed.Before(currentDate))
		//log.Println(endDateParsed.Before(currentDate))

		if startDateParsed.After(endDateParsed) || startDateParsed.Before(currentDate) || endDateParsed.Before(currentDate) {
			v.Errors["DateLogic"] = "Dates are not selected correctly"
		}

		// Validate Price
		v.ValidateField("MaxNumOfVisitors", strconv.Itoa(availability.Price), IsNumber)

		foundErrors := v.GetErrors()

		if len(foundErrors) > 0 {
			for field, message := range foundErrors {
				fmt.Printf("%s: %s\n", field, message)
			}
		}
	}

}

func (v *Validator) ValidateName(name string) {
	v.ValidateField("Name", name, MinLength(2), IsName)

	foundErrors := v.GetErrors()

	if len(foundErrors) > 0 {
		for field, message := range foundErrors {
			fmt.Printf("%s: %s\n", field, message)
		}
	}
}
func (v *Validator) ValidateLocation(location string) {
	v.ValidateField("Location", location, MinLength(2), IsLocationOrConvenience)

	foundErrors := v.GetErrors()

	if len(foundErrors) > 0 {
		for field, message := range foundErrors {
			fmt.Printf("%s: %s\n", field, message)
		}
	}
}
func (v *Validator) ValidateConvenience(convenience string) {
	v.ValidateField("Convenience", convenience, MinLength(2), IsLocationOrConvenience)

	foundErrors := v.GetErrors()

	if len(foundErrors) > 0 {
		for field, message := range foundErrors {
			fmt.Printf("%s: %s\n", field, message)
		}
	}
}

func (v *Validator) ValidateMinNum(minNumOfVisitors string) {
	v.ValidateField("MinNumOfVisitors", minNumOfVisitors, IsNumber)

	foundErrors := v.GetErrors()

	if len(foundErrors) > 0 {
		for field, message := range foundErrors {
			fmt.Printf("%s: %s\n", field, message)
		}
	}
}

func (v *Validator) ValidateMaxNum(maxNumOfVisitors string) {
	v.ValidateField("MinNumOfVisitors", maxNumOfVisitors, IsNumber)

	foundErrors := v.GetErrors()

	if len(foundErrors) > 0 {
		for field, message := range foundErrors {
			fmt.Printf("%s: %s\n", field, message)
		}
	}
}
