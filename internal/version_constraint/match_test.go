package version_constraint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionMatch(t *testing.T) {
	tests := []struct {
		constraint  string
		version     string
		expected    bool
		expectError bool
	}{
		{"^1.0.0", "1.0.0", true, false},
		{"^1.0.0", "1.1.0", true, false},
		{"^1.0.0", "2.0.0", false, false},
		{"~1.2.3", "1.2.4", true, false},
		{"~1.2.3", "1.3.0", false, false},
		{"1.2.x", "1.2.3", true, false},
		{"1.2.x", "1.3.0", false, false},
		{"invalid", "1.0.0", false, true},
		{"^1.0.0", "invalid", false, true},
		{"<=1.2.3", "1.2.3", true, false},
		{"<=1.2.3", "1.2.4", false, false},
		{">=1.2.3", "1.2.3", true, false},
		{">=1.2.3", "1.2.2", false, false},
		{"1.2.3 - 1.2.5", "1.2.4", true, false},
		{"1.2.3 - 1.2.5", "1.2.6", false, false},
		{"*", "1.0.0", true, false},
		{"*", "2.0.0", true, false},
	}

	for _, test := range tests {
		result, err := Match(test.constraint, test.version)
		if test.expectError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, result)
		}
	}
}
