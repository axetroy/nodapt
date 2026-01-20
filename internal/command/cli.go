package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/axetroy/nodapt/internal/util"
)

var nodapt_dir string

func init() {

	nodaptDirFromEnv := util.GetEnvsWithFallback("", "NODE_ENV_DIR")

	if nodaptDirFromEnv != "" {
		nodapt_dir = nodaptDirFromEnv
		return
	}

	homeDir, err := os.UserHomeDir()

	if err != nil {
		// Fallback to temporary directory if home directory is not available
		fmt.Fprintf(os.Stderr, "Warning: Unable to get user home directory: %v\n", err)
		fmt.Fprintf(os.Stderr, "Using temporary directory as fallback\n")
		nodapt_dir = filepath.Join(os.TempDir(), ".nodapt")
		return
	}

	nodapt_dir = filepath.Join(homeDir, ".nodapt")
}
