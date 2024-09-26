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
	entries, err := os.ReadDir(dir)

	if err != nil {
		return false, err
	}

	executableExtensions := []string{}

	if runtime.GOOS == "windows" {
		executableExtensions = []string{".exe", ".bat", ".cmd"}
	}

	// 默认情况下，windows 和 macOS 都对大小写不敏感
	isCaseInsensitive := runtime.GOOS == "windows" || runtime.GOOS == "darwin"

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()

	inner:
		for _, ext := range executableExtensions {
			var isSameFile bool

			if isCaseInsensitive {
				// 大小写不敏感的比较
				isSameFile = strings.EqualFold(fileName, executableName+ext)
			} else {
				// 大小写敏感的比较
				isSameFile = fileName == executableName+ext
			}

			if isSameFile {
				info, err := entry.Info()

				if err != nil {
					continue inner
				}

				return IsExecutable(info, fileName), nil
			}
		}
	}

	return false, nil
}
