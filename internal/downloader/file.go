package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/axetroy/virtual_node_env/internal/util"
	pb "github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
)

func DownloadFile(url string, dest string) error {
	// Perform HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return errors.WithStack(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return errors.Errorf("download file '%s' with status code %d", url, resp.StatusCode)
	}

	if err := util.EnsureDir(filepath.Dir(dest)); err != nil {
		return errors.WithStack(err)
	}

	// Create the destination file
	file, err := os.Create(dest)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	// Set up progress bar
	tmpl := fmt.Sprintf(`{{string . "prefix"}}{{ "%s" }} {{counters . }} {{ bar . "[" "=" ">" "-" "]"}} {{percent . }} {{speed . }}{{string . "suffix"}}`, filepath.Base(dest))
	bar := pb.ProgressBarTemplate(tmpl).Start64(resp.ContentLength)
	bar.SetWriter(os.Stdout)
	defer bar.Finish()

	// Use proxy reader for progress bar
	barReader := bar.NewProxyReader(resp.Body)

	// Copy the response body to the file
	if _, err := io.Copy(file, barReader); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
