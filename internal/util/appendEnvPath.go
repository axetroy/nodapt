package util

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func removeFromSlice(s []string, r string) []string {
	result := s[:0] // Use the same slice to avoid extra allocations
	for _, v := range s {
		if v != r {
			result = append(result, v)
		}
	}
	return result
}

func removeNodePath(paths string) string {
	// Split the PATH into directories
	dirs := strings.Split(paths, string(os.PathListSeparator))

	// Remove the node binary directory from the PATH
	for _, dir := range dirs {
		stat, err := os.Stat(dir)

		if err != nil {
			continue
		}

		if !stat.IsDir() {
			continue
		}

		// Check if the directory contains the node binary
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

	fileLoop:
		for _, file := range files {
			isContainsNodeBinary := (runtime.GOOS == "windows" && file.Name() == "node.exe") || (runtime.GOOS != "windows" && file.Name() == "node")

			if isContainsNodeBinary {
				// Remove the directory from the PATH
				dirs = removeFromSlice(dirs, dir)
				break fileLoop
			}
		}
	}

	return strings.Join(dirs, string(os.PathListSeparator))
}

func AppendEnvPath(pathDir string) string {
	oldPath := removeNodePath(os.Getenv("PATH"))

	newPath := fmt.Sprintf("%s%c%s%c%s", pathDir, os.PathListSeparator, oldPath, os.PathListSeparator, pathDir)

	return newPath
}
