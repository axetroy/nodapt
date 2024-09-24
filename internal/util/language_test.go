package util

import (
	"os"
	"testing"
)

func TestGetLanguageViaEnv(t *testing.T) {
	// Test case when LANG environment variable is set
	t.Run("LANG environment variable is set", func(t *testing.T) {
		expectedLang := "en-US"
		os.Setenv("LANG", expectedLang)
		defer os.Unsetenv("LANG")

		lang := getLanguageViaEnv()
		if lang == nil || *lang != expectedLang {
			t.Errorf("Expected %s, but got %v", expectedLang, lang)
		}
	})

	// Test case when LANG environment variable is not set
	t.Run("LANG environment variable is not set", func(t *testing.T) {
		os.Unsetenv("LANG")

		lang := getLanguageViaEnv()
		if lang != nil {
			t.Errorf("Expected nil, but got %v", lang)
		}
	})

	// Test case when LANG environment variable is empty
	t.Run("LANG environment variable is empty", func(t *testing.T) {
		os.Setenv("LANG", "")
		defer os.Unsetenv("LANG")

		lang := getLanguageViaEnv()
		if lang != nil {
			t.Errorf("Expected nil, but got %v", lang)
		}
	})
}
