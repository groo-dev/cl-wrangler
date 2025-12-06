package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version follows semantic versioning (https://semver.org)
// Set via ldflags: -ldflags "-X github.com/groo-dev/cl-wrangler/cli/cmd.Version=1.0.0"
var Version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cl v%s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
