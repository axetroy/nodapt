package util

import (
	"os"
	"strings"
)

// GetEnvsWithFallback retrieves the value of the specified environment variables,
// checking both uppercase and lowercase versions of the keys. If none of the
// specified keys are found, it returns the provided fallback value.
//
// Parameters:
//   - fallback: The value to return if none of the environment variables are set.
//   - keys: A variadic list of environment variable names to check.
//
// Returns:
//   - The value of the first found environment variable or the fallback value if none are found.
func GetEnvsWithFallback(fallback string, keys ...string) string {
	for _, key := range keys {
		if value := os.Getenv(strings.ToUpper(key)); value != "" {
			return value
		}

		if value := os.Getenv(strings.ToLower(key)); value != "" {
			return value
		}
	}

	return fallback
}
