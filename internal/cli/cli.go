package cli

import (
	"os"
	"path/filepath"

	"github.com/axetroy/virtual_node_env/internal/util"
)

var virtual_node_env_dir string

func init() {

	virtualNodeEnvDirFromEnv := util.GetEnvsWithFallback("", "NODE_ENV_DIR")

	if virtualNodeEnvDirFromEnv != "" {
		virtual_node_env_dir = virtualNodeEnvDirFromEnv
		return
	}

	homeDir, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	virtual_node_env_dir = filepath.Join(homeDir, ".nodapt")
}
