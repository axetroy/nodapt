package virtualnodeenv

import (
	"strings"
	"testing"
)

func Test_getShell(t *testing.T) {
	shellPath, err := getShell()

	if err != nil {
		t.Errorf("getShell() error = %v", err)
		return
	}

	t.Logf("getShell() = %v", shellPath)

	if !strings.HasPrefix(shellPath, "/") && shellPath != "go" {
		t.Errorf("getShell() expected to return an absolute path, got %v", shellPath)
	}

}
