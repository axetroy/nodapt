package node

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/axetroy/nodapt/internal/downloader"
	"github.com/axetroy/nodapt/internal/extractor"
	"github.com/axetroy/nodapt/internal/util"
	"github.com/pkg/errors"
)

func Download(version string, dir string) (string, error) {
	// Remove the 'v' prefix from the version string
	version = strings.TrimPrefix(version, "v")

	// Get the artifact information for the given version
	artifact := GetRemoteArtifactTarget(version)
	if artifact == nil {
		return "", errors.New("unsupported node version")
	}

	extractFolder := filepath.Join(dir, "node", artifact.FileName)

	// Skip download if the folder already exists and contains files
	if _, err := os.Stat(extractFolder); err == nil {
		if files, err := os.ReadDir(extractFolder); err == nil && len(files) > 0 {
			return extractFolder, nil
		}
	}

	url := fmt.Sprintf("%sv%s/%s", NODE_MIRROR, version, artifact.FullName)
	util.Debug("downloadURL: %s\n", url)

	destFile := filepath.Join(dir, "download", artifact.FullName)

	// Download the file
	if err := downloader.DownloadFile(url, destFile); err != nil {
		return "", errors.WithStack(err)
	}

	// Decompress the file into the node folder
	if err := extractor.Extract(destFile, filepath.Dir(extractFolder)); err != nil {
		// If extraction fails, the downloaded file remains for debugging
		return "", errors.WithStack(err)
	}

	// Remove the downloaded file after successful extraction
	if err := os.Remove(destFile); err != nil {
		// Log warning but don't fail - extraction was successful
		util.Debug("Warning: failed to remove temporary file %s: %v\n", destFile, err)
	}

	return extractFolder, nil
}
