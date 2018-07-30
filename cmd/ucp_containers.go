package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

func init() {

	// Sub commands
	ucpContainer.AddCommand(ucpContainerTop)
	ucpContainer.AddCommand(ucpContainerList)

	UCPRoot.AddCommand(ucpContainer)
}

var ucpContainer = &cobra.Command{
	Use:   "containers",
	Short: "Interact with containers",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpContainerTop = &cobra.Command{
	Use:   "top",
	Short: "A list of containers and their CPU usage like the top command on linux",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.ContainerTop()
		if err != nil {
			log.Fatalf("%v", err)
		}
		return
	},
}

var ucpContainerList = &cobra.Command{
	Use:   "list",
	Short: "List all containers across all nodes in UCP",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.GetContainerNames()
		if err != nil {
			log.Fatalf("%v", err)
		}
		return
	},
}
