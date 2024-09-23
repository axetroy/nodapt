package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoopUpFile(t *testing.T) {
	// Setup temporary directories and files for testing
	rootDir := t.TempDir()

	subDirectory := filepath.Join(rootDir, "subdir")
	_ = os.Mkdir(subDirectory, 0755)

	configFileName := "config.yaml"
	configFilePath := filepath.Join(rootDir, configFileName)
	file, err := os.Create(configFilePath)

	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	defer file.Close()

	tests := []struct {
		name     string
		root     string
		fileName string
		expected *string
	}{
		{
			name:     "File exists in root directory",
			root:     rootDir,
			fileName: configFileName,
			expected: &configFilePath,
		},
		{
			name:     "File exists in parent directory",
			root:     subDirectory,
			fileName: configFileName,
			expected: &configFilePath,
		},
		{
			name:     "File does not exist",
			root:     subDirectory,
			fileName: "nonexistent.yaml",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LoopUpFile(tt.root, tt.fileName)
			if (result == nil && tt.expected != nil) || (result != nil && tt.expected == nil) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			} else if result != nil && *result != *tt.expected {
				t.Errorf("Expected %v, got %v", *tt.expected, *result)
			}
		})
	}
}
