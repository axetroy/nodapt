package cli

import (
	"fmt"

	"github.com/axetroy/virtual_node_env/internal/node"
	"github.com/pkg/errors"
)

func ListRemote() error {
	versions, err := node.GetAllVersions()

	if err != nil {
		return errors.WithMessage(err, "failed to get node versions")
	}

	for _, version := range versions {
		fmt.Println(version.Version)
	}

	return nil
}
