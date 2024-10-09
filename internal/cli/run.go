package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"

	"github.com/axetroy/nodapt/internal/node"
	"github.com/axetroy/nodapt/internal/util"
	"github.com/axetroy/nodapt/internal/version_constraint"
	"github.com/pkg/errors"
)

type RunOptions struct {
	Version string   `json:"version"` // The version of Node.js to use
	Cmd     []string `json:"cmd"`     // The command to execute
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
func run(options *RunOptions) error {
	if len(options.Cmd) == 0 {
		return errors.New("no command provided")
	}

	util.Debug("Run command: %s with node %s.\n", options.Cmd, options.Version)

	nodeEnvPath, err := node.Download(options.Version, virtual_node_env_dir)

	if err != nil {
		return errors.WithStack(err)
	}

	var binaryFileDir string

	if runtime.GOOS == "windows" {
		binaryFileDir = nodeEnvPath
	} else {
		binaryFileDir = filepath.Join(nodeEnvPath, "bin")
	}

	// Check if the node executable exists
	if ok, err := util.FindExecutable(binaryFileDir, "node"); err != nil {
		return errors.WithStack(err)
	} else if !ok {
		return errors.Errorf("node executable not found in %s, You should try to remove it.", binaryFileDir)
	}

	var process *exec.Cmd

	command := options.Cmd[0]

	oldPath := os.Getenv("PATH")

	defer os.Setenv("PATH", oldPath)

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

// RunWithConstraint executes a command with a specified version constraint.
//
// It takes a semantic version constraint as a string and a slice of strings representing
// the command to be executed. The function attempts to find a matching version based on
// the provided constraint. If a matching version is found, it runs the command with that
// version. If no matching version is found or if an error occurs while retrieving the
// version, it returns an error.
//
// Parameters:
//   - constraint: A string representing the semantic version constraint.
//   - command: A slice of strings representing the command to be executed.he
//
// Returns:
//   - error: Returns an error if the version cannot be matched or if the command fails to execute.md[1:]...)
func RunWithConstraint(constraint string, command []string) error {
	installedVersion := node.GetCurrentVersion()

	// If the node version is installed and the version satisfies the constraint, then run the command directly
	if installedVersion != nil {
		util.Debug("Current node version: %s\n", *installedVersion)
		if ok, err := version_constraint.Match(constraint, *installedVersion); err != nil {
			return errors.WithStack(err)
		} else if ok {
			util.Debug("Current node version %s is match the constraint, run command directly.\n", *installedVersion)
			return RunDirectly(command)
		}
	}

	// Found cached node version
	if cachedNodes, err := node.GetCachedVersions(virtual_node_env_dir); err != nil {
		return errors.WithStack(err)
	} else {
		// Sort versions in descending order
		sort.Sort(sort.Reverse(node.ByVersion(cachedNodes)))

		for _, node := range cachedNodes {
			if ok, err := version_constraint.Match(constraint, node.Version); err != nil {
				return errors.WithStack(err)
			} else if ok {
				util.Debug("Found cached node version %s is match the constraint.\n", node.Version)
				// Found the match version
				return run(&RunOptions{
					Version: node.Version,
					Cmd:     command,
				})
			}
		}
	}

	matchVersion, err := node.GetMatchVersion(constraint)

	if err != nil {
		return errors.WithMessage(err, "failed to get match version")
	}

	if matchVersion == nil {
		return errors.Errorf("no match version found for %s", constraint)
	}

	return run(&RunOptions{
		Version: *matchVersion,
		Cmd:     command,
	})
}

// RunDirectly executes a command specified by the cmd slice.
// The first element of cmd is the command to run, and the subsequent
// elements are the arguments to that command. The function sets the
// standard input, output, and error streams of the process to the
// corresponding streams of the calling process. If the command
// execution fails, it returns an error wrapped with stack trace
// information. If successful, it returns nil.
//
// Parameters:
//
//	cmd []string: A slice of strings where the first element is
//	the command to execute and the rest are its arguments.
//
// Returns:
//
//	error: An error if the command fails to execute, otherwise nil.
func RunDirectly(cmd []string) error {
	if len(cmd) == 0 {
		return errors.New("no command provided")
	}

	process := exec.Command(cmd[0], cmd[1:]...)

	process.Stdin = os.Stdin
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr

	if err := process.Run(); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("failed to run command: %s", cmd))
	}

	return nil
}
