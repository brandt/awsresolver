package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/brandt/awsresolver/internal/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "awsresolver",
	Short: "Resolve AWS internal hostnames",
	Long: `Resolves AWS internal EC2 FQDNs where the IP is embedded in the hostname.
For example: ip-10-78-32-168.us-east-2.compute.internal resolves to 10.78.32.168`,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the resolver",
	Long:  `Run the DNS resolver server in the foreground.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "darwin" {
			fmt.Println("warning: limited support for operating systems other than macOS")
			return
		}

		// Check if the resolver file exists. If it doesn't, this might be the
		// first run, so we give a more direct error message.
		if _, err := os.Stat(macOSResolverFilePath); os.IsNotExist(err) {
			fmt.Println("macOS resolver not configured. Please run: `sudo awsresolver setup`")
		} else if err := checkMacOSResolverConfig(); err != nil {
			fmt.Printf("macOS resolver not correctly configured: %v\n", err.Error())
			fmt.Println("This might be fixable by re-running setup: `sudo awsresolver setup`")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting server...")
		server.Start()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
