package utils

import (
	"crypto/rand"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/oklog/ulid/v2"
	"golang.org/x/text/runes"
	"golang.org/x/text/unicode/norm"
)

// GenerateSlug converts a string into a URL-friendly slug
func GenerateSlug(input string) string {
	// Normalize the string to decompose accented characters
	input = norm.NFKD.String(input)

	// Remove non-ASCII characters while retaining the base characters
	input = runes.Remove(runes.In(unicode.Mn)).String(input)

	// Convert to lowercase
	input = strings.ToLower(input)

	// Replace spaces and underscores with dashes
	input = strings.ReplaceAll(input, " ", "-")
	input = strings.ReplaceAll(input, "_", "-")

	// Remove non-alphanumeric characters except dashes
	reg, _ := regexp.Compile("[^a-z0-9-]+")
	input = reg.ReplaceAllString(input, "")

	// Remove leading or trailing dashes
	input = strings.Trim(input, "-")

	return input + "-" + strings.ToLower(GenerateUniqueId()[:10])
}

func GenerateUniqueId() string {
	t := time.Now().UTC()

	// Create an entropy source for random number generation TODO
	entropy := ulid.Monotonic(rand.Reader, 0)

	// Generate a ULID
	id := ulid.MustNew(ulid.Timestamp(t), entropy)

	return strings.ToLower(id.String())
}
