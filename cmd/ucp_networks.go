package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/joeabbey/diver/pkg/ucp"
)

var force bool
var containerID, addressv4, addressv6 string

func init() {

	ucpNetworkAttach.Flags().StringVar(&id, "id", "", "ID of the network")
	ucpNetworkAttach.Flags().StringVar(&containerID, "container", "", "ID/Name of the container")
	ucpNetworkAttach.Flags().StringVar(&addressv4, "ipv4", "", "The IPv4 address to give the container on the network")
	ucpNetworkAttach.Flags().StringVar(&addressv6, "ipv6", "", "The IPv6 address to give the container on the network")

	ucpNetworkDetach.Flags().StringVar(&id, "id", "", "ID of the network")
	ucpNetworkDetach.Flags().StringVar(&containerID, "container", "", "ID/Name of the container")
	ucpNetworkDetach.Flags().BoolVar(&force, "force", false, "Force the removal of the container from the network")

	ucpNetwork.AddCommand(ucpNetworkList)
	ucpNetwork.AddCommand(ucpNetworkAttach)
	ucpNetwork.AddCommand(ucpNetworkDetach)

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

		networks, err := client.GetNetworks()
		if err != nil {
			log.Fatalf("%v", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintf(w, "Name\tID\n")
		for i := range networks {
			fmt.Fprintf(w, "%s\t%s\n", networks[i].Name, networks[i].ID)
		}
		w.Flush()
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
		err = client.NetworkConnectContainer(containerID, id, addressv4, addressv6)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Infof("Succesfully added [%s] to Network [%s]", containerID, id)
	},
}

var ucpNetworkDetach = &cobra.Command{
	Use:   "detach",
	Short: "Detach a container from a network",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		if id == "" {
			log.Fatalf("No network ID has been specified")
		}

		if containerID == "" {
			log.Fatalf("No container ID has been specified")
		}

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.NetworkDisconnectContainer(containerID, id, force)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Infof("Succesfully removed [%s] from Network [%s]", containerID, id)
	},
}
