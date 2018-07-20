package cmd

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var diverCmd = &cobra.Command{
	Use:   "diver",
	Short: "This tool uses the native APIs to \"dive\" into Docker EE",
}

var logLevel int

// DiverVersion is the release TAG
var DiverVersion string

// DiverBuild is the current GIT commit
var DiverBuild string

// DiverRO Sets Diver to READ ONLY (all SET/CREATE commands are disabled)
var DiverRO bool

// Defines the padding for tabwritter through all diver output
const tabPadding = 3

func init() {
	// Global flag across all subcommands
	diverCmd.PersistentFlags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")
	diverCmd.AddCommand(diverVersion)
}

// Execute - starts the command parsing process
func Execute() {
	log.SetLevel(log.Level(logLevel))
	if DiverRO {
		log.Debugf("ReadWrite commands disabled")
	}
	if err := diverCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var diverVersion = &cobra.Command{
	Use:   "version",
	Short: "Version and Release information about the diver tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Diver, a CLI tool to interact with the Docker EE APIs\n")
		fmt.Printf("Version:  %s\n", DiverVersion)
		fmt.Printf("Build:    %s\n", DiverBuild)
		fmt.Printf("ReadOnly: %t\n", DiverRO)

	},
}
