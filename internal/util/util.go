package util

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
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
	// Iterate through the provided keys
	for _, key := range keys {
		// Prepare the uppercase and lowercase versions of the key
		keysToCheck := []string{strings.ToUpper(key), strings.ToLower(key)}

		// Check both versions for a non-empty value
		for _, k := range keysToCheck {
			if value := os.Getenv(k); value != "" {
				return value
			}
		}
	}

	return fallback
}

// FindExecutable searches for an executable file in the specified directory.
// It checks for the presence of a file with the given executable name and
// appropriate extensions based on the operating system (e.g., ".exe", ".bat", ".cmd"
// for Windows). The function performs a case-insensitive comparison on Windows and macOS,
// while it is case-sensitive on other operating systems.
//
// Parameters:
//   - dir: The directory path to search for the executable.
//   - executableName: The name of the executable file without the extension.
//
// Returns:
//   - bool: A boolean indicating whether the executable was found.
//   - error: An error if the directory cannot be read or any other issue occurs during the search.
func FindExecutable(dir, executableName string) (bool, error) {
	// Read directory entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, errors.WithStack(err)
	}

	// Determine extensions and case sensitivity based on OS
	var executableExtensions []string
	isCaseInsensitive := false

	switch runtime.GOOS {
	case "windows":
		executableExtensions = []string{".exe", ".bat", ".cmd"}
		isCaseInsensitive = true
	case "darwin":
		isCaseInsensitive = true
		executableExtensions = []string{""}
	default:
		executableExtensions = []string{""}
	}

	// Normalize the executable name for case-insensitive comparison
	if isCaseInsensitive {
		executableName = strings.ToLower(executableName)
	}

	// Check each file in the directory
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Compare the file name with possible executable names
		fileName := entry.Name()

		if isCaseInsensitive {
			fileName = strings.ToLower(fileName)
		}

		for _, ext := range executableExtensions {
			if fileName == executableName+ext {
				info, err := entry.Info()

				if err != nil {
					return false, errors.WithStack(err)
				}

				return isExecutable(info, filepath.Join(dir, fileName)), nil
			}
		}
	}

	return false, nil
}
