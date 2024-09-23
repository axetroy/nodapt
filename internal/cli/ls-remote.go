package cli

import (
	"github.com/axetroy/virtual_node_env/internal/node"
	"github.com/pkg/errors"
)

func ListRemote() error {
	versions, err := node.GetAllVersions()

	if err != nil {
		return errors.WithMessage(err, "failed to get node versions")
	}

	for _, version := range versions {
		println(version)
	}

	return nil
}
