package util

import "os"

func AppendEnvPath(pathDir string) string {
	oldPath := os.Getenv("PATH")

	newPath := pathDir + string(os.PathListSeparator) + oldPath

	return newPath
}
