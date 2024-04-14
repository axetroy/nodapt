package virtualnodeenv

import (
	"fmt"
	"runtime"
)

func getNodeFileName(version string) string {
	if runtime.GOARCH == "amd64" {
		return fmt.Sprintf("node-v%s-win-x64", version)
	} else if runtime.GOARCH == "arm64" {
		return fmt.Sprintf("node-v%s-win-arm64", version)
	} else {
		return ""
	}
}

func getNodeDownloadName(version string) string {
	fileName := getNodeFileName(version)

	if runtime.GOARCH == "amd64" {
		return fmt.Sprintf("%s.zip", fileName)
	} else if runtime.GOARCH == "arm64" {
		return fmt.Sprintf("%s.zip", fileName)
	} else {
		return ""
	}
}
