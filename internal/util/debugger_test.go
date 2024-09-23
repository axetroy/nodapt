package util

import (
	"bytes"
	"os"
	"testing"
)

func TestDebug(t *testing.T) {
	// Save the original stderr
	origStderr := os.Stderr
	defer func() { os.Stderr = origStderr }()

	// Create a pipe to capture stderr output
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Set the DEBUG environment variable to "1"
	os.Setenv("DEBUG", "1")
	defer os.Unsetenv("DEBUG")

	// Call the Debug function
	Debug("Test message: %d", 123)

	// Close the writer and read the output
	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	// Check if the output is as expected
	expected := "Test message: 123"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestDebugNoOutput(t *testing.T) {
	// Save the original stderr
	origStderr := os.Stderr
	defer func() { os.Stderr = origStderr }()

	// Create a pipe to capture stderr output
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Ensure the DEBUG environment variable is not set
	os.Unsetenv("DEBUG")

	// Call the Debug function
	Debug("This should not be printed")

	// Close the writer and read the output
	w.Close()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	// Check if the output is empty
	if buf.String() != "" {
		t.Errorf("expected no output, got %q", buf.String())
	}
}
