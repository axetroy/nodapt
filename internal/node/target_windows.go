package node

import (
	"fmt"
	"runtime"
)

func GetRemoteArtifactTarget(version string) *RemoteArtifactTarget {
	fileName := getNodeFileName(version)

	if fileName == nil {
		return nil
	}

	ext := ".7z"

	return &RemoteArtifactTarget{
		FileName: *fileName,
		FullName: fmt.Sprintf("%s%s", *fileName, ext),
		Ext:      ext,
	}
}

func getNodeFileName(version string) *string {
	if runtime.GOARCH == "amd64" {
		str := fmt.Sprintf("node-v%s-win-x64", version)
		return &str
	} else if runtime.GOARCH == "arm64" {
		str := fmt.Sprintf("node-v%s-win-arm64", version)
		return &str
	} else {
		return nil
	}
}
