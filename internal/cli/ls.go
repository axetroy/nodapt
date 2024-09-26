package cli

import (
	"fmt"

	"github.com/axetroy/virtual_node_env/internal/node"
	"github.com/pkg/errors"
)

func List() error {
	cached, err := node.GetCachedVersions(virtual_node_env_dir)

	if err != nil {
		return errors.WithStack(err)
	}

	for _, c := range cached {
		fmt.Println(c.Version)
	}

	return nil
}
