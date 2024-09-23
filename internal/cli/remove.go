package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/axetroy/virtual_node_env/internal/node"
)

func Remove(version string) error {
	target := node.GetArtifactName(version)

	dest := filepath.Join(virtual_node_env_dir, "node", target)

	// 检查文件是否存在
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Node version %s not found\n", version)
		return nil
	}

	return os.RemoveAll(dest)
}
