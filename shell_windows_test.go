package virtualnodeenv

import (
	"path/filepath"
	"testing"
)

func Test_getShell(t *testing.T) {
	shellPath, err := getShell()

	if err != nil {
		t.Errorf("getShell() error = %v", err)
		return
	}

	t.Logf("getShell() = %v", shellPath)

	if shellPath != "go.exe" && !filepath.IsAbs(shellPath) {
		t.Errorf("getShell() expected to return an absolute path, got %v", shellPath)
	}
}
