package node

import (
	"encoding/json"
	"io"
	"os"
)

type PackageJSONEngine struct {
	Node *string `json:"node"`
}

type PackageJSON struct {
	Engines *PackageJSONEngine `json:"engines"`
}

func readPackageJSON(path string) (PackageJSON, error) {
	file, err := os.Open(path)

	if err != nil {
		return PackageJSON{}, err
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)

	if err != nil {
		return PackageJSON{}, err
	}

	var packageJSON PackageJSON

	if err := json.Unmarshal(bytes, &packageJSON); err != nil {
		return PackageJSON{}, err
	}

	return packageJSON, nil
}

// GetConstraintFromPackageJSON retrieves the Node.js version specified in the
// "engines" field of a package.json file located at the given path.
//
// Parameters:
//   - packageJSONPath: A string representing the file path to the package.json.
//
// Returns:
//   - A pointer to a string containing the Node.js version if found, or nil
//     if the version is not specified.
//   - An error if there was an issue reading the package.json file.
func GetConstraintFromPackageJSON(packageJSONPath string) (*string, error) {
	packageJSON, err := readPackageJSON(packageJSONPath)

	if err != nil {
		return nil, err
	}

	if packageJSON.Engines != nil && packageJSON.Engines.Node != nil {
		return packageJSON.Engines.Node, nil
	}

	return nil, nil
}
