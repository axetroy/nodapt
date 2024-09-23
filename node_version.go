package virtualnodeenv

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// GetAllNodeVersions returns all available node versions
func GetAllNodeVersions() ([]string, error) {
	resp, err := http.Get("https://nodejs.org/dist/index.json")

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
