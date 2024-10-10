package command

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

func Clean() error {
	if err := os.RemoveAll(nodapt_dir); err != nil {
		return errors.WithStack(err)
	}

	fmt.Fprintf(os.Stderr, "Cleaned up all node versions\n")

	return nil
}
