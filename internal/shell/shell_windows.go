package shell

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

func GetPath() (string, error) {
	var shellPath string

	if os.Getenv("SHELL") != "" {
		shellPath = os.Getenv("SHELL")
	} else if os.Getenv("COMSPEC") != "" {
		// Windows
		comspec := os.Getenv("COMSPEC")
		if strings.Contains(strings.ToLower(comspec), "powershell") {
			shellPath = comspec
		} else if strings.Contains(strings.ToLower(comspec), "cmd") {
			shellPath = comspec
		}
	}

	if shellPath == "" {
		return "", errors.New("Unknown shell")
	}

	return shellPath, nil
}
