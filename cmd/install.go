package cmd

import (
	"fmt"
	"os"

	"github.com/hubble-works/pim/internal/config"
	"github.com/hubble-works/pim/internal/installer"
	"github.com/spf13/cobra"
)

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

		configPath := "pim.yaml"
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			configPath = ".pim.yaml"
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				return fmt.Errorf("configuration file not found (pim.yaml or .pim.yaml)")
			}
		}

		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		inst := installer.New(cfg)
		if err := inst.Install(); err != nil {
			return fmt.Errorf("installation failed: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
