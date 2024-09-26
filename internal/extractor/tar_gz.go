package extractor

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func extractTarGz(tarGzFilePath, distFolder string) error {
	// 打开压缩文件
	file, err := os.Open(tarGzFilePath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	// 创建一个 gzip.Reader 从压缩文件中读取数据
	gz, err := gzip.NewReader(file)
	if err != nil {
		return errors.WithStack(err)
	}
	defer gz.Close()

	// 创建一个 tar.Reader 从 gzip.Reader 中读取数据
	tarReader := tar.NewReader(gz)

	// 逐个文件解压并写入目标目录
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.WithStack(err)
		}

		// 目标文件路径
		target := filepath.Join(distFolder, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// 如果是目录，创建目录
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return errors.WithStack(err)
			}
		case tar.TypeReg:
			// 如果是普通文件，创建并写入文件内容
			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return errors.WithStack(err)
			}
			defer file.Close()

			if _, err := io.Copy(file, tarReader); err != nil {
				return errors.WithStack(err)
			}
		case tar.TypeSymlink:
			// 如果是软链接，创建软链接
			if err := os.Symlink(header.Linkname, target); err != nil {
				return errors.WithStack(err)
			}
		case tar.TypeLink:
			// 如果是硬链接，创建硬链接
			if err := os.Link(filepath.Join(distFolder, header.Linkname), target); err != nil {
				return errors.WithStack(err)
			}
		default:
			return fmt.Errorf("Unsupported type: %v in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}
