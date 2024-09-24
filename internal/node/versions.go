package node

import (
	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/axetroy/virtual_node_env/internal/mirrors"
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
	resp, err := http.Get(mirrors.NODE_MIRROR + "index.json")

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

// GetLatestLTSVersion retrieves the latest Long Term Support (LTS) version
// from the available versions of the software. It returns the version as a
// string and an error if any issues occur while fetching the versions.
//
// Returns:
//   - A string representing the latest LTS version, or an empty string if none is found.
//   - An error if there was a problem retrieving the versions.
func GetLatestLTSVersion() (string, error) {
	var latestLTSVersion = ""

	versions, err := GetAllVersions()

	if err != nil {
		return "", errors.WithStack(err)
	}

	for _, version := range versions {
		// Check if the LTS field is a string and not empty
		if str, ok := version.LTS.(string); ok && str != "" {
			latestLTSVersion = version.Version
			break
		}

		// Check if the LTS field is a boolean and set to true
		if lsLTS, ok := version.LTS.(bool); ok && lsLTS {
			latestLTSVersion = version.Version
			break
		}
	}

	return latestLTSVersion, nil
}
