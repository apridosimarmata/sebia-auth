package utils

import "strings"

// MaskName partially masks a name, e.g., "John Doe" -> "Jo** D*e"
func MaskName(name string) string {
	words := strings.Fields(name) // Split by space
	for i, word := range words {
		if len(word) > 2 {
			words[i] = word[:2] + strings.Repeat("*", len(word)-3) + word[len(word)-1:]
		} else if len(word) > 1 {
			words[i] = word[:1] + "*"
		} // Single letter words remain as is
	}
	return strings.Join(words, " ")
}

// MaskPhone partially masks a phone number, e.g., "081212341234" -> "08*****1234"
func MaskPhone(phone string) string {
	if len(phone) > 6 {
		return "+" + phone[:3] + strings.Repeat("*", len(phone)-6) + phone[len(phone)-4:]
	}
	return "+" + phone
}

// MaskEmail partially masks an email, e.g., "johndoe@example.com" -> "j****e@example.com"
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email // Invalid email format, return as is
	}

	localPart := parts[0]
	domain := parts[1]

	if len(localPart) <= 2 {
		// If the local part is too short, mask the middle entirely
		return localPart[:1] + strings.Repeat("*", len(localPart)-1) + "@" + domain
	}

	// Mask all but the first and last character of the local part
	maskedLocal := localPart[:1] + strings.Repeat("*", len(localPart)-2) + localPart[len(localPart)-1:]

	return maskedLocal + "@" + domain
}
