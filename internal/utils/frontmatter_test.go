package utils

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestReadFrontmatter(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		wantKey     string
		wantValue   string
		wantErr     bool
	}{
		{
			name: "valid frontmatter",
			fileContent: `---
key: value
---
Regular content here`,
			wantKey:   "key",
			wantValue: "value",
			wantErr:   false,
		},
		{
			name: "valid frontmatter with trailing content",
			fileContent: `---
title: Test
---
# Content`,
			wantKey:   "title",
			wantValue: "Test",
			wantErr:   false,
		},
		{
			name: "empty frontmatter",
			fileContent: `---
---
Content`,
			wantErr: false,
		},
		{
			name: "no frontmatter delimiter at start",
			fileContent: `Regular content
---
Not frontmatter
---`,
			wantErr: false,
		},
		{
			name: "missing closing delimiter",
			fileContent: `---
key: value
no closing delimiter`,
			wantErr: false,
		},
		{
			name:        "empty file",
			fileContent: "",
			wantErr:     false,
		},
		{
			name: "frontmatter with whitespace",
			fileContent: `   ---   
key: value
   ---   
Content`,
			wantKey:   "key",
			wantValue: "value",
			wantErr:   false,
		},
		{
			name: "invalid YAML",
			fileContent: `---
key: value
  invalid: yaml: syntax
---
Content`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.md")

			err := os.WriteFile(tmpFile, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}

			var got map[string]string
			err = ReadFrontmatter(tmpFile, &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFrontmatter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.wantKey != "" {
				if got[tt.wantKey] != tt.wantValue {
					t.Errorf("ReadFrontmatter() got[%q] = %q, want %q", tt.wantKey, got[tt.wantKey], tt.wantValue)
				}
			}
		})
	}
}

func TestReadFrontmatter_FileNotFound(t *testing.T) {
	nonExistentPath := filepath.Join(t.TempDir(), "does_not_exist.md")

	var result map[string]string
	err := ReadFrontmatter(nonExistentPath, &result)

	if err == nil {
		t.Error("ReadFrontmatter() expected error for non-existent file, got nil")
	}
}

func TestWriteFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]string
		wantErr bool
	}{
		{
			name: "simple map",
			input: map[string]string{
				"key":   "value",
				"title": "Test",
			},
			wantErr: false,
		},
		{
			name:    "empty map",
			input:   map[string]string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := WriteFrontmatter(&buf, tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("WriteFrontmatter() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				output := buf.String()
				if len(output) < 10 || output[:3] != "---" {
					t.Errorf("WriteFrontmatter() output doesn't start with delimiter: %q", output)
				}
			}
		})
	}
}

func TestWriteFrontmatter_RoundTrip(t *testing.T) {
	original := map[string]string{
		"title":       "Test Document",
		"description": "A test document",
		"author":      "Test Author",
	}

	var buf bytes.Buffer
	err := WriteFrontmatter(&buf, original)
	if err != nil {
		t.Fatalf("WriteFrontmatter() error = %v", err)
	}

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")

	err = os.WriteFile(tmpFile, buf.Bytes(), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	var result map[string]string
	err = ReadFrontmatter(tmpFile, &result)
	if err != nil {
		t.Fatalf("ReadFrontmatter() error = %v", err)
	}

	for k, v := range original {
		if result[k] != v {
			t.Errorf("Round trip failed for key %q: got %q, want %q", k, result[k], v)
		}
	}
}
