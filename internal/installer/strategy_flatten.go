package installer

import (
	"fmt"
	"os"
	"path/filepath"
)

type FlattenStrategy struct {
	outputPath string
}

func NewFlattenStrategy(path string) *FlattenStrategy {
	return &FlattenStrategy{
		outputPath: path,
	}
}

func (s *FlattenStrategy) Prepare() error {
	if err := os.RemoveAll(s.outputPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete output directory '%s': %w", s.outputPath, err)
	}

	if err := os.MkdirAll(s.outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory '%s': %w", s.outputPath, err)
	}
	return nil
}

func (s *FlattenStrategy) AddFile(srcPath, relativePath string) error {
	dstPath := filepath.Join(s.outputPath, filepath.Base(relativePath))
	return CopyFile(srcPath, dstPath)
}

func (s *FlattenStrategy) Close() error {
	return nil
}
