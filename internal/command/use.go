package command

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/axetroy/nodapt/internal/crosspty"
	"github.com/axetroy/nodapt/internal/node"
	"github.com/axetroy/nodapt/internal/shell"
	"github.com/axetroy/nodapt/internal/util"
	"github.com/pkg/errors"
)

func Use(constraint string) error {
	util.Debug("Use constraint: %s\n", constraint)

	version, err := node.GetMatchVersion(constraint)

	if err != nil {
		return errors.WithStack(err)
	}

	if version == nil {
		return errors.Errorf("Cannot find the version of node which matches the constraint: %s", constraint)
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
	defer os.Setenv("PATH", oldPath) // 确保在函数返回时恢复原始的 PATH

	// 设置新的 PATH 变量
	os.Setenv("PATH", util.AppendEnvPath(binaryFileDir))

	if err := crosspty.Start(shellPath); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
