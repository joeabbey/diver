package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/joeabbey/diver/pkg/ucp"

	log "github.com/Sirupsen/logrus"
)

var ucpClient ucp.Client
var importPath, exportPath, action string
var top, exampleFile bool

func init() {

	// Add UCP and subcommands to the main application
	diverCmd.AddCommand(UCPRoot)

}

// UCPRoot - This is the root of all UCP commands / flags
var UCPRoot = &cobra.Command{
	Use:   "ucp",
	Short: "Universal Control Plane ",
	Run: func(cmd *cobra.Command, args []string) {

		existingClient, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			cmd.Help()
			log.Errorf("%v", err)
			return
		}
		currentAccount, err := existingClient.AuthStatus()
		if err != nil {
			cmd.Help()
			log.Warn("Session has expired, please login")
			return
		}
		cmd.Help()
		fmt.Printf("\n\n")
		log.Infof("Current user [%s]", currentAccount.Name)
		return
	},
}
