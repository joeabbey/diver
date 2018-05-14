package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

func init() {
	// Service flags
	ucpService.Flags().StringVar(&name, "name", "", "Examine a service by name")
	ucpService.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")
	UCPRoot.AddCommand(ucpService)
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
		log.Debugf("Looking for service [%s]", name)

		if name != "" {
			query := ucp.ServiceQuery{
				ServiceName: name,
			}
			err = client.QueryServiceContainers(&query)
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
