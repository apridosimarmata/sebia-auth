package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidatePhoneNumber checks if the provided phone number is valid.
func ValidatePhoneNumber(phoneNumber string) (newPhoneNumber *string, err error) {
	phoneRegex := regexp.MustCompile(`^(62|08)\d{7,11}$`)
	if phoneRegex.MatchString(phoneNumber) {
		if phoneNumber[0] == '0' {
			phoneNumber = "62" + phoneNumber[1:]

		}

		return &phoneNumber, nil
	}

	return nil, fmt.Errorf("Diawali dengan 62 atau 08")
}

// ValidateEmail checks if the provided email is valid.
func ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if emailRegex.MatchString(email) {
		return nil // Valid email
	}
	return fmt.Errorf("Format email belum benar") // Invalid email
}

// ValidatePassword checks if the provided password is valid (at least 8 characters).
func ValidatePassword(password string) error {
	if len(password) >= 8 {
		return nil // Valid password
	}
	return fmt.Errorf("Kata sandi minimal 8 karakter") // Invalid password
}

// ValidateRequired checks if the provided value is not empty.
func ValidateRequired(value string) error {
	if strings.TrimSpace(value) != "" {
		return nil
	}
	return fmt.Errorf("Empty value provided")
}

// ValidateRequired checks if the provided value is not empty.
func ValidateRequiredSlice(value []interface{}) error {
	if len(value) < 1 {
		return nil
	}
	return fmt.Errorf("Empty value provided")
}

// ValidateRequired checks if the provided value is not empty.
func ValidateRequiredInt(value int) error {
	if value > 0 {
		return nil
	}

	return fmt.Errorf("Empty value provided")
}

func ValidateRequiredIntAllowsZero(value int) error {
	if value > -1 {
		return nil
	}

	return fmt.Errorf("Empty value provided")
}

// ValidateFullName checks if the provided full name is valid.
func ValidateFullName(fullName string) error {
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	if nameRegex.MatchString(fullName) {
		return nil
	}
	return fmt.Errorf("Nama belum benar")
}
