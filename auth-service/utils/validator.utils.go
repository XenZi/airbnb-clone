package utils

import (
	"auth-service/domains"
	"fmt"
	"regexp"
	"strings"
)

const (
	Email    = "Email format is not correct"
	Password = "Password must contain one uppercase, one special character, one number and minimum length must be 8."
	Username = "Username must be longer than 6 characters."
)

var errorMessages = map[string]string{
	"Email":    Email,
	"Password": Password,
	"Username": Username,
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

func (v *Validator) ValidateRegisterUser(registerUser *domains.RegisterUser) {
	v.ValidateField("FirstName", registerUser.FirstName, MinLength(2))
	v.ValidateField("LastName", registerUser.LastName, MinLength(2))
	v.ValidateField("Email", registerUser.Email, IsEmail)
	v.ValidateField("Password", registerUser.Password, MinLength(6), MinDigits(1), MinSpecialChars(1), MinUpperCaseChars(1))
	v.ValidateField("CurrentPlace", registerUser.CurrentPlace, MinLength(3))
	foundErrors := v.GetErrors()
	if len(foundErrors) > 0 {
		for field, message := range foundErrors {
			fmt.Printf("%s: %s\n", field, message)
		}
	}
}

func (v *Validator) ValidatePassword(password string) {
	v.ValidateField("Password", password, MinLength(6), MinDigits(1), MinSpecialChars(1), MinUpperCaseChars(1))
	foundErrors := v.GetErrors()
	if len(foundErrors) > 0 {
		for field, message := range foundErrors {
			fmt.Printf("%s: %s\n", field, message)
		}
	}
}