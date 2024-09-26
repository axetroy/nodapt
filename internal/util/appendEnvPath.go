package util

import (
	"fmt"
	"os"
)

func AppendEnvPath(pathDir string) string {
	oldPath := os.Getenv("PATH")

	newPath := fmt.Sprintf("%s%c%s", pathDir, os.PathListSeparator, oldPath)

	return newPath
}
