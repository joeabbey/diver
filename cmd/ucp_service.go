package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

var svc ucp.ServiceQuery

func init() {
	// Service flags
	ucpService.Flags().StringVar(&svc.ServiceName, "name", "", "Examine a service by name")

	// Query options
	ucpService.Flags().BoolVar(&svc.ID, "id", false, "Display task ID")
	ucpService.Flags().BoolVar(&svc.Networks, "networks", false, "Display task Network connections")
	ucpService.Flags().BoolVar(&svc.State, "state", false, "Display task state")
	ucpService.Flags().BoolVar(&svc.Node, "node", false, "Display Node running task")
	ucpService.Flags().BoolVar(&svc.Resolve, "resolve", false, "Resolve Task IDs to human readable names")

	// Set logging
	ucpService.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")

	// Service Reap flags
	ucpServiceReap.Flags().StringVar(&svc.ServiceName, "name", "", "Examine a service by name")
	ucpServiceReap.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")

	// Add Service to UCP root commands
	UCPRoot.AddCommand(ucpService)

	// Add reap to service subcommands
	ucpService.AddCommand(ucpServiceReap)
}

var ucpService = &cobra.Command{
	Use:   "service",
	Short: "Interact with services",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		log.Debugf("Looking for service [%s]", svc.ServiceName)

		if svc.ServiceName != "" {
			err = client.QueryServiceContainers(&svc)
			if err != nil {
				log.Fatalf("%v", err)
			}
			return
		}

		err = client.GetServices()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var ucpServiceReap = &cobra.Command{
	Use:   "reap",
	Short: "Clean a service",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		log.Debugf("Looking for service [%s]", svc.ServiceName)

		if svc.ServiceName != "" {
			err = client.QueryServiceContainers(&svc)
			if err != nil {
				log.Fatalf("%v", err)
			}
			return
		}
	},
}
