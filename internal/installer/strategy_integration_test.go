package installer

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

type mockPrompter struct {
	allowOverwrite bool
}

func (m *mockPrompter) ConfirmOverwrite(string) (bool, error) {
	return m.allowOverwrite, nil
}

func TestConcatStrategyIntegration(t *testing.T) {
	tests := []struct {
		name              string
		files             []struct{ path, content string }
		outputPath        string
		expectError       bool
		expectedFragments []string
	}{
		{
			name: "concatenate multiple files",
			files: []struct{ path, content string }{
				{"file1.md", "# File 1\nContent 1"},
				{"file2.md", "# File 2\nContent 2"},
			},
			outputPath:  "output.md",
			expectError: false,
			expectedFragments: []string{
				"---\ngeneratedBy: github.com/hubblew/pim-cli\n---",
				"# File 1\nContent 1",
				"# File 2\nContent 2",
			},
		},
		{
			name: "single file",
			files: []struct{ path, content string }{
				{"single.md", "# Single File\nSingle content"},
			},
			outputPath:  "output.md",
			expectError: false,
			expectedFragments: []string{
				"---\ngeneratedBy: github.com/hubblew/pim-cli\n---",
				"# Single File\nSingle content",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()

			for _, f := range tt.files {
				if err := afero.WriteFile(fs, f.path, []byte(f.content), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			strategy := NewConcatStrategy(fs, tt.outputPath)
			prompter := &mockPrompter{allowOverwrite: true}

			if err := strategy.Initialize(prompter); err != nil {
				if !tt.expectError {
					t.Fatalf("unexpected error on Initialize: %v", err)
				}
				return
			}

			for _, f := range tt.files {
				if err := strategy.AddFile(f.path, f.path); err != nil {
					if !tt.expectError {
						t.Fatalf("unexpected error on AddFile: %v", err)
					}
					return
				}
			}

			if err := strategy.Close(); err != nil {
				t.Fatalf("unexpected error on Close: %v", err)
			}

			output, err := afero.ReadFile(fs, tt.outputPath)
			if err != nil {
				t.Fatalf("failed to read output file: %v", err)
			}

			outputStr := string(output)
			for _, fragment := range tt.expectedFragments {
				if !strings.Contains(outputStr, fragment) {
					t.Errorf("output missing expected fragment:\n%s\n\nFull output:\n%s", fragment, outputStr)
				}
			}
		})
	}
}

func TestFlattenStrategyIntegration(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"dir1/file1.md":        "Content 1",
		"dir1/file2.md":        "Content 2",
		"dir2/nested/file3.md": "Content 3",
	}

	for path, content := range files {
		if err := fs.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := afero.WriteFile(fs, path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	strategy := NewFlattenStrategy(fs, "output")
	prompter := &mockPrompter{allowOverwrite: true}

	if err := strategy.Initialize(prompter); err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}

	for path := range files {
		if err := strategy.AddFile(path, path); err != nil {
			t.Fatalf("failed to add file %s: %v", path, err)
		}
	}

	if err := strategy.Close(); err != nil {
		t.Fatalf("failed to close: %v", err)
	}

	expectedFiles := []string{"file1.md", "file2.md", "file3.md"}
	for _, fname := range expectedFiles {
		outputPath := filepath.Join("output", fname)
		if exists, _ := afero.Exists(fs, outputPath); !exists {
			t.Errorf("expected file %s to exist", outputPath)
		}

		content, err := afero.ReadFile(fs, outputPath)
		if err != nil {
			t.Errorf("failed to read %s: %v", outputPath, err)
		}

		if !strings.Contains(string(content), "Content") {
			t.Errorf("file %s has unexpected content: %s", fname, string(content))
		}
	}
}

func TestPreserveStrategyIntegration(t *testing.T) {
	fs := afero.NewMemMapFs()

	files := map[string]string{
		"dir1/file1.md":        "Content 1",
		"dir1/file2.md":        "Content 2",
		"dir2/nested/file3.md": "Content 3",
	}

	for path, content := range files {
		if err := fs.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := afero.WriteFile(fs, path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create file: %v", err)
		}
	}

	strategy := NewPreserveStrategy(fs, "output")
	prompter := &mockPrompter{allowOverwrite: true}

	if err := strategy.Initialize(prompter); err != nil {
		t.Fatalf("failed to initialize: %v", err)
	}

	for path := range files {
		if err := strategy.AddFile(path, path); err != nil {
			t.Fatalf("failed to add file %s: %v", path, err)
		}
	}

	if err := strategy.Close(); err != nil {
		t.Fatalf("failed to close: %v", err)
	}

	for origPath, expectedContent := range files {
		outputPath := filepath.Join("output", origPath)

		if exists, _ := afero.Exists(fs, outputPath); !exists {
			t.Errorf("expected file %s to exist", outputPath)
			continue
		}

		content, err := afero.ReadFile(fs, outputPath)
		if err != nil {
			t.Errorf("failed to read %s: %v", outputPath, err)
			continue
		}

		if string(content) != expectedContent {
			t.Errorf("file %s content mismatch\nexpected: %s\ngot: %s",
				outputPath, expectedContent, string(content))
		}
	}
}
