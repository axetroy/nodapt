//go:build !windows

package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func getFishShell() *string {
	cmd := exec.Command("status", "fish-path")

	out, err := cmd.Output()

	if err != nil {
		return nil
	}

	shell := strings.TrimSpace(string(out))

	return &shell
}

func getShellFromEnv() *string {
	shell := os.Getenv("SHELL")

	if shell == "" {
		return nil
	}

	return &shell
}

func getShellFromProcess() (string, error) {
	out, err := exec.Command("ps", "-p", fmt.Sprint(os.Getppid()), "-o", "comm=").Output()

	if err != nil {
		return "", errors.WithStack(err)
	}

	shellName := strings.TrimSpace(string(out))

	return shellName, nil
}

func GetPath() (string, error) {
	if shell := getFishShell(); shell != nil {
		return *shell, nil
	}

	if shell := getShellFromEnv(); shell != nil {
		return *shell, nil
	}

	return getShellFromProcess()
}
