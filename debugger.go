package virtualnodeenv

import (
	"fmt"
	"os"
)

func debug(format string, a ...any) {
	if getEnvsWithFallback("", "DEBUG") == "1" {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}
