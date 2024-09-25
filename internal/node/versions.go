package node

import (
	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/axetroy/virtual_node_env/internal/version_match"
	"github.com/pkg/errors"
)

type Versions []struct {
	Version string `json:"version"`
	LTS     any    `json:"lts"`
}

// GetAllVersions retrieves a list of all available Node.js versions from the official Node.js distribution index.
// It makes an HTTP GET request to the Node.js API endpoint and decodes the JSON response into a slice of version strings.
//
// Returns:
// - A slice of strings containing the Node.js versions, or
// - An error if the request fails or if there is an issue decoding the response.
func GetAllVersions() (Versions, error) {
	resp, err := http.Get(NODE_MIRROR + "index.json")

	if err != nil {
		return nil, errors.WithMessage(err, "failed to get node versions")
	}

	defer resp.Body.Close()

	var versions Versions

	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, errors.WithMessage(err, "failed to decode node versions")
	}

	return versions, nil
}

// GetMatchVersion returns the first version that matches the provided semantic version constraint.
// It retrieves all available node versions and checks each one against the given constraint.
//
// Parameters:
//   - semverVersionConstraint: A string representing the semantic version constraint to match against.
//
// Returns:
//   - A pointer to a string containing the matching version if found, or nil if no match is found.
//   - An error if there was a failure in retrieving the node versions.
func GetMatchVersion(semverVersionConstraint string) (*string, error) {
	versions, err := GetAllVersions()

	if err != nil {
		return nil, errors.WithMessage(err, "failed to get node versions")
	}

	for _, version := range versions {
		isMatch, _ := version_match.VersionMatch(semverVersionConstraint, version.Version)

		if isMatch {
			return &version.Version, nil
		}
	}

	return nil, nil
}

// GetCurrentVersion retrieves the current version of Node.js installed on the system.
// It executes the command "node -v" and returns the version as a string pointer.
// If there is an error during the execution of the command, it returns an error with a descriptive message.
//
// Returns:
//   - A pointer to a string containing the Node.js version, or nil if an error occurred.
//   - An error if the command fails to execute or if there is an issue retrieving the version.
func GetCurrentVersion() (*string, error) {
	cmd := exec.Command("node", "-v")

	output, err := cmd.Output()

	if err != nil {
		return nil, errors.WithMessage(err, "failed to get current node version")
	}

	version := string(output)

	return &version, nil
}
