package util

import (
	"os"
	"path/filepath"
)

// 检查文件是否存在
func IsFileExist(filePath string) bool {
	_, err := os.Stat(filePath)

	return err == nil
}

// LoopUpFile searches for a file with the specified name starting from the given root directory.
// It traverses up the directory tree until it finds the file or reaches the root directory.
//
// Parameters:
//   - root: The starting directory path from which to begin the search.
//   - fileName: The name of the file to look for.
//
// Returns:
//   - A pointer to the string containing the full path of the found file, or nil if the file does not exist.
func LoopUpFile(root string, fileName string) *string {
	for {
		Debug("Look up file in '%s'\n", root)

		configFilePath := filepath.Join(root, fileName)

		if IsFileExist(configFilePath) {
			return &configFilePath
		}

		// 获取上一级目录
		parentDir := filepath.Dir(root)

		if parentDir == root {
			break // 已经到达根目录
		}

		root = parentDir
	}

	return nil
}
