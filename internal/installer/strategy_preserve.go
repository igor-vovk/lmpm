package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type PreserveStrategy struct {
	fs         afero.Fs
	outputPath string
}

var _ Strategy = (*PreserveStrategy)(nil)

func NewPreserveStrategy(fs afero.Fs, path string) *PreserveStrategy {
	return &PreserveStrategy{
		fs:         fs,
		outputPath: path,
	}
}

func (s *PreserveStrategy) Initialize(_ UserPrompter) error {
	if err := s.fs.RemoveAll(s.outputPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete output directory '%s': %w", s.outputPath, err)
	}

	if err := s.fs.MkdirAll(s.outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory '%s': %w", s.outputPath, err)
	}
	return nil
}

func (s *PreserveStrategy) AddFile(srcPath, relativePath string) error {
	dstPath := filepath.Join(s.outputPath, relativePath)
	return CopyFile(s.fs, srcPath, dstPath)
}

func (s *PreserveStrategy) Close() error {
	return nil
}
