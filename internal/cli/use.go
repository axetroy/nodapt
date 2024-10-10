package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

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
		return errors.Errorf("Can not found the version of node which match the constraint: %s", constraint)
	}

	shell, err := shell.GetPath()

	if err != nil {
		return errors.WithMessage(err, "Can not found shell")
	}

	util.Debug("Current shell: %s\n", shell)

	nodePath, err := node.Download(*version, nodapt_dir)

	if err != nil {
		return errors.WithStack(err)
	}

	// 创建一个 *exec.Cmd 对象来运行 Fish shell
	cmd := exec.Command(shell)

	var binaryFileDir string

	if runtime.GOOS == "windows" {
		binaryFileDir = nodePath
	} else {
		binaryFileDir = filepath.Join(nodePath, "bin")
	}

	oldPath := os.Getenv("PATH")

	defer os.Setenv("PATH", oldPath)

	// 设置新的 PATH 变量
	os.Setenv("PATH", util.AppendEnvPath(binaryFileDir))

	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// 启动命令
	if err := cmd.Start(); err != nil {
		return errors.WithStack(err)
	}

	// Write to the stdin of the shell and ignore error
	_, _ = fmt.Fprintf(os.Stdin, "Now you are using node %s\n", *version)

	if err := cmd.Wait(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
