package command

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/axetroy/nodapt/internal/crosspty"
	"github.com/axetroy/nodapt/internal/node"
	"github.com/axetroy/nodapt/internal/shell"
	"github.com/axetroy/nodapt/internal/util"
	"github.com/pkg/errors"
)

func Use(constraint *string) error {
	if constraint == nil {
		cwd, err := os.Getwd()

		if err != nil {
			return errors.WithStack(err)
		}

		packageJSONPath := util.LoopUpFile(cwd, "package.json")

		// If the package.json file is found, then use the node constraint in the package.json to run the command
		if packageJSONPath != nil {
			util.Debug("Use node constraint from %s\n", *packageJSONPath)

			c, err := node.GetConstraintFromPackageJSON(*packageJSONPath)

			if err != nil {
				return errors.WithMessagef(err, "failed to get node constraint from %s", *packageJSONPath)
			}

			constraint = c
		}
	}

	if constraint == nil {
		return errors.New("constraint is required")
	}

	util.Debug("Use constraint: %s\n", *constraint)

	version, err := node.GetMatchVersion(*constraint)

	if err != nil {
		return errors.WithStack(err)
	}

	if version == nil {
		return errors.Errorf("Cannot find the version of node which matches the constraint: %s", *constraint)
	}

	shellPath, err := shell.GetPath()
	if err != nil {
		return errors.WithMessage(err, "Cannot find shell")
	}

	util.Debug("Current shell: %s\n", shellPath)

	nodePath, err := node.Download(*version, nodapt_dir)
	if err != nil {
		return errors.WithStack(err)
	}

	var binaryFileDir string
	if runtime.GOOS == "windows" {
		binaryFileDir = nodePath
	} else {
		binaryFileDir = filepath.Join(nodePath, "bin")
	}

	oldPath := os.Getenv("PATH")
	oldNpmConfigPrefix := os.Getenv("NPM_CONFIG_PREFIX")
	defer os.Setenv("PATH", oldPath)                         // 确保在函数返回时恢复原始的 PATH
	defer os.Setenv("NPM_CONFIG_PREFIX", oldNpmConfigPrefix) // 确保在函数返回时恢复原始的 NPM_CONFIG_PREFIX

	// 设置新的 PATH 变量
	os.Setenv("PATH", util.AppendEnvPath(binaryFileDir))
	os.Setenv("NPM_CONFIG_PREFIX", nodePath)

	if err := crosspty.Start(shellPath, map[string]string{
		"NPM_CONFIG_PREFIX": os.Getenv("NPM_CONFIG_PREFIX"),
		"PATH":              os.Getenv("PATH"),
	}, fmt.Sprintf("Welcome to the nodapt shell, Currently using node %s!", *version)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
