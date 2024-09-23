//go:build linux || darwin

package shell

import (
	"path/filepath"
	"testing"
)

func Test_getShell(t *testing.T) {
	shellPath, err := GetPath()

	if err != nil {
		t.Errorf("getShell() error = %v", err)
		return
	}

	t.Logf("getShell() = %v", shellPath)

	if shellPath != "go" && !filepath.IsAbs(shellPath) {
		t.Errorf("getShell() expected to return an absolute path, got %v", shellPath)
	}

}
