//go:build unix

package util

import "os"

func isExecutable(fileInfo os.FileInfo, _ string) bool {
	if fileInfo.IsDir() {
		return false
	}

	mode := fileInfo.Mode()

	const execPermissionBits = 0111 // Owner, group, and others execute permission bits

	return mode&execPermissionBits != 0 // Check if any execute bit is set (owner, group, or others)
}
