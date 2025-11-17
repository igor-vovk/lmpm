package installer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter"
	"github.com/hubblew/pim/internal/config"
	"github.com/hubblew/pim/internal/ui"
	"github.com/spf13/afero"
)

type Installer struct {
	fs afero.Fs
}

type Options struct {
	Config       *config.Config
	UserPrompter UserPrompter
}

func NewInstaller(fs afero.Fs) *Installer {
	return &Installer{
		fs: fs,
	}
}

func (i *Installer) Install(options *Options) error {
	tempDir, err := os.MkdirTemp("", "pim-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Printf("failed to remove temp directory '%s': %v\n", path, err)
		}
	}(tempDir)

	sourceDirsByName := make(map[string]string, len(options.Config.Sources))

	for _, source := range options.Config.Sources {
		if info, err := os.Stat(source.URL); err == nil && info.IsDir() {
			sourceDirsByName[source.Name] = source.URL

			continue
		}

		var sourceDir = filepath.Join(tempDir, source.Name)

		err = ui.RunWithSpinner(
			fmt.Sprintf("Fetching source '%s' from %s...\n", source.Name, source.URL),
			func() error {
				client := &getter.Client{
					Src:  source.URL,
					Dst:  sourceDir,
					Mode: getter.ClientModeDir,
				}

				if err := client.Get(); err != nil {
					return err
				}

				sourceDirsByName[source.Name] = sourceDir
				return nil
			})
		if err != nil {
			return fmt.Errorf("failed to fetch source '%s': %w", source.Name, err)
		}

		fmt.Printf("Source '%s' fetched successfully.\n", source.Name)
	}

	for _, target := range options.Config.Targets {
		if err := InstallTarget(i, &target, sourceDirsByName, options.UserPrompter); err != nil {
			return err
		}
	}

	fmt.Println("Installation complete!")
	return nil
}

func InstallTarget(i *Installer, target *config.Target, sourceDirsByName map[string]string, prompter UserPrompter) error {
	fmt.Printf("Installing target '%s' to %s...\n", target.Name, target.Output)

	strategy, err := NewStrategy(i.fs, target.StrategyType, target.Output)
	if err != nil {
		return fmt.Errorf("failed to create strategy for target '%s': %w", target.Name, err)
	}

	if err := strategy.Initialize(prompter); err != nil {
		return err
	}
	defer func(strategy Strategy) {
		if err := strategy.Close(); err != nil {
			fmt.Printf("failed to close strategy for target '%s': %v\n", target.Name, err)
		}
	}(strategy)

	for _, include := range target.IncludeParsed {
		sourceDir, ok := sourceDirsByName[include.Source]
		if !ok {
			return fmt.Errorf("source '%s' not found", include.Source)
		}

		srcPath := filepath.Join(sourceDir, include.File)

		// Use Glob to handle both literal paths and wildcard patterns
		matches, err := afero.Glob(i.fs, srcPath)
		if err != nil {
			return fmt.Errorf("failed to expand pattern '%s': %w", include.File, err)
		}

		if len(matches) == 0 {
			return fmt.Errorf("no files matched pattern '%s'", include.File)
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

	return nil
}
