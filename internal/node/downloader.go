package node

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

	"github.com/axetroy/virtual_node_env/internal/mirrors"
	"github.com/axetroy/virtual_node_env/internal/util"
	pb "github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
)

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
		return errors.New(fmt.Sprintf("download file '%s' with status code %d", url, resp.StatusCode))
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

	defer bar.Finish()

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
	// 打开压缩文件
	file, err := os.Open(tarGzFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建一个 gzip.Reader 从压缩文件中读取数据
	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
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
			return err
		}

		// 目标文件路径
		target := filepath.Join(extractToDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// 如果是目录，创建目录
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			// 如果是普通文件，创建并写入文件内容
			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
		case tar.TypeSymlink:
			// 如果是软链接，创建软链接
			if err := os.Symlink(header.Linkname, target); err != nil {
				return err
			}
		case tar.TypeLink:
			// 如果是硬链接，创建硬链接
			if err := os.Link(filepath.Join(extractToDir, header.Linkname), target); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Unsupported type: %v in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}

func GetArtifactName(version string) string {
	return getNodeFileName(version)
}

func Download(version string, dir string) (string, error) {
	fileNameWithoutExt := GetArtifactName(version)
	fileNameWithExt := getNodeDownloadName(version)

	util.Debug("fileNameWithoutExt: %s\n", fileNameWithoutExt)
	util.Debug("fileNameWithExt: %s\n", fileNameWithExt)

	if fileNameWithoutExt == "" || fileNameWithExt == "" {
		return "", errors.New("unsupported node version")
	}

	destFolder := filepath.Join(dir, "node", fileNameWithoutExt)

	// skip download if folder exists
	if _, err := os.Stat(destFolder); err == nil {
		return destFolder, nil
	}

	url := fmt.Sprintf("%sv%s/%s", mirrors.NODE_MIRROR, version, fileNameWithExt)

	util.Debug("downloadURL: %s\n", url)

	destFile := filepath.Join(dir, "download", fileNameWithExt)

	// download file
	if err := downloadFile(url, destFile); err != nil {
		return "", errors.WithStack(err)
	}

	// decompress file
	if err := decompressFile(destFile, filepath.Join(dir, "node")); err != nil {
		return "", nil
	}

	// remove downloaded file
	defer os.Remove(destFile)

	return destFolder, nil
}
