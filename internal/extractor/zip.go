package extractor

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func extractZip(zipFilePath, distFolder string) error {
	// 打开 ZIP 文件
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer r.Close()

	// 创建解压目录
	err = os.MkdirAll(distFolder, os.ModePerm)
	if err != nil {
		return err
	}

	// 遍历 ZIP 文件中的文件并解压缩
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// 创建解压后的文件
		path := filepath.Join(distFolder, f.Name)
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(path, os.ModePerm)
		} else {
			_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)
			w, err := os.Create(path)
			if err != nil {
				return err
			}
			defer w.Close()

			// 复制文件内容
			_, err = io.Copy(w, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
