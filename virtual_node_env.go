package virtualnodeenv

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

type Options struct {
	Version string   `json:"version"`
	Cmd     []string `json:"cmd"`
}

var envPathDelimiter = ":"
var virtualNodeEnvDir string

func init() {
	if runtime.GOOS == "windows" {
		envPathDelimiter = ";"
	}

	virtualNodeEnvDirFromEnv := getEnvsWithFallback("", "NODE_ENV_DIR")

	if virtualNodeEnvDirFromEnv != "" {
		virtualNodeEnvDir = virtualNodeEnvDirFromEnv
		return
	}

	homeDir, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	virtualNodeEnvDir = filepath.Join(homeDir, ".virtual-node-env")
}

func Run(options *Options) error {
	nodeEnvPath, err := DownloadNodeJs(options.Version)

	if err != nil {
		return errors.WithStack(err)
	}

	var binaryFileDir string

	if runtime.GOOS == "windows" {
		binaryFileDir = nodeEnvPath
	} else {
		binaryFileDir = filepath.Join(nodeEnvPath, "bin")
	}

	var process *exec.Cmd

	command := options.Cmd[0]

	path := os.Getenv("PATH")

	newPath := binaryFileDir + envPathDelimiter + path

	os.Setenv("PATH", newPath)

	if len(options.Cmd) == 1 {
		process = exec.Command(command)
	} else {
		process = exec.Command(command, options.Cmd[1:]...)
	}

	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	if err := process.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Use(version string) error {

	shell, err := getShell()

	if err != nil {
		return errors.WithStack(err)
	}

	Debug("shell: %s\n", shell)

	nodeEnvPath, err := DownloadNodeJs(version)

	if err != nil {
		return errors.WithStack(err)
	}

	// 创建一个 *exec.Cmd 对象来运行 Fish shell
	cmd := exec.Command(shell)

	var binaryFileDir string

	if runtime.GOOS == "windows" {
		binaryFileDir = nodeEnvPath
	} else {
		binaryFileDir = filepath.Join(nodeEnvPath, "bin")
	}

	// 获取当前的 PATH 变量
	path := os.Getenv("PATH")

	newPath := binaryFileDir + envPathDelimiter + path
	// 设置新的 PATH 变量
	os.Setenv("PATH", newPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// 启动命令
	if err := cmd.Start(); err != nil {
		return errors.WithStack(err)
	}

	// Write to the stdin of the shell and ignore error
	_, _ = fmt.Fprintf(os.Stdin, "Now you are using node v%s\n", version)

	if err := cmd.Wait(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Remove(version string) error {
	target := getNodeFileName(version)

	dest := filepath.Join(virtualNodeEnvDir, "node", target)

	// 检查文件是否存在
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Node version %s not found\n", version)
		return nil
	}

	return os.RemoveAll(dest)
}

func Clean() error {
	if err := os.RemoveAll(virtualNodeEnvDir); err != nil {
		return errors.WithStack(err)
	}

	fmt.Fprintf(os.Stderr, "Cleaned up all node versions\n")

	return nil
}

func List() error {
	if _, err := os.Stat(filepath.Join(virtualNodeEnvDir, "node")); os.IsNotExist(err) {
		return nil
	}

	files, err := os.ReadDir(filepath.Join(virtualNodeEnvDir, "node"))

	if err != nil {
		return errors.WithStack(err)
	}

	for _, file := range files {
		fName := file.Name()
		if file.IsDir() && strings.HasPrefix(fName, "node-v") {
			n := strings.SplitN(fName, "-", -1)

			version := n[1]
			println(version)
		}
	}

	return nil
}

func ListRemote() error {
	// Request the HTML page.
	res, err := http.Get("https://nodejs.org/dist/")

	if err != nil {
		return errors.WithStack(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return errors.WithStack(err)
	}

	vs := []*version.Version{}

	firstStableVersion, _ := version.NewVersion("v1.0.0")

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimRight(s.Text(), "/")

		if strings.HasPrefix(text, "v") {
			ver, err := version.NewVersion(text)

			if err != nil {
				return
			}

			if ver.GreaterThanOrEqual(firstStableVersion) {
				vs = append(vs, ver)
			}
		}
	})

	// 对版本号进行排序
	sort.Sort(version.Collection(vs))

	for _, v := range vs {
		println(v.String())
	}

	return nil
}

type Config struct {
	Node string `json:"node"`
}

func LoadConfig(filePath string) (*Config, error) {
	content, err := os.ReadFile(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.WithMessage(err, "Can not found `.virtual-node-env.json` file.")
		}

		return nil, errors.WithStack(err)
	}

	c := &Config{}

	if err := json.Unmarshal(content, c); err != nil {
		return nil, errors.WithMessagef(err, "Read config file %s failed", filePath)
	}

	return c, nil
}
