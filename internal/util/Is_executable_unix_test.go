//go:build unix

package util

import (
	"os"
	"testing"
	"time"
)

func TestIsExecutable(t *testing.T) {
	tests := []struct {
		name     string
		fileInfo os.FileInfo
		expected bool
	}{
		{
			name: "directory",
			fileInfo: &mockFileInfo{
				isDir: true,
				mode:  0755,
			},
			expected: false,
		},
		{
			name: "non-executable file",
			fileInfo: &mockFileInfo{
				isDir: false,
				mode:  0644,
			},
			expected: false,
		},
		{
			name: "executable file",
			fileInfo: &mockFileInfo{
				isDir: false,
				mode:  0755,
			},
			expected: true,
		},
		{
			name: "executable by owner only",
			fileInfo: &mockFileInfo{
				isDir: false,
				mode:  0700,
			},
			expected: true,
		},
		{
			name: "executable by group only",
			fileInfo: &mockFileInfo{
				isDir: false,
				mode:  0070,
			},
			expected: true,
		},
		{
			name: "executable by others only",
			fileInfo: &mockFileInfo{
				isDir: false,
				mode:  0001,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isExecutable(tt.fileInfo, "")
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// mockFileInfo is a mock implementation of os.FileInfo for testing purposes.
type mockFileInfo struct {
	isDir bool
	mode  os.FileMode
}

func (m *mockFileInfo) Name() string       { return "" }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() os.FileMode  { return m.mode }
func (m *mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m *mockFileInfo) IsDir() bool        { return m.isDir }
func (m *mockFileInfo) Sys() interface{}   { return nil }
