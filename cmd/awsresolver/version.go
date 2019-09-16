package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// Version contains the current version.
	Version = "dev"

	// BuildDate contains a string with the build date.
	BuildDate = "unknown"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Long:  `Print the version number and build date of this binary.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("awsresolver %s\n", Version)
		fmt.Printf("  Build date: %s\n", BuildDate)
		fmt.Printf("  Built with: %s\n", runtime.Version())
	},
}
