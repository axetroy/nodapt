//go:build windows

package crosspty

import (
	"os"

	"github.com/aymanbagabas/go-pty"
)

func notifySizeChanges(chan os.Signal) {}

func handlePtySize(p pty.Pty, _ chan os.Signal) {
	// TODO
}

func initSizeChange(chan os.Signal) {}
