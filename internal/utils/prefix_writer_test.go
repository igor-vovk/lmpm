package utils

import (
	"bytes"
	"testing"
)

func TestPrefixWriter(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		input    string
		expected string
	}{
		{
			name:     "single line",
			prefix:   "  ",
			input:    "hello",
			expected: "  hello",
		},
		{
			name:     "multiple lines",
			prefix:   "  ",
			input:    "line1\nline2\nline3",
			expected: "  line1\n  line2\n  line3",
		},
		{
			name:     "empty line in middle",
			prefix:   "  ",
			input:    "line1\n\nline3",
			expected: "  line1\n  \n  line3",
		},
		{
			name:     "trailing newline",
			prefix:   "  ",
			input:    "line1\nline2\n",
			expected: "  line1\n  line2\n",
		},
		{
			name:     "custom prefix",
			prefix:   ">>> ",
			input:    "output\nmore output",
			expected: ">>> output\n>>> more output",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			pw := NewPrefixWriter(&buf, tt.prefix)

			n, err := pw.Write([]byte(tt.input))
			if err != nil {
				t.Fatalf("Write failed: %v", err)
			}
			if n != len(tt.input) {
				t.Errorf("Write returned %d, expected %d", n, len(tt.input))
			}

			got := buf.String()
			if got != tt.expected {
				t.Errorf("got %q, expected %q", got, tt.expected)
			}
		})
	}
}

func TestPrefixWriter_MultipleWrites(t *testing.T) {
	var buf bytes.Buffer
	pw := NewPrefixWriter(&buf, "  ")

	writes := []string{
		"first",
		" line\n",
		"second line\n",
		"third",
	}

	for _, w := range writes {
		if _, err := pw.Write([]byte(w)); err != nil {
			t.Fatalf("Write failed: %v", err)
		}
	}

	expected := "  first line\n  second line\n  third"
	got := buf.String()
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}
