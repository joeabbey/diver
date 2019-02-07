package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/joeabbey/diver/pkg/ucp"
)

func init() {
	ucpContainerProcesses.Flags().StringVar(&id, "id", "", "Container ID to retrieve processes from")
	ucpContainerNetworks.Flags().StringVar(&id, "id", "", "Container ID to retrieve processes from")

	ucpContainerGet.AddCommand(ucpContainerProcesses)
	ucpContainerGet.AddCommand(ucpContainerNetworks)

	// Sub commands
	ucpContainer.AddCommand(ucpContainerTop)
	ucpContainer.AddCommand(ucpContainerList)
	ucpContainer.AddCommand(ucpContainerGet)
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
		containers, err := client.GetAllContainers()
		if err != nil {
			log.Fatalf("%v", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Name\tID\tImage\tNode")

		for i := range containers {
			// String returned from UCP is /<node>/<container name>
			splitUCPName := strings.Split(containers[i].Names[0], "/")

			// Image name is  image:tag@sha256
			splitImageSha := strings.Split(containers[i].Image, "@")

			// Shrink the imageID to a manageble size
			shrunkID := fmt.Sprintf("%.16s", containers[i].ID)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", splitUCPName[2], shrunkID, splitImageSha[0], splitUCPName[1])
		}
		w.Flush()

	},
}

var ucpContainerGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve specific information about a container",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpContainerNetworks = &cobra.Command{
	Use:   "networks",
	Short: "List all networks a running container is attached too",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		if id == "" {
			log.Fatalf("No Container ID specified")
		}
		container, err := client.GetContainerFromID(id)
		if err != nil {
			log.Fatalf("%v", err)
		}
		if container.NetworkSettings == nil {
			log.Info("No networks")
		}
		networks := container.NetworkSettings.Networks
		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Name\tID\tAddress")
		for name, network := range networks {
			fmt.Fprintf(w, "%s\t%s\t%s\n", name, network.NetworkID, network.IPAddress)
		}
		w.Flush()
	},
}

var ucpContainerProcesses = &cobra.Command{
	Use:   "processes",
	Short: "List all processes in a running container",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		if id == "" {
			log.Fatalf("No Container ID specified")
		}
		processes, err := client.GetContainerProcesses(id)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Debugf("Found %d process lines", len(processes.Processes))
		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		// Prepare Title headers
		for i := range processes.Titles {
			fmt.Fprintf(w, "%s\t", processes.Titles[i])
		}
		// End of line for title
		fmt.Fprintf(w, "\n")

		// Process lines
		for i := range processes.Processes {
			for x := range processes.Processes[i] {
				fmt.Fprintf(w, "%s\t", processes.Processes[i][x])
			}
			fmt.Fprintf(w, "\n")
		}
		w.Flush() // Flush adds a \n
	},
}
