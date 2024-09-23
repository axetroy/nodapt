package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/axetroy/virtual_node_env/internal/node"
	"github.com/axetroy/virtual_node_env/internal/util"
	"github.com/pkg/errors"
)

type RunOptions struct {
	Version string   `json:"version"`
	Cmd     []string `json:"cmd"`
}

// Run executes a command using a specified version of Node.js.
// It downloads the Node.js version if it is not already available,
// sets the appropriate environment variables, and runs the command
// with the provided options.
//
// Parameters:
//   - options: A pointer to RunOptions containing the version of Node.js
//     to use and the command to execute. The version should be prefixed
//     with 'v' (e.g., "v14.17.0").
//
// Returns:
//   - An error if the command fails to execute or if there is an issue
//     downloading the specified Node.js version; otherwise, it returns nil.
func Run(options *RunOptions) error {
	nodeVersion := strings.TrimPrefix(options.Version, "v")

	nodeEnvPath, err := node.Download(nodeVersion, virtual_node_env_dir)

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

	os.Setenv("PATH", util.AppendEnvPath(binaryFileDir))

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

// RunWithConfig executes a command with the specified configuration.
//
// It loads the configuration from the provided configPath. If the configuration
// file cannot be loaded or if the required 'Node' field is missing, it returns
// an error. If successful, it runs the command with the options derived from
// the configuration.
//
// Parameters:
//   - configPath: The path to the configuration file.
//   - cmd: A slice of strings representing the command to be executed.
//
// Returns:
//   - An error if the configuration loading fails or if the 'Node' field is missing;
//     otherwise, it returns nil.
func RunWithConfig(configPath string, cmd []string) error {
	config, err := LoadConfig(configPath)

	if err != nil {
		return errors.WithMessage(err, "load configuration file failed")
	}

	if config.Node == "" {
		return errors.New("missing node field in the configuration file")
	}

	return Run(&RunOptions{
		Version: config.Node,
		Cmd:     cmd,
	})
}
