package cli

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/axetroy/virtual_node_env/internal/util"
	"github.com/pkg/errors"
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

	virtual_node_env_dir = filepath.Join(homeDir, ".virtual-node-env")
}

type Config struct {
	Node string `json:"node"`
}

func LoadConfig(filePath string) (*Config, error) {
	content, err := os.ReadFile(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.WithMessage(err, "Can not found `.virtual-node-env.json` file.")
		}

		return nil, errors.WithStack(err)
	}

	c := &Config{}

	if err := json.Unmarshal(content, c); err != nil {
		return nil, errors.WithMessagef(err, "Read config file %s failed", filePath)
	}

	return c, nil
}
