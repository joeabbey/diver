package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/swarm"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

var svc ucp.ServiceQuery

var prevSpec, colour bool

func init() {
	// Service flags
	ucpServiceList.Flags().StringVar(&svc.ServiceName, "name", "", "Examine a service by name")

	ucpServiceGetTasks.Flags().StringVar(&svc.ServiceName, "name", "", "Examine a service by name")
	ucpServiceGetTasks.Flags().BoolVar(&colour, "colour", false, "Use Colour in Task output")
	ucpServiceGetHealth.Flags().StringVar(&svc.ServiceName, "name", "", "Examine a service by name")

	// Query options
	ucpServiceList.Flags().BoolVar(&svc.ID, "id", false, "Display task ID")
	ucpServiceList.Flags().BoolVar(&svc.Networks, "networks", false, "Display task Network connections")
	ucpServiceList.Flags().BoolVar(&svc.State, "state", false, "Display task state")
	ucpServiceList.Flags().BoolVar(&svc.Node, "node", false, "Display Node running task")
	ucpServiceList.Flags().BoolVar(&svc.Resolve, "resolve", false, "Resolve Task IDs to human readable names")

	// Service Reap flags
	ucpServiceReap.Flags().StringVar(&svc.ServiceName, "name", "", "Reap tasks from this Service")

	// Service Architecture flags
	ucpServiceArchitecture.Flags().BoolVar(&prevSpec, "previousSpec", false, "Display the previous Service specification")

	// Service Configuration Flags
	ucpServiceGetConfig.Flags().StringVar(&svc.ServiceName, "name", "", "The name of service to retrieve configurations from")

	// Add Service to UCP root commands
	UCPRoot.AddCommand(ucpService)

	ucpServiceGet.AddCommand(ucpServiceGetTasks)
	ucpServiceGet.AddCommand(ucpServiceGetHealth)

	// Add reap to service subcommands
	ucpService.AddCommand(ucpServiceList)
	ucpService.AddCommand(ucpServiceReap)
	ucpService.AddCommand(ucpServiceArchitecture)
	ucpService.AddCommand(ucpServiceGet)
	ucpService.AddCommand(ucpServiceGetConfig)

}

var ucpService = &cobra.Command{
	Use:   "services",
	Short: "Interact with services",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpServiceList = &cobra.Command{
	Use:   "list",
	Short: "List services",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		svcs, err := client.GetAllServices()
		if err != nil {
			log.Fatalf("%v", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Name\tID")
		for i := range svcs {

			fmt.Fprintf(w, "%s\t%s\n", svcs[i].Spec.Name, svcs[i].ID)
		}
		w.Flush()

	},
}

var ucpServiceGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve information about a Service",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpServiceGetHealth = &cobra.Command{
	Use:   "health",
	Short: "Retrieve the health of a service or all services",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		services, err := client.GetAllServices()
		if err != nil {
			log.Fatalf("%v", err)
		}
		// Find the specific service and output the health

		if svc.ServiceName != "" {
			for i := range services {
				if services[i].Spec.Name == svc.ServiceName {
					// Service has been found, get the tasks
					tasks, err := client.GetServiceTasks(svc.ServiceName)
					if err != nil {
						log.Fatalf("%v", err)
					}

					var running, failed, shutdown int

					for x := range tasks {
						// Loop through the tasks and work out the health
						log.Debugf("%s")
						switch tasks[x].Status.State {
						case "running":
							running++
						case "failed":
							failed++
						case "shutdown":
							shutdown++
						}
					}

					w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
					fmt.Fprintln(w, "Service\tExpected\tRunning\tShutdown\tFailed\tTasks\tStatus")

					// Check if Replica or Global Service

					if services[i].Spec.Mode.Replicated != nil && services[i].Spec.Mode.Replicated.Replicas != nil {
						fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t%d\n", svc.ServiceName,
							*services[i].Spec.Mode.Replicated.Replicas,
							running,
							shutdown,
							failed,
							len(tasks))
					} else {
						fmt.Fprintf(w, "%s\t---\t%d\t%d\t%d\t%d\n", svc.ServiceName,
							running,
							shutdown,
							failed,
							len(tasks))
					}

					w.Flush()
					return
				}
			}
			log.Fatalf("Service [%s] couldn't be found", svc.ServiceName)
		}

	},
}

var ucpServiceGetTasks = &cobra.Command{
	Use:   "tasks",
	Short: "Retrieve all tasks from a service",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		if svc.ServiceName == "" {
			cmd.Help()
			log.Fatalf("Please specify either a Service Name")
		}
		//retrieve all tasks
		tasks, err := client.GetServiceTasks(svc.ServiceName)
		if err != nil {
			log.Fatalf("%v", err)
		}

		log.Debugf("Found %d tasks for service %s", len(tasks), svc.ServiceName)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Hostame\tID\tNode\tIP Address\tNetwork\tState")

		for i := range tasks {
			// Retrieve task ID
			task := tasks[i].Status.ContainerStatus.ContainerID
			resolvedTask, err := client.GetContainerFromID(task)
			if colour {
				switch tasks[i].Status.State {
				case "running":
					fmt.Fprintf(w, "\x1b[32;1m")

				case "failed":
					fmt.Fprintf(w, "\x1b[31;1m")

				case "shutdown":
					fmt.Fprintf(w, "\x1b[34;1m")
				}
			}

			// Print details about container, unless it has already been removed
			if err != nil {
				fmt.Fprintf(w, "%s\t%s\t", "Removed", task)
			} else {
				fmt.Fprintf(w, "%s\t%s\t", resolvedTask.Name, resolvedTask.ID)
			}

			// Retrieve node for container
			containerNode, err := client.GetContainerFromID(tasks[i].Status.ContainerStatus.ContainerID)
			if err != nil {
				fmt.Fprintf(w, "Removed\t")
			} else {
				fmt.Fprintf(w, "%s\t", containerNode.Node.Name)
			}

			// Print all networks attached to task (Only if attachements exist)
			if len(tasks[i].NetworksAttachments) != 0 {
				var networkString string
				for n := range tasks[i].NetworksAttachments {
					for a := range tasks[i].NetworksAttachments[n].Addresses {

						address := tasks[i].NetworksAttachments[n].Addresses[a]
						networkID := tasks[i].NetworksAttachments[n].Network.ID

						// build output from query
						resolvedNetwork, err := client.GetNetworkFromID(networkID)
						if err != nil {
							return
						}
						// Build from resolved name
						networkString = networkString + address + "  " + resolvedNetwork.Name + "  "
					}
				}
				fmt.Fprintf(w, "%s", networkString)
			} else {
				fmt.Fprintf(w, "Unattached\t")
			}

			fmt.Fprintf(w, "%s\t", tasks[i].Status.State)

			// Create a newline for the next task
			fmt.Fprintf(w, "\n")
		}
		w.Flush()
		fmt.Printf("\x1b[0m")
	},
}

var ucpServiceReap = &cobra.Command{
	Use:   "reap",
	Short: "Clean a service (not implemented)",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		if svc.ServiceName != "" {
			err := client.ReapFailedTasks(svc.ServiceName, false, false)
			if err != nil {
				log.Fatalf("%v", err)
			}
			return
		}
	},
}

var ucpServiceArchitecture = &cobra.Command{
	Use:   "architecture",
	Short: "Retrieve the \"design\" of a service",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		if len(args) == 0 {
			log.Fatalf("Please specify a service name")
		}
		if len(args) > 1 {
			log.Fatalf("Please only specify a single service to view the architecture")
		}

		log.Infof("Inspecting service [%s]", args[0])

		service, err := client.GetService(args[0])

		if err != nil {
			ucp.ParseUCPError([]byte(err.Error()))
			return
		}

		var spec *swarm.ServiceSpec

		if prevSpec == true {
			spec = service.PreviousSpec
		} else {
			spec = &service.Spec
		}

		printServiceSpec(service, spec)
	},
}

// This will read through the service spec and print out the details
func printServiceSpec(service *swarm.Service, spec *swarm.ServiceSpec) {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

	fmt.Fprintf(w, "ID:\t%s\n", service.ID)
	fmt.Fprintf(w, "Version:\t%d\n", service.Version.Index)
	fmt.Fprintf(w, "Name:\t%s\n", spec.Name)
	fmt.Fprintf(w, "Image:\t%s\n", spec.TaskTemplate.ContainerSpec.Image)
	//Print out the command used for the image
	fmt.Fprintf(w, "Cmd:")
	for i := range spec.TaskTemplate.ContainerSpec.Command {
		fmt.Fprintf(w, " %s", spec.TaskTemplate.ContainerSpec.Command[i])
	}
	fmt.Fprintf(w, "\n")
	// Print all arguments to the command
	fmt.Fprintf(w, "Args:")
	for i := range spec.TaskTemplate.ContainerSpec.Args {
		fmt.Fprintf(w, " %s", spec.TaskTemplate.ContainerSpec.Args[i])
	}
	fmt.Fprintf(w, "\n")

	// Print the labels from the key/map
	fmt.Fprintf(w, "Labels:\n")
	for key, value := range spec.TaskTemplate.ContainerSpec.Labels {
		fmt.Fprintf(w, "\t%s", key)
		fmt.Fprintf(w, "\t%s\n", value)
	}

	//Print reservations
	if spec.TaskTemplate.Resources != nil {
		// Check if this struct exists
		if spec.TaskTemplate.Resources.Reservations != nil {
			fmt.Fprintf(w, "Memory Reservation:\t%d\n", spec.TaskTemplate.Resources.Reservations.MemoryBytes)
			fmt.Fprintf(w, "CPU Reservation:\t%d\n", spec.TaskTemplate.Resources.Reservations.NanoCPUs)
		}
		//Print limits
		if spec.TaskTemplate.Resources.Limits != nil {
			fmt.Fprintf(w, "Memory Limits:\t%d\n", spec.TaskTemplate.Resources.Limits.MemoryBytes)
			fmt.Fprintf(w, "CPU Limits:\t%d\n", spec.TaskTemplate.Resources.Limits.NanoCPUs)
		}
	}
	// Check the pointer to the replica struct isn't nil and read replica count
	if spec.Mode.Replicated != nil {
		fmt.Fprintf(w, "Replicas:\t%d\n", *spec.Mode.Replicated.Replicas)
	}
	if spec.Mode.Global != nil {
		fmt.Fprintf(w, "Global:\ttrue\n")
	}
	w.Flush()

}

var ucpServiceGetConfig = &cobra.Command{
	Use:   "config",
	Short: "Retrieve the config used by a service",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		if svc.ServiceName == "" {
			cmd.Help()
			log.Fatalf("Please specify a Service Name")
		}
		service, err := client.GetService(svc.ServiceName)

		if err != nil {
			ucp.ParseUCPError([]byte(err.Error()))
			log.Fatalf("Could not retrieve configuration from [%s]", svc.ServiceName)
			return
		}
		if service.Spec.TaskTemplate.ContainerSpec != nil && service.Spec.TaskTemplate.ContainerSpec.Configs != nil {

			w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
			fmt.Fprintf(w, "Name\tID\tFile\n")
			for i := range service.Spec.TaskTemplate.ContainerSpec.Configs {
				fmt.Fprintf(w, "%s\t%s\t%s\n", service.Spec.TaskTemplate.ContainerSpec.Configs[i].ConfigName, service.Spec.TaskTemplate.ContainerSpec.Configs[i].ConfigID, service.Spec.TaskTemplate.ContainerSpec.Configs[i].File.Name)
			}
			w.Flush()

		}
	},
}
