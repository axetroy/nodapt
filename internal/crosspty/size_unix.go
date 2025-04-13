//go:build !windows

package crosspty

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aymanbagabas/go-pty"
)

func listenOnResize(ch chan os.Signal, p pty.Pty, onResize func(p pty.Pty) error) {
	signal.Notify(ch, syscall.SIGWINCH)

	go func() {
		for range ch {
			_ = onResize(p)
		}
	}()
}
