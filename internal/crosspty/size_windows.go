//go:build windows

package crosspty

import (
	"os"
	"time"

	"github.com/aymanbagabas/go-pty"
	"golang.org/x/sys/windows"
)

func listenOnResize(ch chan os.Signal, p pty.Pty, onResize func(p pty.Pty) error) {
	// Windows does not support resizing pty, so we do nothing here.
	// This is a no-op function to satisfy the interface.
	var prevCols, prevRows int

	for {
		time.Sleep(500 * time.Millisecond)

		cols, rows, err := getConsoleSize(p)

		if err != nil {
			continue
		}

		if cols != prevCols || rows != prevRows {
			err := onResize(p)

			if err == nil {
				prevCols = cols
				prevRows = rows
			}
		}
	}
}

func getConsoleSize(p pty.Pty) (int, int, error) {
	handle := windows.Handle(os.Stdout.Fd())
	var info windows.ConsoleScreenBufferInfo

	if err := windows.GetConsoleScreenBufferInfo(handle, &info); err != nil {
		return 0, 0, err
	}

	cols := int(info.Window.Right - info.Window.Left + 1)
	rows := int(info.Window.Bottom - info.Window.Top + 1)

	return cols, rows, nil
}
