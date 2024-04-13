package virtualnodeenv

import (
	"fmt"
	"runtime"
)

func getNodeFileName(version string) string {
	if runtime.GOARCH == "x64" {
		return fmt.Sprintf("node-v%s-darwin-x64", version)
	} else if runtime.GOARCH == "arm64" {
		return fmt.Sprintf("node-v%s-darwin-arm64", version)
	} else {
		return ""
	}
}

func getNodeDownloadName(version string) string {
	fileName := getNodeFileName(version)

	if runtime.GOARCH == "x64" {
		return fmt.Sprintf("%s.tar.gz", fileName)
	} else if runtime.GOARCH == "arm64" {
		return fmt.Sprintf("%s.tar.gz", fileName)
	} else {
		return ""
	}
}
