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

// Execute -
func Execute() {
	if err := diverCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
