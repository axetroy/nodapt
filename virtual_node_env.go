package virtualnodeenv

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

type Options struct {
	Version string   `json:"version"`
	Cmd     []string `json:"cmd"`
}

func generateNewEnvs(paths []string) []string {
	envs := []string{}

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, "PATH=") {
			envs = append(envs, env)
		}
	}

	newPath := strings.Join(paths, ":") + ":" + os.Getenv("PATH")

	return append(envs, "PATH="+newPath)
}

func Setup(options *Options) error {
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
					break
				}
			}
		} else {
			if command == exeName {
				command = filepath.Join(binaryFileDir, exeName)
				break
			}
		}
	}

	if len(options.Cmd) == 1 {
		process = exec.Command(command)
	} else {
		process = exec.Command(command, options.Cmd[1:]...)
	}

	process.Env = generateNewEnvs([]string{filepath.Join(nodeEnvPath, "bin")})
	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	if err := process.Run(); err != nil {

		exitCode := process.ProcessState.ExitCode()

		if exitCode != 0 {
			os.Exit(exitCode)
		}

		return errors.WithStack(err)
	}

	return nil
}
