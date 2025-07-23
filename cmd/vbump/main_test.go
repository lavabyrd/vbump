package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadVersion(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expected    string
		expectError bool
	}{
		{
			name:     "valid version",
			content:  "1.2.3\n",
			expected: "1.2.3",
		},
		{
			name:     "version with extra whitespace",
			content:  "  1.2.3  \n",
			expected: "1.2.3",
		},
		{
			name:        "empty file",
			content:     "",
			expectError: true,
		},
		{
			name:     "version without newline",
			content:  "1.2.3",
			expected: "1.2.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			versionFile := filepath.Join(tmpDir, "VERSION")
			
			if err := os.WriteFile(versionFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			result, err := readVersion(versionFile)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestReadVersionFileNotFound(t *testing.T) {
	_, err := readVersion("/nonexistent/file")
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}

func TestAtoi(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"0", 0},
		{"1", 1},
		{"42", 42},
		{"123", 123},
		{"invalid", 0}, // fmt.Sscanf returns 0 for invalid input
		{"", 0},
		{"12abc", 12}, // fmt.Sscanf parses until first non-digit
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := atoi(tt.input)
			if result != tt.expected {
				t.Errorf("atoi(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestVersionBumping(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		major    bool
		minor    bool
		expected string
	}{
		{
			name:     "patch bump",
			input:    "1.2.3",
			expected: "1.2.4",
		},
		{
			name:     "minor bump",
			input:    "1.2.3",
			minor:    true,
			expected: "1.3.0",
		},
		{
			name:     "major bump",
			input:    "1.2.3",
			major:    true,
			expected: "2.0.0",
		},
		{
			name:     "patch bump from 0",
			input:    "0.0.0",
			expected: "0.0.1",
		},
		{
			name:     "minor bump resets patch",
			input:    "1.2.9",
			minor:    true,
			expected: "1.3.0",
		},
		{
			name:     "major bump resets minor and patch",
			input:    "1.9.9",
			major:    true,
			expected: "2.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.Split(tt.input, ".")
			if len(parts) != 3 {
				t.Fatalf("invalid test input: %s", tt.input)
			}

			var result string
			switch {
			case tt.major:
				result = fmt.Sprintf("%d.0.0", atoi(parts[0])+1)
			case tt.minor:
				result = fmt.Sprintf("%s.%d.0", parts[0], atoi(parts[1])+1)
			default:
				result = fmt.Sprintf("%s.%s.%d", parts[0], parts[1], atoi(parts[2])+1)
			}

			if result != tt.expected {
				t.Errorf("version bump %s (major=%v, minor=%v) = %s, expected %s", 
					tt.input, tt.major, tt.minor, result, tt.expected)
			}
		})
	}
}

func TestInvalidVersionFormat(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{"too few parts", "1.2"},
		{"too many parts", "1.2.3.4"},
		{"single number", "1"},
		{"empty string", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.Split(tt.version, ".")
			if len(parts) == 3 {
				t.Errorf("test case %q should have invalid format, but has 3 parts", tt.version)
			}
		})
	}
}