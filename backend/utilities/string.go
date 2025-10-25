package utilities

import (
	"regexp"
	"strings"
)

// ToSnakeCase converts a string to snake_case
// e.g., "Digital Wallet" -> "digital_wallet"
func ToSnakeCase(str string) string {
	// Replace spaces with underscores
	result := strings.ReplaceAll(str, " ", "_")

	// Convert to lowercase
	result = strings.ToLower(result)

	// Remove any characters that aren't alphanumeric or underscore
	reg := regexp.MustCompile("[^a-z0-9_]+")
	result = reg.ReplaceAllString(result, "")

	return result
}
