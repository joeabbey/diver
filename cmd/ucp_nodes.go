package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

// Used to enable or disable orchestrator type
var orchestratorKube, orchestratorSwarm bool

func init() {
	ucpNodesGet.Flags().StringVar(&id, "id", "", "ID of the Docker Node")

	ucpNodesOrchestrator.Flags().StringVar(&id, "id", "", "ID of the Docker Node")
	ucpNodesOrchestrator.Flags().BoolVar(&orchestratorKube, "kubernetes", false, "Enable Kubernetes to use this node")
	ucpNodesOrchestrator.Flags().BoolVar(&orchestratorSwarm, "swarm", false, "Enable Swarm to use this node")

	ucpNodes.AddCommand(ucpNodesList)
	ucpNodes.AddCommand(ucpNodesGet)
	ucpNodes.AddCommand(ucpNodesOrchestrator)

	// Add nodes to UCP root commands
	UCPRoot.AddCommand(ucpNodes)

}

var ucpNodes = &cobra.Command{
	Use:   "nodes",
	Short: "Interact with Docker Nodes",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpNodesList = &cobra.Command{
	Use:   "list",
	Short: "Retrieve all Docker Nodes",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		nodes, err := client.ListAllNodes()
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Debugf("Found [%d] nodes", len(nodes))
		if len(nodes) == 0 {
			log.Fatalf("No Nodes found")
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Name\tID\tRole\tVersion\tPlatform")
		for i := range nodes {
			// Combine OS/Arch
			platform := nodes[i].Description.Platform.OS + "/" + nodes[i].Description.Platform.Architecture
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", nodes[i].Description.Hostname, nodes[i].ID, nodes[i].Spec.Role, nodes[i].Description.Engine.EngineVersion, platform)
		}
		w.Flush()
	},
}

var ucpNodesGet = &cobra.Command{
	Use:   "get",
	Short: "Get information about a particular Docker Node",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if id == "" {
			cmd.Help()
			log.Fatalln("No Node ID specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		node, err := client.GetNode(id)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Debugf("Retrieved information about [%s]", node.Description.Hostname)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Label Key\tLabel Value")
		for k, v := range node.Spec.Labels {
			fmt.Fprintf(w, "%s\t%s\n", k, v)
		}
		w.Flush()
	},
}

var ucpNodesOrchestrator = &cobra.Command{
	Use:   "orchestrator",
	Short: "Enable node for different orchestrators",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if id == "" {
			cmd.Help()
			log.Fatalln("No Node ID specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.SetNodeLabel(id, "com.docker.ucp.orchestrator.kubernetes", strconv.FormatBool(orchestratorKube))
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.SetNodeLabel(id, "com.docker.ucp.orchestrator.swarm", strconv.FormatBool(orchestratorSwarm))
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		log.Infof("Configured Node [%s] to use orchestrator kubernetes=[%t] / swarm=[%t]", id, orchestratorKube, orchestratorSwarm)
	},
}
