package node

import (
	"os"
	"testing"
)

func TestGetVersionFromPackageJSON(t *testing.T) {
	tests := []struct {
		name           string
		packageJSON    string
		expectedResult *string
		expectError    bool
	}{
		{
			name: "Valid Node version",
			packageJSON: `{
				"engines": {
					"node": "14.x"
				}
			}`,
			expectedResult: strPtr("14.x"),
			expectError:    false,
		},
		{
			name: "No engines field",
			packageJSON: `{
				"name": "example"
			}`,
			expectedResult: nil,
			expectError:    false,
		},
		{
			name: "No node field in engines",
			packageJSON: `{
				"engines": {
					"npm": "6.x"
				}
			}`,
			expectedResult: nil,
			expectError:    false,
		},
		{
			name: "Invalid JSON",
			packageJSON: `{
				"engines": {
					"node": "14.x"
				`,
			expectedResult: nil,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tmpfile, err := os.CreateTemp("", "package.json")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			// Write the packageJSON content to the temp file
			if _, err := tmpfile.Write([]byte(tt.packageJSON)); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			if err := tmpfile.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}

			// Call the function under test
			result, err := GetConstraintFromPackageJSON(tmpfile.Name())

			// Check for unexpected errors
			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			// Check the result
			if (result == nil && tt.expectedResult != nil) || (result != nil && tt.expectedResult == nil) || (result != nil && *result != *tt.expectedResult) {
				t.Errorf("Expected result: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
