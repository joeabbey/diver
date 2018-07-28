package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

func init() {

	ucpNetwork.AddCommand(ucpNetworkList)
	UCPRoot.AddCommand(ucpNetwork)
}

var ucpNetwork = &cobra.Command{
	Use:   "network",
	Short: "Interact with container networks",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var ucpNetworkList = &cobra.Command{
	Use:   "list",
	Short: "list all container networks",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		err = client.GetNetworks()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var ucpNetworkAttach = &cobra.Command{
	Use:   "attach",
	Short: "Attach a container to a network",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		client.NetworkConnectContainer("", "", "", "")
	},
}

var ucpNetworkDetach = &cobra.Command{
	Use:   "detach",
	Short: "Detach a container from a network",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		client.NetworkConnectContainer("", "", "", "")
	},
}
