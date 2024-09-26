package util

import (
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
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

func TestIsExecutable(t *testing.T) {
	tests := []struct {
		name     string
		fileInfo os.FileInfo
		filePath string
		expected bool
	}{
		{
			name:     "Windows executable file",
			fileInfo: mockFileInfo{mode: 0777, isDir: false},
			filePath: "test.exe",
			expected: true,
		},
		{
			name:     "Windows non-executable file",
			fileInfo: mockFileInfo{mode: 0666, isDir: false},
			filePath: "test.txt",
			expected: false,
		},
		{
			name:     "Unix executable file",
			fileInfo: mockFileInfo{mode: 0755, isDir: false},
			filePath: "test.sh",
			expected: true,
		},
		{
			name:     "Unix non-executable file",
			fileInfo: mockFileInfo{mode: 0644, isDir: false},
			filePath: "test.txt",
			expected: false,
		},
		{
			name:     "Directory",
			fileInfo: mockFileInfo{mode: 0755, isDir: true},
			filePath: "testdir",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if runtime.GOOS == "windows" && !strings.HasSuffix(tt.filePath, ".exe") {
				tt.expected = false
			}
			result := IsExecutable(tt.fileInfo, tt.filePath)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// mockFileInfo is a mock implementation of os.FileInfo for testing purposes.
type mockFileInfo struct {
	mode  os.FileMode
	isDir bool
}

func (m mockFileInfo) Name() string       { return "" }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return nil }
