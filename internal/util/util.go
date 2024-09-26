package util

import (
	"os"
	"path/filepath"
	"runtime"
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

func IsExecutable(fileInfo os.FileInfo, filePath string) bool {
	if runtime.GOOS == "windows" {
		// Windows: check for common executable extensions
		ext := filepath.Ext(filePath)
		executableExtensions := []string{".exe", ".bat", ".cmd"}
		for _, e := range executableExtensions {
			if strings.EqualFold(ext, e) {
				return true
			}
		}
		return false
	}
	// Unix-based systems: check executable permission
	mode := fileInfo.Mode()
	return !fileInfo.IsDir() && (mode&0111 != 0) // Check if any execute bit is set (owner, group, or others)
}

// FindExecutable checks if a specified file exists in a given directory
// and determines if it is executable.
//
// Parameters:
//   - dir: The directory in which to search for the file.
//   - fileName: The name of the file to check for.
//
// Returns:
//   - bool: True if the file exists and is executable, false otherwise.
//   - error: An error if there was an issue checking the file, or nil if no error occurred.
func FindExecutable(dir, fileName string) (bool, error) {
	// Get the full path of the file in the directory
	filePath := filepath.Join(dir, fileName)

	// Check if the file exists and retrieve its info
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	// Check if the file is executable
	if IsExecutable(fileInfo, filePath) {
		return true, nil
	}

	return false, nil
}
