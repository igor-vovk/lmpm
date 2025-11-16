package utils

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"
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
			fs := afero.NewMemMapFs()
			testFile := "test.md"

			if err := afero.WriteFile(fs, testFile, []byte(tt.fileContent), 0644); err != nil {
				t.Fatalf("failed to create test file: %v", err)
			}

			var got map[string]string
			err := ReadFrontmatter(fs, testFile, &got)

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
	fs := afero.NewMemMapFs()

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

	testFile := "test.md"
	if err := afero.WriteFile(fs, testFile, buf.Bytes(), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	var result map[string]string
	err = ReadFrontmatter(fs, testFile, &result)
	if err != nil {
		t.Fatalf("ReadFrontmatter() error = %v", err)
	}

	for k, v := range original {
		if result[k] != v {
			t.Errorf("Round trip failed for key %q: got %q, want %q", k, result[k], v)
		}
	}
}
