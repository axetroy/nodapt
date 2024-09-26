//go:build windows

package util

import (
	"os"
	"path/filepath"
	"strings"
)

// executableExtensions is a set of common executable file extensions on Windows.
var executableExtensions = map[string]struct{}{
	".exe": {},
	".bat": {},
	".cmd": {},
}

// isExecutable checks if a file is considered executable on Windows based on its extension.
func isExecutable(fileInfo os.FileInfo, filePath string) bool {
	// Skip directories
	if fileInfo.IsDir() {
		return false
	}

	// Get the file extension
	ext := filepath.Ext(filePath)

	// Check if the extension matches a known executable extension (case-insensitive)
	_, exists := executableExtensions[strings.ToLower(ext)]

	return exists
}
