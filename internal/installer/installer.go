package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	getter "github.com/hashicorp/go-getter"
	"github.com/hubble-works/pim/internal/config"
)

type Installer struct {
	config *config.Config
}

func New(cfg *config.Config) *Installer {
	return &Installer{
		config: cfg,
	}
}

func (i *Installer) Install() error {
	tempDir, err := os.MkdirTemp("", "pim-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	sourceCache := make(map[string]string)

	for _, source := range i.config.Sources {
		sourceDir := filepath.Join(tempDir, source.Name)

		fmt.Printf("Fetching source '%s' from %s...\n", source.Name, source.URL)

		client := &getter.Client{
			Src:  source.URL,
			Dst:  sourceDir,
			Mode: getter.ClientModeDir,
		}

		if err := client.Get(); err != nil {
			return fmt.Errorf("failed to fetch source '%s': %w", source.Name, err)
		}

		sourceCache[source.Name] = sourceDir
	}

	for _, target := range i.config.Targets {
		fmt.Printf("Installing target '%s' to %s...\n", target.Name, target.Output)

		if target.Strategy == config.StrategyConcat {
			// For concat strategy, create the output file
			if err := os.MkdirAll(filepath.Dir(target.Output), 0755); err != nil {
				return fmt.Errorf("failed to create output directory '%s': %w", filepath.Dir(target.Output), err)
			}

			outFile, err := os.Create(target.Output)
			if err != nil {
				return fmt.Errorf("failed to create output file '%s': %w", target.Output, err)
			}
			defer outFile.Close()

			// Concatenate all files
			for _, include := range target.Include {
				sourceDir, ok := sourceCache[include.Source]
				if !ok {
					return fmt.Errorf("source '%s' not found", include.Source)
				}

				for _, file := range include.Files {
					srcPath := filepath.Join(sourceDir, file)

					if err := appendFileToOutput(srcPath, outFile); err != nil {
						return fmt.Errorf("failed to append file '%s': %w", file, err)
					}

					fmt.Printf("  ✓ %s\n", file)
				}
			}
		} else {
			// For flatten and preserve strategies
			if err := os.MkdirAll(target.Output, 0755); err != nil {
				return fmt.Errorf("failed to create output directory '%s': %w", target.Output, err)
			}

			for _, include := range target.Include {
				sourceDir, ok := sourceCache[include.Source]
				if !ok {
					return fmt.Errorf("source '%s' not found", include.Source)
				}

				for _, file := range include.Files {
					srcPath := filepath.Join(sourceDir, file)

					var dstPath string
					if target.Strategy == config.StrategyFlatten {
						dstPath = filepath.Join(target.Output, filepath.Base(file))
					} else {
						dstPath = filepath.Join(target.Output, file)
					}

					if err := copyFile(srcPath, dstPath); err != nil {
						return fmt.Errorf("failed to copy file '%s': %w", file, err)
					}

					fmt.Printf("  ✓ %s\n", file)
				}
			}
		}
	}

	fmt.Println("Installation complete!")
	return nil
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

func appendFileToOutput(srcPath string, outFile *os.File) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Copy file content
	if _, err := io.Copy(outFile, srcFile); err != nil {
		return err
	}

	// Add newline at the end
	if _, err := outFile.WriteString("\n"); err != nil {
		return err
	}

	return nil
}
