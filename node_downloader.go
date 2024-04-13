package virtualnodeenv

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
)

func DownloadNodeJs(version string) (string, error) {
	virtualNodeEnvPath := filepath.Join(os.Getenv("HOME"), ".vne")

	fileNameWithoutExt := getNodeFileName(version)
	fileNameWithExt := getNodeDownloadName(version)

	destFolder := filepath.Join(virtualNodeEnvPath, "node", fileNameWithoutExt)

	// skip download if folder exists
	if _, err := os.Stat(destFolder); err == nil {
		return destFolder, nil
	}

	// url := fmt.Sprintf("https://nodejs.org/dist/v%s/%s", version, targetName)
	url := fmt.Sprintf("https://registry.npmmirror.com/-/binary/node/v%s/%s", version, fileNameWithExt)

	destFile := filepath.Join(virtualNodeEnvPath, "versions", fileNameWithExt)

	// download file
	if err := downloadFile(url, destFile); err != nil {
		return "", errors.WithStack(err)
	}

	// decompress file
	if err := decompressFile(destFile, filepath.Join(virtualNodeEnvPath, "node")); err != nil {
		return "", nil
	}

	// remove downloaded file
	defer os.Remove(destFile)

	return destFolder, nil
}

func ensureDirExists(dirPath string) error {
	// 检查目录是否已存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 目录不存在，创建它
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		// 其他错误
		return err
	}

	return nil
}

func downloadFile(url string, fileName string) error {
	// 发起 HTTP GET 请求
	resp, err := http.Get(url)

	if err != nil {
		return errors.WithStack(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return errors.New(fmt.Sprintf("download file with status code %d", resp.StatusCode))
	}

	if err := ensureDirExists(filepath.Dir(fileName)); err != nil {
		return errors.WithStack(err)
	}

	// 创建文件
	file, err := os.Create(fileName)

	if err != nil {
		return errors.WithStack(err)
	}

	defer file.Close()

	tmpl := fmt.Sprintf(`{{string . "prefix"}}{{ "%s" }} {{counters . }} {{ bar . "[" "=" ">" "-" "]"}} {{percent . }} {{speed . }}{{string . "suffix"}}`, filepath.Base(fileName))

	bar := pb.ProgressBarTemplate(tmpl).Start64(resp.ContentLength)

	bar.SetWriter(os.Stdout)

	barReader := bar.NewProxyReader(resp.Body)

	// 将响应体内容复制到文件
	_, err = io.Copy(file, barReader)

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func decompressFile(fileName string, extractToDir string) error {
	name := filepath.Base(fileName)

	if strings.HasSuffix(name, ".zip") {
		return unzip(fileName, extractToDir)
	} else if strings.HasSuffix(name, ".tar.gz") {
		return extractTarGz(fileName, extractToDir)
	} else {
		return errors.New("unsupported file format")
	}
}

func unzip(zipFilePath, extractToDir string) error {
	// 打开 ZIP 文件
	r, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer r.Close()

	// 创建解压目录
	err = os.MkdirAll(extractToDir, os.ModePerm)
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
		path := filepath.Join(extractToDir, f.Name)
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

func extractTarGz(tarGzFilePath, extractToDir string) error {
	// 打开 .tar.gz 文件
	f, err := os.Open(tarGzFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// 创建 gzip 读取器
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	// 创建 tar 读取器
	tr := tar.NewReader(gz)

	// 创建硬链接映射
	linkMap := make(map[string]string)

	// 遍历 tar 文件中的文件并解压缩
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // 所有文件都已解压完毕
		}
		if err != nil {
			return err
		}

		// 构建解压后的文件路径
		targetPath := filepath.Join(extractToDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录
			err = os.MkdirAll(targetPath, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
		case tar.TypeReg:
			// 创建文件并写入内容
			w, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, tr)
			if err != nil {
				return err
			}
			w.Close()
			// 设置可执行文件权限
			if header.Mode&0111 != 0 { // 文件有可执行权限
				err = os.Chmod(targetPath, os.FileMode(header.Mode)|0111)
				if err != nil {
					return err
				}
			}
		case tar.TypeSymlink:
			// 创建软链接
			err = os.Symlink(header.Linkname, targetPath)
			if err != nil {
				return err
			}
			// 设置软链接权限
			linkMode := os.FileMode(header.Mode)

			if header.Mode&0111 != 0 { // 链接有可执行权限
				linkMode |= 0111
			}

			err = os.Chmod(targetPath, linkMode)
			if err != nil {
				return err
			}
		case tar.TypeLink:
			// 存储硬链接映射
			linkMap[targetPath] = header.Linkname
		}
	}

	// 创建硬链接
	for newLink, origLink := range linkMap {
		err := os.Link(origLink, newLink)
		if err != nil {
			return err
		}
	}

	return nil
}
