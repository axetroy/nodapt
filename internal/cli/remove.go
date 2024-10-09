package cli

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver"
	"github.com/axetroy/nodapt/internal/node"
	"github.com/pkg/errors"
)

func Remove(constraint string) error {
	cachedNodes, err := node.GetCachedVersions(virtual_node_env_dir)

	if err != nil {
		return errors.WithStack(err)
	}

	constraintVer, err := semver.NewConstraint(constraint)

	if err != nil {
		return errors.WithStack(err)
	}

	for _, cache := range cachedNodes {
		v, err := semver.NewVersion(cache.Version)

		if err != nil {
			return errors.WithStack(err)
		}

		if constraintVer.Check(v) {
			err := os.RemoveAll(cache.FilePath)

			if err != nil {
				return errors.WithStack(err)
			}

			fmt.Fprintf(os.Stderr, "Node version %s has been removed\n", cache.Version)
		}
	}

	return nil
}
