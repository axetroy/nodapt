package mirrors

import (
	"os"
	"strings"

	"github.com/axetroy/virtual_node_env/internal/util"
)

var NODE_MIRROR string = "https://nodejs.org/dist/"

func init() {
	NODE_MIRROR = getNodeMirror()

	util.Debug("nodeMirrorURL: %s\n", NODE_MIRROR)
}

func getNodeMirror() string {
	var mirrorUrl = "https://nodejs.org/dist/"

	lang := strings.ToLower(os.Getenv("LANG"))

	if strings.Contains(lang, "zh_cn") {
		mirrorUrl = "https://registry.npmmirror.com/-/binary/node/"
	}

	return util.GetEnvsWithFallback(mirrorUrl, "NODE_MIRROR")
}
