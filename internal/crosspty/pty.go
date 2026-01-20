package crosspty

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/aymanbagabas/go-pty"
	"github.com/axetroy/nodapt/internal/util"
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
	if _, err := os.Stderr.WriteString(welcome + "\n"); err != nil {
		// Non-fatal, just log the error
		util.Debug("Warning: failed to write welcome message: %v\n", err)
	}

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

	time.Sleep(1000 * time.Millisecond) // Give the shell some time to start.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	shellName := filepath.Base(shellPath)

	newLine := getNewLine(shellName)

	// Set environment variables
	for k, v := range env {
		var command string
		
		switch shellName {
		case "bash", "zsh":
			// Escape single quotes for bash/zsh: replace ' with '"'"'
			escapedValue := strings.Replace(v, "'", "'\"'\"'", -1)
			command = fmt.Sprintf("export %s='%s'", k, escapedValue) + newLine
		case "fish":
			// Escape single quotes for fish: replace ' with '"'"'
			escapedValue := strings.Replace(v, "'", "'\"'\"'", -1)
			command = fmt.Sprintf("set -gx %s '%s'", k, escapedValue) + newLine
		case "powershell", "powershell.exe":
			// For PowerShell, escape single quotes by doubling them: ' becomes ''
			escapedValue := strings.Replace(v, "'", "''", -1)
			command = fmt.Sprintf("$env:%s='%s'", k, escapedValue) + newLine
		case "cmd", "cmd.exe":
			// For CMD, wrap in quotes and escape quotes: " becomes ""
			// Also escape special chars like &, |, <, >, ^, %, !
			escapedValue := strings.NewReplacer(
				"\"", "\"\"",
				"&", "^&",
				"|", "^|",
				"<", "^<",
				">", "^>",
				"^", "^^",
				"%", "%%",
			).Replace(v)
			command = fmt.Sprintf("set \"%s=%s\"", k, escapedValue) + newLine
		default:
			// Default to bash-style escaping
			escapedValue := strings.Replace(v, "'", "'\"'\"'", -1)
			command = fmt.Sprintf("export %s='%s'", k, escapedValue) + newLine
		}
		
		_, _ = ptmx.Write([]byte(command))
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

	// Clear the PTY output buffer with a timeout to prevent blocking
	// Use channels to coordinate goroutine cleanup
	buf := make([]byte, 1024)
	done := make(chan struct{})
	stop := make(chan struct{})
	
	go func() {
		defer close(done)
		// Read and discard output until stopped or error
		for {
			select {
			case <-stop:
				return
			default:
				n, err := ptmx.Read(buf)
				if err != nil || n == 0 {
					return
				}
			}
		}
	}()

	// Wait for either completion or timeout
	select {
	case <-done:
		// Goroutine completed successfully
	case <-time.After(200 * time.Millisecond):
		// Timeout - signal goroutine to stop
		close(stop)
		// Wait for goroutine to acknowledge the stop signal
		<-done
	}

	_, _ = ptmx.Write([]byte(newLine))

	go func() {
		_, _ = io.Copy(ptmx, os.Stdin)
	}()

	go func() {
		_, _ = io.Copy(os.Stdout, ptmx)
	}()

	return c.Wait()
}
