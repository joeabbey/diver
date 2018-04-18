package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	diverCmd.AddCommand(dtrCmd)

}

var dtrCmd = &cobra.Command{
	Use:   "dtr",
	Short: "Docker Trusted Registry",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
