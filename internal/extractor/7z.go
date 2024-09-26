package extractor

import (
	"io"
	"os"
	"path/filepath"

	"github.com/bodgit/sevenzip"
	"github.com/pkg/errors"
)

func extract7ZFile(f *sevenzip.File, destFolder string) error {
	rc, err := f.Open()
	if err != nil {
		return errors.WithStack(err)
	}
	defer rc.Close()

	// 创建解压后的文件
	path := filepath.Join(destFolder, f.Name)

	if f.FileInfo().IsDir() {
		_ = os.MkdirAll(path, os.ModePerm)
	} else {
		_ = os.MkdirAll(filepath.Dir(path), os.ModePerm)
		w, err := os.Create(path)
		if err != nil {
			return errors.WithStack(err)
		}
		defer w.Close()

		// 复制文件内容
		_, err = io.Copy(w, rc)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func extract7Z(zip7FilePath, destFolder string) error {
	// 打开 7z 文件
	r, err := sevenzip.OpenReader(zip7FilePath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer r.Close()

	// 创建解压目录
	err = os.MkdirAll(destFolder, os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}

	// 遍历 7z 文件中的文件并解压缩
	for _, f := range r.File {
		if err := extract7ZFile(f, destFolder); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
