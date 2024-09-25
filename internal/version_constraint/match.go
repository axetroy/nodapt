package version_constraint

import (
	"github.com/pkg/errors"

	"github.com/Masterminds/semver"
)

// Match checks if the given version satisfies the specified version constraint.
// It takes a constraint string and a version string as input parameters.
//
// Parameters:
//   - constraint: A string representing the version range to check against.
//   - version: A string representing the version to be checked.
//
// Returns:
//   - bool: true if the version satisfies the constraint, false otherwise.
//   - error: An error if the constraint or version cannot be parsed, with a descriptive message.
func Match(constraint string, version string) (bool, error) {
	c, err := semver.NewConstraint(constraint)

	if err != nil {
		return false, errors.WithMessage(err, "failed to parse version range")
	}

	v, err := semver.NewVersion(version)

	if err != nil {
		return false, errors.WithMessage(err, "failed to parse version")
	}

	return c.Check(v), nil
}
