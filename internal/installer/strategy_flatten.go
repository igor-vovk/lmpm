package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

type FlattenStrategy struct {
	fs         afero.Fs
	outputPath string
}

var _ Strategy = (*FlattenStrategy)(nil)

func NewFlattenStrategy(fs afero.Fs, path string) *FlattenStrategy {
	return &FlattenStrategy{
		fs:         fs,
		outputPath: path,
	}
}

func (s *FlattenStrategy) Initialize(_ UserPrompter) error {
	if err := s.fs.RemoveAll(s.outputPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete output directory '%s': %w", s.outputPath, err)
	}

	if err := s.fs.MkdirAll(s.outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory '%s': %w", s.outputPath, err)
	}
	return nil
}

func (s *FlattenStrategy) AddFile(srcPath, relativePath string) error {
	dstPath := filepath.Join(s.outputPath, filepath.Base(relativePath))
	return CopyFile(s.fs, srcPath, dstPath)
}

func (s *FlattenStrategy) Close() error {
	return nil
}
