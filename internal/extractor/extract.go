package extractor

import (
	"path/filepath"
	"strings"

	"github.com/axetroy/nodapt/internal/util"
	"github.com/pkg/errors"
)

// Extract extracts the contents of a compressed file to a specified directory.
// It supports files with the ".7z" and ".tar.xz" extensions.
//
// Parameters:
//   - fileName: The path to the compressed file to be extracted.
//   - destFolder: The directory where the contents will be extracted.
//
// Returns:
//   - An error if the extraction fails or if the file format is unsupported.
func Extract(fileName string, destFolder string) error {
	name := filepath.Base(fileName)

	util.Debug("Extracting %s to %s\n", name, destFolder)

	if strings.HasSuffix(name, ".tar.xz") {
		return extractTarXz(fileName, destFolder)
	} else if strings.HasSuffix(name, ".7z") {
		return extract7Z(fileName, destFolder)
	} else {
		return errors.Errorf("unsupported file format: %s", name)
	}
}
