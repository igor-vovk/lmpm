package installer

import (
	"fmt"
	"os"
	"path/filepath"
)

type PreserveStrategy struct {
	outputPath string
}

func NewPreserveStrategy(path string) *PreserveStrategy {
	return &PreserveStrategy{
		outputPath: path,
	}
}

func (s *PreserveStrategy) Prepare() error {
	if err := os.RemoveAll(s.outputPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete output directory '%s': %w", s.outputPath, err)
	}

	if err := os.MkdirAll(s.outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory '%s': %w", s.outputPath, err)
	}
	return nil
}

func (s *PreserveStrategy) AddFile(srcPath, relativePath string) error {
	dstPath := filepath.Join(s.outputPath, relativePath)
	return CopyFile(srcPath, dstPath)
}

func (s *PreserveStrategy) Close() error {
	return nil
}
