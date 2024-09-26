package extractor

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/pkg/errors"
)

func extract7ZFile(f *sevenzip.File, destFolder string) error {
	rc, err := f.Open()
	if err != nil {
		return errors.WithStack(err)
	}
	defer rc.Close()

	// 构建解压后的文件路径
	path := filepath.Join(destFolder, f.Name)

	// 确保路径安全，防止路径遍历攻击
	if !strings.HasPrefix(path, filepath.Clean(destFolder)+string(os.PathSeparator)) {
		return errors.Errorf("invalid file path: %s", f.Name)
	}

	// 处理文件类型
	switch {
	case f.FileInfo().IsDir():
		// 创建目录
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return errors.WithStack(err)
		}
	case (f.Mode() & os.ModeSymlink) != 0:
		// 如果是符号链接，创建软链接
		linkname, err := io.ReadAll(rc)
		if err != nil {
			return errors.WithStack(err)
		}

		if err := os.Symlink(string(linkname), path); err != nil {
			return errors.WithStack(err)
		}
	// case (f.Mode() & os.ModeNamedPipe) != 0:
	// 	// 如果是 FIFO (命名管道)
	// 	if err := createNamedPipe(path, f.Mode()); err != nil {
	// 		return errors.WithStack(err)
	// 	}
	case (f.Mode() & os.ModeDevice) != 0:
		// 如果是设备文件 (字符设备/块设备)
		// if err := createSpecialFile(path, f.Mode(), f.Major(), f.Minor()); err != nil {
		// 	return errors.WithStack(err)
		// }
	default:
		// 创建普通文件并写入内容
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return errors.WithStack(err)
		}

		w, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, f.Mode().Perm())
		if err != nil {
			return errors.WithStack(err)
		}
		defer w.Close()

		// 复制文件内容
		if _, err := io.Copy(w, rc); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

// extract7Z 解压 7z 文件到指定目录
func extract7Z(zip7FilePath, destFolder string) error {
	// 打开 7z 文件
	r, err := sevenzip.OpenReader(zip7FilePath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer r.Close()

	// 创建解压目录
	if err := os.MkdirAll(destFolder, os.ModePerm); err != nil {
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
