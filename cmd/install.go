package cmd

import (
	"fmt"
	"os"

	"github.com/hubblew/pim/internal/config"
	"github.com/hubblew/pim/internal/installer"
	"github.com/spf13/cobra"
)

var configFlag string

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

		configPath := configFlag
		if configPath == "" {
			// Auto-detect config file if not specified
			configPath = "pim.yaml"
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				configPath = ".pim.yaml"
				if _, err := os.Stat(configPath); os.IsNotExist(err) {
					return fmt.Errorf("configuration file not found (pim.yaml or .pim.yaml)")
				}
			}
		} else {
			// Use specified config file
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				return fmt.Errorf("configuration file not found: %s", configPath)
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
	installCmd.Flags().StringVarP(
		&configFlag,
		"config",
		"c",
		"",
		"Path to configuration file (default: pim.yaml or .pim.yaml)",
	)

	rootCmd.AddCommand(installCmd)
}
