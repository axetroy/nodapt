package node

import (
	"github.com/axetroy/virtual_node_env/internal/util"
)

var NODE_MIRROR string = getNodeMirror("https://nodejs.org/dist/")

func init() {
	util.Debug("nodeMirrorURL: %s\n", NODE_MIRROR)
}

func getNodeMirror(defaultMirror string) string {
	var mirrorUrl = defaultMirror

	if util.IsSimplifiedChinese() {
		mirrorUrl = "https://registry.npmmirror.com/-/binary/node/"
	}

	return util.GetEnvsWithFallback(mirrorUrl, "NODE_MIRROR")
}
