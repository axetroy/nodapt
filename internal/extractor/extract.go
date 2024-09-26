package extractor

import (
	"errors"
	"path/filepath"
	"strings"
)

// Extract extracts the contents of a compressed file to a specified directory.
// It supports files with the ".zip" and ".tar.gz" extensions.
//
// Parameters:
//   - fileName: The path to the compressed file to be extracted.
//   - distFolder: The directory where the contents will be extracted.
//
// Returns:
//   - An error if the extraction fails or if the file format is unsupported.
func Extract(fileName string, distFolder string) error {
	name := filepath.Base(fileName)

	if strings.HasSuffix(name, ".zip") {
		return extractZip(fileName, distFolder)
	} else if strings.HasSuffix(name, ".tar.gz") {
		return extractTarGz(fileName, distFolder)
	} else if strings.HasSuffix(name, ".tar.xz") {
		return extractTarXz(fileName, distFolder)
	} else {
		return errors.New("unsupported file format")
	}
}
