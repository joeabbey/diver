package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/joeabbey/diver/pkg/ucp"
)

// Used to enable or disable orchestrator type
var orchestratorKube, orchestratorSwarm bool

// Set a node to a swarm availability state
var availability string

// Set a node to a specific role type
var role string

// Set a label on a node
var labelKey, labelValue string

func init() {

	ucpNodesAvailability.Flags().StringVar(&id, "id", "", "ID of the Docker Node")
	ucpNodesAvailability.Flags().StringVar(&availability, "state", "active", "Node availability [active/drain/pause]")

	ucpNodesDelete.Flags().StringVar(&id, "id", "", "ID of the Docker Node")
	ucpNodesDelete.Flags().BoolVar(&force, "force", false, "Force the removal of this node")

	ucpNodesGet.Flags().StringVar(&id, "id", "", "ID of the Docker Node")

	ucpNodesLabel.Flags().StringVar(&id, "id", "", "ID of the Docker Node")
	ucpNodesLabel.Flags().StringVar(&labelKey, "key", "", "The label Key")
	ucpNodesLabel.Flags().StringVar(&labelValue, "value", "", "The label Value")

	ucpNodesOrchestrator.Flags().StringVar(&id, "id", "", "ID of the Docker Node")
	ucpNodesOrchestrator.Flags().BoolVar(&orchestratorKube, "kubernetes", false, "Enable Kubernetes to use this node")
	ucpNodesOrchestrator.Flags().BoolVar(&orchestratorSwarm, "swarm", false, "Enable Swarm to use this node")

	ucpNodesRole.Flags().StringVar(&id, "id", "", "ID of the Docker Node")
	ucpNodesRole.Flags().StringVar(&role, "role", "", "Node role [manager/worker]")

	ucpNodes.AddCommand(ucpNodesAvailability)
	ucpNodes.AddCommand(ucpNodesDelete)
	ucpNodes.AddCommand(ucpNodesGet)
	ucpNodes.AddCommand(ucpNodesLabel)
	ucpNodes.AddCommand(ucpNodesList)
	ucpNodes.AddCommand(ucpNodesOrchestrator)
	ucpNodes.AddCommand(ucpNodesRole)

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
		fmt.Fprintln(w, "Name\tID\tRole\tVersion\tPlatform\tSwarm\tKubernetes")
		for i := range nodes {
			// Combine OS/Arch
			platform := nodes[i].Description.Platform.OS + "/" + nodes[i].Description.Platform.Architecture

			// Determine Orchestrator configuration
			orchestratorKube, err = strconv.ParseBool(nodes[i].Spec.Labels["com.docker.ucp.orchestrator.kubernetes"])
			if err != nil {
				// If there is an error it means that the label isn't part of the spec, default to disabled
				orchestratorKube = false
			}

			orchestratorSwarm, err = strconv.ParseBool(nodes[i].Spec.Labels["com.docker.ucp.orchestrator.swarm"])
			if err != nil {
				// If there is an error it means that the label isn't part of the spec, default to disabled
				orchestratorSwarm = false
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%t\t%t\n", nodes[i].Description.Hostname,
				nodes[i].ID,
				nodes[i].Spec.Role,
				nodes[i].Description.Engine.EngineVersion,
				platform,
				orchestratorSwarm,
				orchestratorKube)
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
	Short: "Configure which orchestrators can utilise a node",
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

		// If both orchestrators are false, then neither can schedule workloads (display a warning)
		if orchestratorKube == false && orchestratorSwarm == false {
			log.Warn("This node has no orchestrators defined and wont be scheduled any workload")
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
		log.Infof("Configured Node [%s] to allow kubernetes=%t and swarm=%t", id, orchestratorKube, orchestratorSwarm)
	},
}

var ucpNodesAvailability = &cobra.Command{
	Use:   "availability",
	Short: "Set the node availability [active/drain/pause]",
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

		err = client.SetNodeAvailability(id, availability)
		if err != nil {
			log.Fatalf("%v", err)
		}

		log.Infof("Succesfully set node [%s] to state [%s]", id, availability)
	},
}

var ucpNodesRole = &cobra.Command{
	Use:   "role",
	Short: "Set the node role [manager/worker]",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if id == "" {
			cmd.Help()
			log.Fatalln("No Node ID specified")
		}

		if role == "" {
			cmd.Help()
			log.Fatalln("No Node Role specified, should be either manager or worker")
		}

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.SetNodeRole(id, role)
		if err != nil {
			log.Fatalf("%v", err)
		}

		log.Infof("Succesfully set node [%s] to swarm role [%s]", id, role)
	},
}

var ucpNodesLabel = &cobra.Command{
	Use:   "label",
	Short: "Set a label and value on a node",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if id == "" {
			cmd.Help()
			log.Fatalln("No Node ID specified")
		}
		if labelKey == "" {
			cmd.Help()
			log.Fatalln("No label key has been specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.SetNodeLabel(id, labelKey, labelValue)
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		log.Infof("Succesfully updated node [%s] with the label [%s=%s]", id, labelKey, labelValue)
	},
}

var ucpNodesDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Docker node, including removing from UCP (KV/Auth stores etc.)",
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
		err = client.DeleteNode(id, force)
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		log.Infof("Succesfully removed node [%s]", id)

	},
}
