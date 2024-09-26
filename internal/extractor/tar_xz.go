package extractor

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/ulikunitz/xz"
)

// extractTarXzFile extracts a single file from the tar archive.
func extractTarXzFile(reader *tar.Reader, header *tar.Header, destFolder string) error {
	// Resolve the destination path.
	destPath := filepath.Join(destFolder, header.Name)

	// Ensure no file path traversal attacks by sanitizing the path.
	if !strings.HasPrefix(destPath, filepath.Clean(destFolder)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", header.Name)
	}

	switch header.Typeflag {
	case tar.TypeDir:
		// If it's a directory, create the directory.
		if err := os.MkdirAll(destPath, os.FileMode(header.Mode)); err != nil {
			return errors.WithStack(err)
		}
	case tar.TypeReg:
		// If it's a regular file, create and write the file content.
		file, err := os.OpenFile(destPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(header.Mode))
		if err != nil {
			return errors.WithStack(err)
		}
		defer file.Close()

		if _, err := io.Copy(file, reader); err != nil {
			return errors.WithStack(err)
		}
	case tar.TypeSymlink:
		// If it's a symbolic link, create a symlink.
		if err := os.Symlink(header.Linkname, destPath); err != nil {
			return errors.WithStack(err)
		}
	case tar.TypeLink:
		// If it's a hard link, create a hard link.
		linkTarget := filepath.Join(destFolder, header.Linkname)
		if err := os.Link(linkTarget, destPath); err != nil {
			return errors.WithStack(err)
		}
	// case tar.TypeChar, tar.TypeBlock, tar.TypeFifo:
	// 	// Handle special file types like character devices, block devices, and FIFOs.
	// 	if err := createSpecialFile(destPath, os.FileMode(header.Mode), header.Devmajor, header.Devminor); err != nil {
	// 		return errors.WithStack(err)
	// 	}
	default:
		return fmt.Errorf("unsupported file type: %v in %s", header.Typeflag, header.Name)
	}

	// If the file is executable, ensure proper permissions.
	if header.FileInfo().Mode()&0111 != 0 {
		if err := os.Chmod(destPath, os.FileMode(header.Mode)); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// extractTarXz extracts a .tar.xz archive into the specified destination folder.
func extractTarXz(tarXzFilePath, destFolder string) error {
	// Open the .tar.xz file.
	file, err := os.Open(tarXzFilePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", tarXzFilePath, err)
	}
	defer file.Close()

	// Create xz.Reader to decompress the .xz data.
	xzReader, err := xz.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create xz reader: %w", err)
	}

	// Create tar.Reader to read the decompressed .tar data.
	tarReader := tar.NewReader(xzReader)

	// Iterate over the files and directories in the .tar archive.
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break // End of archive.
		}

		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Extract each file.
		if err := extractTarXzFile(tarReader, header, destFolder); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
