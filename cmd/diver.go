package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var diverCmd = &cobra.Command{
	Use:   "diver",
	Short: "This tool uses the native APIs to \"dive\" into Docker EE",
}

var logLevel int

func init() {
	// Global flag across all subcommands
	diverCmd.PersistentFlags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")
}

// Execute - starts the command parsing process
func Execute() {
	if err := diverCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
