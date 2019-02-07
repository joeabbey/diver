package cmd

import (
	"github.com/joeabbey/diver/pkg/ucp"
	"github.com/joeabbey/diver/pkg/ucp/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	// UCP Team flags
	ucpAuthTeamsList.Flags().StringVar(&name, "org", "", "Name of the organisation to query for teams")

	ucpAuthTeamsCreate.Flags().StringVar(&org, "org", "", "Name of the organisation to add a team to")
	ucpAuthTeamsCreate.Flags().StringVar(&name, "team", "", "Name of the new team")
	ucpAuthTeamsCreate.Flags().StringVar(&description, "description", "", "Description for the new team")

	ucpAuthTeamsDelete.Flags().StringVar(&org, "org", "", "Name of the organisation to remove a team from")
	ucpAuthTeamsDelete.Flags().StringVar(&name, "team", "", "Name of the team to be removed")

	ucpAuthTeamsAddUser.Flags().StringVar(&org, "org", "", "Name of the organisation")
	ucpAuthTeamsAddUser.Flags().StringVar(&name, "team", "", "Team that will the user will be added to")
	ucpAuthTeamsAddUser.Flags().StringVar(&user, "user", "", "Username to be added to the team")

	ucpAuthTeamsDelUser.Flags().StringVar(&org, "org", "", "Name of the organisation")
	ucpAuthTeamsDelUser.Flags().StringVar(&name, "team", "", "Team that will the user will be added to")
	ucpAuthTeamsDelUser.Flags().StringVar(&user, "user", "", "Username to be added to the team")

	// UCP Team
	ucpAuth.AddCommand(ucpAuthTeams)
	ucpAuthTeams.AddCommand(ucpAuthTeamsList)
	if !DiverRO {
		ucpAuthTeams.AddCommand(ucpAuthTeamsCreate)
		ucpAuthTeams.AddCommand(ucpAuthTeamsDelete)
		ucpAuthTeams.AddCommand(ucpAuthTeamsAddUser)
		ucpAuthTeams.AddCommand(ucpAuthTeamsDelUser)
	}
}

var ucpAuthTeams = &cobra.Command{
	Use:   "teams",
	Short: "Manage Docker EE Teams",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpAuthTeamsCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a Docker EE Team",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if name == "" {
			cmd.Help()
			log.Fatalln("No Team name specified")
		}
		if org == "" {
			cmd.Help()
			log.Fatalln("No Organisation Specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}

		var newTeam ucptypes.Team
		newTeam.Name = name
		newTeam.Description = description

		err = client.AddTeamToOrganisation(&newTeam, org)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var ucpAuthTeamsDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Docker EE Team",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if name == "" {
			cmd.Help()
			log.Fatalln("No Team name specified")
		}
		if org == "" {
			cmd.Help()
			log.Fatalln("No Organisation Specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = client.DeleteTeamFromOrganisation(name, org)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var ucpAuthTeamsAddUser = &cobra.Command{
	Use:   "adduser",
	Short: "Add a User to a team",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if name == "" {
			cmd.Help()
			log.Fatalln("No Team name specified")
		}
		if org == "" {
			cmd.Help()
			log.Fatalln("No Organisation Specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}

		err = client.AddUserToTeam(user, org, name)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var ucpAuthTeamsDelUser = &cobra.Command{
	Use:   "deluser",
	Short: "Delete a User from a team",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if name == "" {
			cmd.Help()
			log.Fatalln("No Team name specified")
		}
		if org == "" {
			cmd.Help()
			log.Fatalln("No Organisation Specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}

		err = client.DelUserFromTeam(user, org, name)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var ucpAuthTeamsList = &cobra.Command{
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
