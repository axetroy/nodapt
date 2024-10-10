package command

import (
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
		panic(err)
	}

	nodapt_dir = filepath.Join(homeDir, ".nodapt")
}
