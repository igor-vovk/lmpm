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

func NewConcatStrategy(path string) *ConcatStrategy {
	return &ConcatStrategy{
		outputPath: path,
	}
}

func (s *ConcatStrategy) Prepare() error {
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
