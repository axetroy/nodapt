package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func List() error {
	if _, err := os.Stat(filepath.Join(virtual_node_env_dir, "node")); os.IsNotExist(err) {
		return nil
	}

	files, err := os.ReadDir(filepath.Join(virtual_node_env_dir, "node"))

	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		fName := file.Name()
		if file.IsDir() && strings.HasPrefix(fName, "node-v") {
			n := strings.SplitN(fName, "-", -1)

			version := n[1]
			fmt.Println(version)
		}
	}

	return nil
}
