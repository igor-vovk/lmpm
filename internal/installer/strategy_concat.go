package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ConcatStrategy struct {
	outputPath string
	outFile    *os.File
}

var _ Strategy = (*ConcatStrategy)(nil)

func NewConcatStrategy(path string) *ConcatStrategy {
	return &ConcatStrategy{
		outputPath: path,
	}
}

func (s *ConcatStrategy) Initialize(prompter UserPrompter) error {
	if _, err := os.Stat(s.outputPath); err == nil {
		isGeneratedByPim, err := IsPimGenerated(s.outputPath)
		if err != nil {
			return fmt.Errorf("failed to check if file can be overridden: %w", err)
		}
		if !isGeneratedByPim {
			allowOverride, err := prompter.ConfirmOverwrite(s.outputPath)
			if err != nil {
				return fmt.Errorf("failed to prompt for file overwrite: %w", err)
			}
			if !allowOverride {
				return fmt.Errorf("user declined to override file '%s'", s.outputPath)
			}
		}
	}

	if err := os.Remove(s.outputPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete output file '%s': %w", s.outputPath, err)
	}

	if err := os.MkdirAll(filepath.Dir(s.outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory '%s': %w", filepath.Dir(s.outputPath), err)
	}

	outFile, err := os.Create(s.outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", s.outputPath, err)
	}

	s.outFile = outFile

	err = AddGeneratedByPimHeader(s.outFile)
	if err != nil {
		return fmt.Errorf("failed to write frontmatter to output file '%s': %w", s.outputPath, err)
	}

	return nil
}

func (s *ConcatStrategy) AddFile(srcPath, _ string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file '%s': %w", srcPath, err)
	}
	defer srcFile.Close()

	if _, err := io.Copy(s.outFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file '%s': %w", srcPath, err)
	}

	if _, err := s.outFile.WriteString("\n"); err != nil {
		return fmt.Errorf("failed to write newline to output file '%s': %w", s.outputPath, err)
	}

	return nil
}

func (s *ConcatStrategy) Close() error {
	if s.outFile != nil {
		return s.outFile.Close()
	}
	return nil
}
