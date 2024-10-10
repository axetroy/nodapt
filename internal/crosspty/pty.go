package crosspty

import (
	"io"
	"os"
	"os/signal"

	"github.com/aymanbagabas/go-pty"
	"github.com/pkg/errors"
	"golang.org/x/term"
)

func Start(shellPath string) error {
	ptmx, err := pty.New()
	if err != nil {
		return err
	}

	defer ptmx.Close()

	c := ptmx.Command(shellPath)
	if err := c.Start(); err != nil {
		return err
	}

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	notifySizeChanges(ch)
	go handlePtySize(ptmx, ch)
	initSizeChange(ch)
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		return errors.WithMessage(err, "failed to set stdin in raw mode")
	}

	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	go func() {
		_, _ = io.Copy(ptmx, os.Stdin)
	}()
	go func() {
		_, _ = io.Copy(os.Stdout, ptmx)
	}()

	return c.Wait()
}
