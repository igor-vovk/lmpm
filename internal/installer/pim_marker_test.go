package installer

import (
	"testing"

	"github.com/spf13/afero"
)

func TestIsPimGeneratedWithMemFS(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expected    bool
		expectError bool
	}{
		{
			name:        "pim generated file",
			content:     "---\ngeneratedBy: github.com/hubblew/pim-cli\n---\n\nContent",
			expected:    true,
			expectError: false,
		},
		{
			name:        "non-pim file",
			content:     "---\nauthor: someone\n---\n\nContent",
			expected:    false,
			expectError: false,
		},
		{
			name:        "no frontmatter",
			content:     "Just content",
			expected:    false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			path := "test.md"

			if err := afero.WriteFile(fs, path, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			result, err := IsPimGenerated(fs, path)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestAddGeneratedByPimHeader(t *testing.T) {
	fs := afero.NewMemMapFs()
	file, err := fs.Create("test.md")
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	if err := AddGeneratedByPimHeader(file); err != nil {
		t.Fatalf("failed to add header: %v", err)
	}

	file.Close()

	content, err := afero.ReadFile(fs, "test.md")
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	expected := "---\ngeneratedBy: github.com/hubblew/pim-cli\n---\n\n"
	if string(content) != expected {
		t.Errorf("header mismatch\nexpected:\n%s\ngot:\n%s", expected, string(content))
	}
}
