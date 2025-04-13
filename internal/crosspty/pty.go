package crosspty

import (
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/aymanbagabas/go-pty"
	"github.com/pkg/errors"
	"golang.org/x/term"
)

func setPytSize(p pty.Pty) error {
	if w, h, err := getConsoleSize(p); err == nil {
		return p.Resize(w, h)
	}

	return nil
}

func Start(shellPath string, env map[string]string) error {
	ptmx, err := pty.New()
	if err != nil {
		return err
	}

	defer ptmx.Close()

	c := ptmx.Command(shellPath)
	if err := c.Start(); err != nil {
		return err
	}

	if err := setPytSize(ptmx); err != nil {
		return errors.WithMessage(err, "failed to set initial pty size")
	}

	ch := make(chan os.Signal, 1)

	go listenOnResize(ch, ptmx, setPytSize)

	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))

	if err != nil {
		return errors.WithMessage(err, "failed to set stdin in raw mode")
	}

	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	time.Sleep(1500 * time.Millisecond) // Give the shell some time to start.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	shellName := filepath.Base(shellPath)

	// Set environment variables
	for k, v := range env {
		switch shellName {
		case "bash", "zsh":
			_, _ = ptmx.Write([]byte("export " + k + "='" + v + "'\n"))
		case "fish":
			_, _ = ptmx.Write([]byte("set -gx " + k + " '" + v + "'\n"))
		case "powershell", "powershell.exe":
			_, _ = ptmx.Write([]byte("$env:" + k + "='" + v + "'\n"))
		case "cmd", "cmd.exe":
			_, _ = ptmx.Write([]byte("set " + k + "=" + v + "\r\n"))
		default:
			_, _ = ptmx.Write([]byte("export " + k + "='" + v + "'\n"))
		}
	}

	// Clear the screen
	switch shellName {
	case "bash", "zsh":
		_, _ = ptmx.Write([]byte("clear\n"))
	case "fish":
		_, _ = ptmx.Write([]byte("clear\n"))
	case "powershell", "powershell.exe":
		_, _ = ptmx.Write([]byte("Clear-Host\n"))
	case "cmd", "cmd.exe":
		_, _ = ptmx.Write([]byte("cls\r\n"))
	default:
		_, _ = ptmx.Write([]byte("clear\n"))
	}

	go func() {
		_, _ = io.Copy(ptmx, os.Stdin)
	}()

	go func() {
		_, _ = io.Copy(os.Stdout, ptmx)
	}()

	return c.Wait()
}
