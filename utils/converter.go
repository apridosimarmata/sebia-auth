package utils

import (
	"strings"
)

func ConvertPhoneNumber(phone string) string {
	// Check if the phone number starts with '0'
	if strings.HasPrefix(phone, "0") {
		// Replace the first '0' with '62'
		return "62" + phone[1:]
	}
	// If it doesn't start with '0', return it as is
	return phone
}
