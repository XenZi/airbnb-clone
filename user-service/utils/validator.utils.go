package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"user-service/domain"
)

const (
	Email     = "Email format is not correct"
	Username  = "Username must be longer than 6 characters."
	FirstName = "Name must start with a capital letter"
	LastName  = "Last name must start with a capital letter"
	Residence = "Residence format is not correct"
	Role      = "Not an approved role"
	Age       = "Not an appropriate age"
)

var errorMessages = map[string]string{
	"Email":     Email,
	"Username":  Username,
	"FirstName": FirstName,
	"LastName":  LastName,
	"Residence": Residence,
	"Role":      Role,
	"Age":       Age,
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

type IntValidationRule func(value int) bool

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

func (v *Validator) ValidateIntField(fieldName string, value int, rules ...IntValidationRule) {
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
}

func (v *Validator) GetErrors() map[string]string {
	return v.Errors
}

func IsEmail(value string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	isValid, _ := regexp.MatchString(emailRegex, value)
	return isValid
}

func MinLength(minLength int) ValidationRule {
	return func(value string) bool {
		return len(value) >= minLength
	}
}

func MinDigits(minDigits int) ValidationRule {
	return func(value string) bool {
		digitCount := 0
		for _, char := range value {
			if strings.Contains("0123456789", string(char)) {
				digitCount++
			}
		}
		return digitCount >= minDigits
	}
}

func MinSpecialChars(minSpecialChars int) ValidationRule {
	return func(value string) bool {
		specialCharCount := 0
		for _, char := range value {
			if strings.Contains("!@#$%^&*()-_+=<>?/", string(char)) {
				specialCharCount++
			}
		}
		return specialCharCount >= minSpecialChars
	}
}

func MinUpperCaseChars(minUpperCaseChars int) ValidationRule {
	return func(value string) bool {
		specialCharCount := 0
		for _, char := range value {
			if strings.Contains("QWERTYUIOPASDFGHJKLZXCVBNM", string(char)) {
				specialCharCount++
			}
		}
		return specialCharCount >= minUpperCaseChars
	}
}

func StartsWithCapital() ValidationRule {
	return func(value string) bool {
		runes := []rune(value)
		return runes[0] == unicode.ToUpper(runes[0])
	}
}

func ValidAge(minAge int) IntValidationRule {
	return func(value int) bool {
		return value > minAge
	}

}

func (v *Validator) ValidateUser(newUser *domain.CreateUser) {
	v.ValidateField("FirstName", newUser.FirstName, MinLength(2), StartsWithCapital())
	v.ValidateField("LastName", newUser.LastName, MinLength(2), StartsWithCapital())
	v.ValidateField("Email", newUser.Email, IsEmail)
	v.ValidateField("Residence", newUser.Residence, MinLength(3), MinDigits(1))
	v.ValidateIntField("Age", newUser.Age, ValidAge(18))
	foundErrors := v.GetErrors()
	if len(foundErrors) > 0 {
		for field, message := range foundErrors {
			fmt.Printf("%s: %s\n", field, message)
		}
	}
}
