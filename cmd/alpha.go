package cmd

import (
	"github.com/spf13/cobra"
)

// alphaCmd represents the alpha command
var alphaCmd = &cobra.Command{
	Use:   "alpha",
	Short: "Alpha features that are not yet ready for production use",
	Long: `Groups commands that are considered experimental.

These commands might change or be removed in future versions.
Use with caution.`,
}

func init() {
	rootCmd.AddCommand(alphaCmd)
}
