package mirrors

import (
	"os"
	"testing"
)

func TestNodeMirror(t *testing.T) {
	// Backup original environment variables
	originalLang := os.Getenv("LANG")
	originalNodeMirror := os.Getenv("NODE_MIRROR")

	// Restore environment variables after test
	defer func() {
		os.Setenv("LANG", originalLang)
		os.Setenv("NODE_MIRROR", originalNodeMirror)
	}()

	tests := []struct {
		langEnv       string
		nodeMirrorEnv string
		expected      string
	}{
		{"en_US", "", "https://nodejs.org/dist/"},
		{"zh_CN", "", "https://registry.npmmirror.com/-/binary/node/"},
		{"en_US", "https://custom-mirror.com/node/", "https://custom-mirror.com/node/"},
		{"zh_CN", "https://custom-mirror.com/node/", "https://custom-mirror.com/node/"},
	}

	for _, test := range tests {
		os.Setenv("LANG", test.langEnv)
		os.Setenv("NODE_MIRROR", test.nodeMirrorEnv)

		nodeMirrorURL := getNodeMirror()

		if nodeMirrorURL != test.expected {
			t.Errorf("For LANG=%s and NODE_MIRROR=%s, expected %s but got %s", test.langEnv, test.nodeMirrorEnv, test.expected, nodeMirrorURL)
		}
	}
}
