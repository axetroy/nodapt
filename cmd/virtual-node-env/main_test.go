package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Input struct {
	Args []string `json:"args"`
	Code int      `json:"code"`
}

func TestMain(t *testing.T) {
	cwd, err := os.Getwd()

	assert.Nil(t, err)

	fixtures := filepath.Join(cwd, "fixtures")

	files, err := os.ReadDir(fixtures)

	assert.Nil(t, err)

	for _, f := range files {
		fName := f.Name()

		if f.IsDir() {
			continue
		}

		if !strings.HasSuffix(fName, ".input") {
			continue
		}

		fPath := filepath.Join(fixtures, fName)

		buff, err := os.ReadFile(fPath)

		assert.Nil(t, err)

		var args = []string{"run", "./main.go"}

		var input Input

		assert.Nil(t, json.Unmarshal(buff, &input))

		args = append(args, input.Args...)

		ps := exec.Command("go", args...)

		output, err := ps.CombinedOutput()

		assert.Nil(t, err)

		outFilepath := filepath.Join(fixtures, strings.Replace(fName, ".input", ".output", 1))

		expectOut, err := os.ReadFile(outFilepath)

		assert.Nil(t, err)

		t.Logf("Running cli test '%s' ...", fName)
		assert.Equal(t, string(expectOut), string(output))
	}

}
