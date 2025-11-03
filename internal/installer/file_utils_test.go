package installer

import "testing"

func TestHasTextExtension(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/output/file.md", true},
		{"/output/file.txt", true},
		{"/output/file.MD", false},  // case sensitive
		{"/output/file.TXT", false}, // case sensitive
		{"/output/file.yaml", false},
		{"/output/file.json", false},
		{"/output/dir", false},
		{"output.md", true},
		{"output.txt", true},
	}

	for _, tt := range tests {
		result := HasTextExtension(tt.path)
		if result != tt.expected {
			t.Errorf("HasTextExtension(%q) = %v, expected %v", tt.path, result, tt.expected)
		}
	}
}
