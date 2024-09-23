package virtualnodeenv

import (
	"github.com/pkg/errors"

	"github.com/Masterminds/semver"
)

func VersionMatch(versionRange string, version string) (bool, error) {
	c, err := semver.NewConstraint(versionRange)

	if err != nil {
		return false, errors.WithMessage(err, "failed to parse version range")
	}

	v, err := semver.NewVersion(version)

	if err != nil {
		return false, errors.WithMessage(err, "failed to parse version")
	}

	return c.Check(v), nil
}
