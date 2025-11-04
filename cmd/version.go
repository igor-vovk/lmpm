package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via -ldflags
	Version = "dev"
	// Commit is set at build time via -ldflags
	Commit = "none"
	// Date is set at build time via -ldflags
	Date = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of PIM",
	Long:  `Display the current version of PIM (Prompt Instruction Manager).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("PIM %s\n", Version)
		fmt.Printf("  Commit: %s\n", Commit)
		fmt.Printf("  Built:  %s\n", Date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
