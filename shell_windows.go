package virtualnodeenv

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

func getShell() (string, error) {
	var shellPath string

	if os.Getenv("COMSPEC") != "" {
		// Windows
		comspec := os.Getenv("COMSPEC")
		if strings.Contains(strings.ToLower(comspec), "powershell") {
			shellPath = comspec
		} else if strings.Contains(strings.ToLower(comspec), "cmd") {
			shellPath = comspec
		} else {
			return "", errors.New("Unknown shell")
		}
	}

	if shellPath != "" {
		return shellPath, nil
	}

	return "", errors.New("Unknown shell")
}