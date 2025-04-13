package crosspty

import (
	"fmt"
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

func getNewLine(shellName string) string {
	switch shellName {
	case "cmd", "cmd.exe", "powershell", "powershell.exe":
		return "\r\n"
	default:
		return "\n"
	}
}

func Start(shellPath string, env map[string]string, welcome string) error {
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

	newLine := getNewLine(shellName)

	// Set environment variables
	for k, v := range env {
		switch shellName {
		case "bash", "zsh":
			_, _ = ptmx.Write([]byte(fmt.Sprintf("export %s='%s'", k, v) + newLine))
		case "fish":
			_, _ = ptmx.Write([]byte(fmt.Sprintf("set -gx %s '%s'", k, v) + newLine))
		case "powershell", "powershell.exe":
			_, _ = ptmx.Write([]byte(fmt.Sprintf("$env:%s='%s'", k, v) + newLine))
		case "cmd", "cmd.exe":
			_, _ = ptmx.Write([]byte("set " + k + "=" + v + newLine))
		default:
			_, _ = ptmx.Write([]byte(fmt.Sprintf("export %s='%s'", k, v) + newLine))
		}
	}

	// Clear the screen
	switch shellName {
	case "bash", "zsh":
		_, _ = ptmx.Write([]byte("clear" + newLine))
	case "fish":
		_, _ = ptmx.Write([]byte("clear" + newLine))
	case "powershell", "powershell.exe":
		_, _ = ptmx.Write([]byte("Clear-Host" + newLine))
	case "cmd", "cmd.exe":
		_, _ = ptmx.Write([]byte("cls" + newLine))
	default:
		_, _ = ptmx.Write([]byte("clear" + newLine))
	}

	// 清空掉 ptmx 的 stdout 缓冲区，使用 channel 并限制最多读取 500 毫秒
	buf := make([]byte, 1024)
	readCh := make(chan struct{})
	stopCh := make(chan struct{})
	go func() {
		defer close(readCh)
		for {
			select {
			case <-stopCh:
				return
			default:
				n, err := ptmx.Read(buf)
				if err != nil || n == 0 {
					return
				}
			}
		}
	}()

	select {
	case <-readCh:
		// Completed reading
	case <-time.After(500 * time.Millisecond):
		// Timeout after 500ms, signal to stop reading
		close(stopCh)
	}

	println(welcome)

	_, _ = ptmx.Write([]byte(newLine))

	go func() {
		_, _ = io.Copy(ptmx, os.Stdin)
	}()

	go func() {
		_, _ = io.Copy(os.Stdout, ptmx)
	}()

	return c.Wait()
}
