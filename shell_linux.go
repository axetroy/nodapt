package virtualnodeenv

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func getShell() (string, error) {
	out, err := exec.Command("ps", "-p", fmt.Sprint(os.Getppid()), "-o", "comm=").Output()

	if err != nil {
		return "", errors.WithStack(err)
	}

	shellName := strings.TrimSpace(string(out))

	return shellName, nil
}
