package installer

import (
	"fmt"
	"os"

	"github.com/hubblew/pim/internal/utils"
)

const generatedByPim = "github.com/hubblew/pim-cli"

type frontmatterHeader struct {
	GeneratedBy string `yaml:"generatedBy"`
}

func defaultHeader() frontmatterHeader {
	return frontmatterHeader{
		GeneratedBy: generatedByPim,
	}
}

// IsPimGenerated checks if the markdown file at the given path
// contains a frontmatter block with the "generatedBy" key set to "github.com/hubblew/pim-cli".
func IsPimGenerated(path string) (bool, error) {
	var frontmatter frontmatterHeader
	err := utils.ReadFrontmatter(path, &frontmatter)
	if err != nil {
		return false, fmt.Errorf("failed to extract frontmatter block: %w", err)
	}

	// Check for generated-by key
	if frontmatter.GeneratedBy == generatedByPim {
		return true, nil
	}
	return false, nil
}

// AddGeneratedByPimHeader writes the PIM generation marker to the given file's frontmatter.
func AddGeneratedByPimHeader(file *os.File) error {
	header := defaultHeader()

	if err := utils.WriteFrontmatter(file, header); err != nil {
		return fmt.Errorf("failed to write frontmatter: %w", err)
	}
	return nil
}
