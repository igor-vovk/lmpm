package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

const delimiter = "---"

// ReadFrontmatter reads the file at the given path and extracts the frontmatter block if present.
// The frontmatter is expected to be in YAML format, enclosed between '---' delimiters.
// The extracted frontmatter is unmarshalled into the provided variable v.
func ReadFrontmatter(path string, v any) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close() // Error on read-only file close is not critical
	}()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return nil
	}

	firstLine := strings.TrimSpace(scanner.Text())
	if firstLine != delimiter {
		return nil
	}

	var builder strings.Builder
	foundClosing := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == delimiter {
			foundClosing = true
			break
		}
		builder.WriteString(line)
		builder.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if !foundClosing {
		return nil
	}

	if err := yaml.Unmarshal([]byte(builder.String()), v); err != nil {
		return fmt.Errorf("failed to parse frontmatter YAML: %w", err)
	}

	return nil
}

func WriteFrontmatter(w io.Writer, v any) error {
	frontmatterBytes, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	_, err = w.Write([]byte(fmt.Sprintf("%s\n%s%s\n\n", delimiter, string(frontmatterBytes), delimiter)))

	return err
}
