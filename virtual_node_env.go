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

var virtualNodeEnvDir string

func init() {
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

func generateNewEnvs(paths []string) []string {
	envs := []string{}

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "PATH=") {
			envs = append(envs, env)
		}
	}

	var newPath string

	if runtime.GOOS == "windows" {
		newPath = strings.Join(paths, ";") + ";" + os.Getenv("PATH") + ";" + strings.Join(paths, ";")
	} else {
		newPath = strings.Join(paths, ":") + ":" + os.Getenv("PATH") + ":" + strings.Join(paths, ":")
	}

	debug("newPath: %s\n", newPath)

	return append(envs, "PATH="+newPath)
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

	files, err := os.ReadDir(binaryFileDir)

	if err != nil {
		return errors.WithStack(err)
	}

	// find the command in the binary file directory
	for _, nodeCommand := range files {
		exeName := nodeCommand.Name()

		if runtime.GOOS == "windows" {
			if strings.HasSuffix(exeName, ".cmd") || strings.HasSuffix(exeName, ".exe") {
				exeName = strings.TrimSuffix(exeName, ".cmd")
				exeName = strings.TrimSuffix(exeName, ".exe")

				if command == exeName {
					command = filepath.Join(binaryFileDir, exeName)
					debug("fallback to command: %s\n", command)
					break
				}
			}
		} else {
			if command == exeName {
				command = filepath.Join(binaryFileDir, exeName)
				debug("fallback to command: %s\n", command)
				break
			}
		}
	}

	if len(options.Cmd) == 1 {
		process = exec.Command(command)
	} else {
		process = exec.Command(command, options.Cmd[1:]...)
	}

	process.Env = generateNewEnvs([]string{binaryFileDir})
	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	if err := process.Run(); err != nil {

		exitCode := process.ProcessState.ExitCode()

		if exitCode != 0 {
			os.Exit(exitCode)
			return nil
		}

		return errors.WithStack(err)
	}

	return nil
}

func Use(version string) error {

	shell, err := getShell()

	if err != nil {
		return errors.WithStack(err)
	}

	debug("shell: %s\n", shell)

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

	cmd.Env = generateNewEnvs([]string{binaryFileDir})
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
		exitCode := cmd.ProcessState.ExitCode()

		if exitCode != 0 {
			os.Exit(exitCode)
			return nil
		}

		return errors.WithStack(err)
	}

	return nil
}

func Clean() error {
	if err := os.RemoveAll(virtualNodeEnvDir); err != nil {
		return errors.WithStack(err)
	}

	fmt.Fprintf(os.Stderr, "Cleaned up all node versions\n")

	return nil
}

func List() error {
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
		return nil, errors.WithStack(err)
	}

	c := &Config{}

	if err := json.Unmarshal(content, c); err != nil {
		return nil, errors.WithMessagef(err, "Read config file %s failed", filePath)
	}

	return c, nil
}
