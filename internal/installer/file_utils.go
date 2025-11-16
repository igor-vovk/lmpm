package installer

import (
	"io"
	"path/filepath"

	"github.com/spf13/afero"
)

func CopyFile(fs afero.Fs, src, dst string) error {
	if err := fs.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	srcFile, err := fs.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := fs.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	srcInfo, err := fs.Stat(src)
	if err != nil {
		return err
	}

	return fs.Chmod(dst, srcInfo.Mode())
}

func HasMdExtension(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".md"
}
