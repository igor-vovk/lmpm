package installer

import (
	"fmt"
	"os"
	"path/filepath"

	getter "github.com/hashicorp/go-getter"
	"github.com/hubblew/pim/internal/config"
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

	sourceDirsByName := make(map[string]string)

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

		sourceDirsByName[source.Name] = sourceDir
	}

	for _, target := range i.config.Targets {
		fmt.Printf("Installing target '%s' to %s...\n", target.Name, target.Output)

		strategy, err := NewStrategy(target.StrategyType, target.Output)
		if err != nil {
			return fmt.Errorf("failed to create strategy for target '%s': %w", target.Name, err)
		}

		if err := strategy.Prepare(); err != nil {
			return err
		}
		defer strategy.Close()

		for _, include := range target.IncludeParsed {
			sourceDir, ok := sourceDirsByName[include.Source]
			if !ok {
				return fmt.Errorf("source '%s' not found", include.Source)
			}

			for _, file := range include.Files {
				srcPath := filepath.Join(sourceDir, file)

				// Use Glob to handle both literal paths and wildcard patterns
				matches, err := filepath.Glob(srcPath)
				if err != nil {
					return fmt.Errorf("failed to expand pattern '%s': %w", file, err)
				}

				if len(matches) == 0 {
					return fmt.Errorf("no files matched pattern '%s'", file)
				}

				for _, match := range matches {
					// Get the relative path from sourceDir
					relPath, err := filepath.Rel(sourceDir, match)
					if err != nil {
						return fmt.Errorf("failed to get relative path for '%s': %w", match, err)
					}

					if err := strategy.AddFile(match, relPath); err != nil {
						return fmt.Errorf("failed to add file '%s': %w", relPath, err)
					}

					fmt.Printf("  ✓ %s\n", relPath)
				}
			}
		}
	}

	fmt.Println("Installation complete!")
	return nil
}
