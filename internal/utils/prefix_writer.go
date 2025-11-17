package utils

import (
	"io"
)

// PrefixWriter wraps an io.Writer and adds a prefix to each line
type PrefixWriter struct {
	writer      io.Writer
	prefix      []byte
	atLineStart bool
}

var _ io.Writer = (*PrefixWriter)(nil)

func NewPrefixWriter(w io.Writer, prefix string) *PrefixWriter {
	return &PrefixWriter{
		writer:      w,
		prefix:      []byte(prefix),
		atLineStart: true,
	}
}

func (pw *PrefixWriter) Write(p []byte) (n int, err error) {
	n = len(p)

	for len(p) > 0 {
		if pw.atLineStart && len(p) > 0 {
			if _, err := pw.writer.Write(pw.prefix); err != nil {
				return n, err
			}
			pw.atLineStart = false
		}

		i := 0
		for i < len(p) && p[i] != '\n' {
			i++
		}

		if i > 0 {
			if _, err := pw.writer.Write(p[:i]); err != nil {
				return n, err
			}
		}

		if i < len(p) && p[i] == '\n' {
			if _, err := pw.writer.Write([]byte{'\n'}); err != nil {
				return n, err
			}
			pw.atLineStart = true
			i++
		}

		p = p[i:]
	}

	return n, nil
}
