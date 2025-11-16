package installer

import "testing"

func TestHasMdExtension(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"/output/file.md", true},
		{"/output/file.txt", false},
		{"/output/file.MD", false}, // case sensitive
		{"/output/dir", false},
		{"output.md", true},
		{"output.txt", false},
	}

	for _, tt := range tests {
		result := HasMdExtension(tt.path)
		if result != tt.expected {
			t.Errorf("HasMdExtension(%q) = %v, expected %v", tt.path, result, tt.expected)
		}
	}
}
