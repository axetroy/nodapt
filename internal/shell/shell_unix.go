//go:build !windows

package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func getShellFromEnv() *string {
	shell := os.Getenv("SHELL")

	if shell == "" {
		return nil
	}

	return &shell
}

var knownShells = []string{"sh", "bash", "zsh", "fish", "dash", "ash", "ksh"}

func isKnownShell(shellPath string) bool {
	parts := strings.Split(shellPath, "/")

	executableFileName := parts[len(parts)-1]

	return slices.Contains(knownShells, executableFileName)
}

func getParentProcessInfo(pid int) (int, string, error) {
	cmd := exec.Command("ps", "-o", "ppid,comm", "-p", fmt.Sprintf("%d", pid))
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return -1, "", err
	}

	lines := strings.Split(out.String(), "\n")
	if len(lines) < 2 {
		return -1, "", fmt.Errorf("no output from ps")
	}

	fmt.Printf("%+v\n", lines)

	fields := strings.Fields(lines[1])
	if len(fields) < 2 {
		return -1, "", fmt.Errorf("unexpected output format")
	}

	ppid := fields[0]
	command := fields[1]

	ppidInt, err := strconv.Atoi(ppid)

	if err != nil {
		return -1, "", err
	}

	return ppidInt, command, nil
}

func getShellFromProcess(pid int) (string, error) {
	for {
		ppid, command, err := getParentProcessInfo(pid)

		if err != nil {
			return "", err
		}

		if isKnownShell(command) {
			return command, nil
		}

		pid = ppid
	}

}

func GetPath() (string, error) {
	if shell, err := getShellFromProcess(os.Getppid()); err == nil {
		if path.IsAbs(shell) {
			return shell, nil
		} else {
			if fullShellPath, err := exec.LookPath(shell); err == nil {
				return fullShellPath, nil
			}
		}
	}

	if shell := getShellFromEnv(); shell != nil {
		return *shell, nil
	}

	return "", errors.New("can't detect shell")
}
