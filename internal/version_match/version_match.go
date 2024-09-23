package version_match

import (
	"github.com/pkg/errors"

	"github.com/Masterminds/semver"
)

func VersionMatch(semverVersionConstraint string, version string) (bool, error) {
	c, err := semver.NewConstraint(semverVersionConstraint)

	if err != nil {
		return false, errors.WithMessage(err, "failed to parse version range")
	}

	v, err := semver.NewVersion(version)

	if err != nil {
		return false, errors.WithMessage(err, "failed to parse version")
	}

	return c.Check(v), nil
}
