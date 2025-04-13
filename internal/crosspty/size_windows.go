//go:build windows

package crosspty

import (
	"os"

	"github.com/aymanbagabas/go-pty"
)

func listenOnResize(ch chan os.Signal, p pty.Pty, onResize func(p pty.Pty) error) {
	// Windows does not support resizing pty, so we do nothing here.
	// This is a no-op function to satisfy the interface.
}
