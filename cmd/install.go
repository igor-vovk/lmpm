package cmd

import (
	"fmt"
	"os"

	"github.com/hubblew/pim/internal/config"
	"github.com/hubblew/pim/internal/installer"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var configPathFlag string
var forceFlag bool

const DefaultConfigFileName = "pim.yaml"

var installCmd = &cobra.Command{
	Use:   "install [directory]",
	Short: "Install packages from sources to targets",
	Long:  `Fetch sources and copy specified files to target directories.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		if err := os.Chdir(dir); err != nil {
			return fmt.Errorf("failed to change to directory %s: %w", dir, err)
		}

		if _, err := os.Stat(configPathFlag); os.IsNotExist(err) {
			return fmt.Errorf("configuration file not found: %s", configPathFlag)
		}

		cfg, err := config.LoadConfig(configPathFlag)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		var prompter installer.UserPrompter
		if forceFlag {
			prompter = installer.NewAcceptAllPrompter()
		} else {
			prompter = installer.NewInteractivePrompter()
		}

		opts := installer.Options{
			Config:       cfg,
			UserPrompter: prompter,
		}

		fs := afero.NewOsFs()
		if err := installer.NewInstaller(fs).Install(&opts); err != nil {
			return fmt.Errorf("installation failed: %w", err)
		}

		return nil
	},
}

func init() {
	installCmd.Flags().StringVarP(
		&configPathFlag,
		"config",
		"c",
		DefaultConfigFileName,
		"Path to configuration file",
	)
	installCmd.Flags().BoolVarP(
		&forceFlag,
		"force",
		"f",
		false,
		"Force overwrite existing files without prompting",
	)

	rootCmd.AddCommand(installCmd)
}
