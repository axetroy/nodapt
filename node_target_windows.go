package virtualnodeenv

import (
	"fmt"
	"runtime"
)

func getNodeFileName(version string) string {
	if runtime.GOARCH == "x86" {
		return fmt.Sprintf("node-v%s-win-x86", version)
	} else if runtime.GOARCH == "x64" {
		return fmt.Sprintf("node-v%s-win-x64", version)
	} else {
		return ""
	}
}

func getNodeDownloadName(version string) string {
	fileName := getNodeFileName(version)

	if runtime.GOARCH == "x86" {
		return fmt.Sprintf("%s.zip", fileName)
	} else if runtime.GOARCH == "x64" {
		return fmt.Sprintf("%s.zip", fileName)
	} else {
		return ""
	}
}
