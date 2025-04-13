package node

import (
	"fmt"
	"runtime"

	"github.com/Masterminds/semver/v3"
)

func GetRemoteArtifactTarget(version string) *RemoteArtifactTarget {
	fileName := getNodeFileName(version)

	if fileName == nil {
		return nil
	}

	ext := ".tar.xz"

	return &RemoteArtifactTarget{
		FileName: *fileName,
		FullName: fmt.Sprintf("%s%s", *fileName, ext),
		Ext:      ext,
	}
}

func getNodeFileName(version string) *string {
	if runtime.GOARCH == "amd64" {
		str := fmt.Sprintf("node-v%s-darwin-x64", version)

		return &str
	} else if runtime.GOARCH == "arm64" {
		// Node.js 16.0.0 and later versions have official support for Apple Silicon
		// https://nodejs.org/en/blog/release/v16.0.0/
		if c, err := semver.NewConstraint("< 16.0.0"); err == nil && c.Check(semver.MustParse(version)) {
			str := fmt.Sprintf("node-v%s-darwin-x64", version)
			return &str
		}

		str := fmt.Sprintf("node-v%s-darwin-arm64", version)

		return &str
	} else {
		return nil
	}
}
