package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

func init() {
	// UCP Team flags
	ucpAuthTeamList.Flags().StringVar(&name, "org", "", "Name of the organisation to query for teams")

	// UCP Team
	ucpAuth.AddCommand(ucpAuthTeams)
	ucpAuthTeams.AddCommand(ucpAuthTeamList)

}

var ucpAuthTeams = &cobra.Command{
	Use:   "teams",
	Short: "Manage Docker EE Teams",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpAuthTeamList = &cobra.Command{
	Use:   "list",
	Short: "List Docker EE Teams",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		if name == "" {
			log.Fatalf("No Organisation Specified")
		}

		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = client.GetTeams(name)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}
