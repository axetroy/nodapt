//go:build !windows

package crosspty

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/aymanbagabas/go-pty"
	"golang.org/x/term"
)

func getConsoleSize(p pty.Pty) (int, int, error) {
	w, h, err := term.GetSize(int(os.Stdin.Fd()))

	if err != nil {
		return 0, 0, err
	}

	return w, h, nil
}

func listenOnResize(ch chan os.Signal, p pty.Pty, onResize func(p pty.Pty) error) {
	signal.Notify(ch, syscall.SIGWINCH)

	for range ch {
		_ = onResize(p)
	}
}
