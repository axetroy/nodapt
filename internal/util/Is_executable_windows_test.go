//go:build windows

package util

import (
	"os"
	"testing"
	"time"
)

// mockFileInfo is a mock implementation of os.FileInfo for testing purposes.
type mockFileInfo struct {
	name string
	dir  bool
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() os.FileMode  { return 0 }
func (m mockFileInfo) ModTime() time.Time { return time.Time{} }
func (m mockFileInfo) IsDir() bool        { return m.dir }
func (m mockFileInfo) Sys() interface{}   { return nil }

func TestIsExecutable(t *testing.T) {
	tests := []struct {
		name     string
		fileInfo os.FileInfo
		filePath string
		want     bool
	}{
		{"Executable .exe", mockFileInfo{name: "test.exe", dir: false}, "test.exe", true},
		{"Executable .bat", mockFileInfo{name: "test.bat", dir: false}, "test.bat", true},
		{"Executable .cmd", mockFileInfo{name: "test.cmd", dir: false}, "test.cmd", true},
		{"Non-executable .txt", mockFileInfo{name: "test.txt", dir: false}, "test.txt", false},
		{"Non-executable directory", mockFileInfo{name: "testdir", dir: true}, "testdir", false},
		{"Executable .EXE (case-insensitive)", mockFileInfo{name: "test.EXE", dir: false}, "test.EXE", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isExecutable(tt.fileInfo, tt.filePath); got != tt.want {
				t.Errorf("isExecutable() = %v, want %v", got, tt.want)
			}
		})
	}
}
