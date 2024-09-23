package node

import (
	"encoding/json"
	"net/http"

	"github.com/axetroy/virtual_node_env/internal/mirrors"
	"github.com/axetroy/virtual_node_env/internal/version_match"
	"github.com/pkg/errors"
)

// GetAllVersions retrieves a list of all available Node.js versions from the official Node.js distribution index.
// It makes an HTTP GET request to the Node.js API endpoint and decodes the JSON response into a slice of version strings.
//
// Returns:
// - A slice of strings containing the Node.js versions, or
// - An error if the request fails or if there is an issue decoding the response.
func GetAllVersions() ([]string, error) {
	resp, err := http.Get(mirrors.NODE_MIRROR + "index.json")

	if err != nil {
		return nil, errors.WithMessage(err, "failed to get node versions")
	}

	defer resp.Body.Close()

	var versions []struct {
		Version string `json:"version"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, errors.WithMessage(err, "failed to decode node versions")
	}

	var result []string
	for _, v := range versions {
		result = append(result, v.Version)
	}

	return result, nil
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
		isMatch, _ := version_match.VersionMatch(semverVersionConstraint, version)

		if isMatch {
			return &version, nil
		}
	}

	return nil, nil
}
