package node

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/axetroy/virtual_node_env/internal/extractor"
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

func Download(version string, dir string) (string, error) {
	// Remove the 'v' prefix from the version
	version = strings.TrimPrefix(version, "v")

	artifact := GetRemoteArtifactTarget(version)

	if artifact == nil {
		return "", errors.New("unsupported node version")
	}

	nodeFolder := filepath.Join(dir, "node")
	destFolder := filepath.Join(nodeFolder, artifact.FileName)

	// skip download if folder exists
	if _, err := os.Stat(destFolder); err == nil {
		// make sure the folder is not empty

		if files, err := os.ReadDir(destFolder); err == nil && len(files) > 0 {
			return destFolder, nil
		}
	}

	url := fmt.Sprintf("%sv%s/%s", NODE_MIRROR, version, artifact.FullName)

	util.Debug("downloadURL: %s\n", url)

	destFile := filepath.Join(dir, "download", artifact.FullName)

	// download file
	if err := downloadFile(url, destFile); err != nil {
		return "", errors.WithStack(err)
	}

	// remove downloaded file
	defer os.Remove(destFile)

	// decompress file
	if err := extractor.Extract(destFile, nodeFolder); err != nil {
		return "", nil
	}

	return destFolder, nil
}
