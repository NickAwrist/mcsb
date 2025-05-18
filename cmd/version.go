package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version     = "dev"
	ShortCommit = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		ExecuteVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func ExecuteVersion() {
	fmt.Printf("Version: %s\nCommit: %s\n", Version, ShortCommit)
}
