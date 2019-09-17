package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const macOSResolverFilePath = "/etc/resolver/internal"
const macOSResolverConfig = `domain internal
nameserver 127.0.0.1
port 1053
timeout 1
search_order 1
`

func checkMacOSResolverConfig() error {
	//  TODO: Check permissions and ownership?
	if _, err := os.Stat(macOSResolverFilePath); err == nil {
		contents, readErr := ioutil.ReadFile(macOSResolverFilePath)
		if readErr != nil {
			return errors.Wrap(readErr, "could not read macOS resolver config")
		}
		conf := string(contents)
		if conf == macOSResolverConfig {
			// Config is correct
			return nil
		} else if conf == "" {
			return fmt.Errorf("resolver config in %s is empty.\nExpected: `%s`", macOSResolverFilePath, macOSResolverConfig)
		} else {
			return fmt.Errorf("unexpected config in %s.\nExpected: `%s`\nGot: `%s`)", macOSResolverFilePath, macOSResolverConfig, conf)
		}
	} else if os.IsNotExist(err) {
		return fmt.Errorf("macOS resolver config does not exist in %s", macOSResolverFilePath)
	} else {
		return errors.Wrap(err, "unknown macOS resolver config error")
	}
}

func setupMacOSResolverConfig() error {
	dirName := filepath.Dir(macOSResolverFilePath)
	errDir := os.MkdirAll(dirName, 0755)
	if errDir != nil {
		return errors.Wrap(errDir, "could not create macOS resolver directory")
	}
	err := ioutil.WriteFile(macOSResolverFilePath, []byte(macOSResolverConfig), 0644)
	if err != nil {
		return errors.Wrap(err, "could not write macOS resolver config")
	}
	return nil
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Perform the initial resolver setup",
	Long:  `Sets up the macOS resolver to route ".internal" DNS queries to awsresolver.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if runtime.GOOS != "darwin" {
			fmt.Println("error: setup currently only supports macOS")
			os.Exit(1)
		}

		if os.Geteuid() != 0 {
			fmt.Println("error: setup must be run as root")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := setupMacOSResolverConfig(); err != nil {
			fmt.Printf("error: could not setup macOS resolver config: %v\n", err.Error())
			os.Exit(1)
		}

		fmt.Println("Checking work...")

		if err := checkMacOSResolverConfig(); err != nil {
			fmt.Printf("error: macOS resolver config still has errors: %v\n", err.Error())
			os.Exit(1)
		}

		// TODO: Add daemon setup.

		fmt.Println("Setup complete. You may now start with: `awsresolver run`")
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
