package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSalt(size int) (string, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func GenerateRandomString(n int) (string, error) {
	// Create a slice of random bytes
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Encode the bytes to a URL-safe base64 string
	return base64.URLEncoding.EncodeToString(b)[:n], nil
}
