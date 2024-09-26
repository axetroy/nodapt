package util

import (
	"os"
	"testing"
)

func TestGetEnvsWithFallback(t *testing.T) {
	// Set up environment variables for testing
	os.Setenv("TEST_KEY_UPPER", "upper_value")
	os.Setenv("test_key_lower", "lower_value")
	defer os.Unsetenv("TEST_KEY_UPPER")
	defer os.Unsetenv("test_key_lower")

	tests := []struct {
		name     string
		fallback string
		keys     []string
		expected string
	}{
		{
			name:     "Upper case key exists",
			fallback: "default_value",
			keys:     []string{"TEST_KEY_UPPER"},
			expected: "upper_value",
		},
		{
			name:     "Lower case key exists",
			fallback: "default_value",
			keys:     []string{"test_key_lower"},
			expected: "lower_value",
		},
		{
			name:     "Mixed case keys, upper case exists",
			fallback: "default_value",
			keys:     []string{"test_key_upper", "TEST_KEY_UPPER"},
			expected: "upper_value",
		},
		{
			name:     "Mixed case keys, lower case exists",
			fallback: "default_value",
			keys:     []string{"TEST_KEY_LOWER", "test_key_lower"},
			expected: "lower_value",
		},
		{
			name:     "No matching keys, return fallback",
			fallback: "default_value",
			keys:     []string{"NON_EXISTENT_KEY"},
			expected: "default_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEnvsWithFallback(tt.fallback, tt.keys...)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}
