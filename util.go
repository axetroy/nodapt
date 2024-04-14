package virtualnodeenv

import (
	"os"
	"strings"
)

func getEnvsWithFallback(fallback string, keys ...string) string {
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
