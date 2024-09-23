package util

import (
	"fmt"
	"os"
)

// Debug prints formatted debug messages to standard error if the DEBUG environment variable is set to "1".
//
// Parameters:
//   - format: A format string that specifies how to format the output.
//   - a: A variadic parameter that allows passing any number of arguments to be formatted according to the format string.
//
// Usage:
//
//	Debug("This is a debug message: %s", "some value")
func Debug(format string, a ...any) {
	if GetEnvsWithFallback("", "DEBUG") == "1" {
		fmt.Fprintf(os.Stderr, format, a...)
	}
}
