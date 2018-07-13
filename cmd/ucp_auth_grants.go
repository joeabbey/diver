package cmd

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
	"github.com/thebsdbox/diver/pkg/ucp/types"
)

func init() {

	// UCP Grant flags
	ucpAuthGrantsList.Flags().BoolVar(&resolve, "resolve", false, "Resolve the UUIDs to subject,role and grant names")

	ucpAuthGrantsSet.Flags().StringVar(&name, "subject", "", "The subject (user/org) that will be used")
	ucpAuthGrantsSet.Flags().StringVar(&ruleset, "role", "", "The role providing the capabilites")
	ucpAuthGrantsSet.Flags().StringVar(&collection, "collection", "", "The collection that will user")
	ucpAuthGrantsSet.Flags().StringVar(&collectionType, "type", "collection", "Type of grant: collection / namespace / all")

	ucpAuthGrantsDelete.Flags().StringVar(&name, "subject", "", "The subject (user/org) that will be used")
	ucpAuthGrantsDelete.Flags().StringVar(&ruleset, "role", "", "The role providing the capabilites")
	ucpAuthGrantsDelete.Flags().StringVar(&collection, "collection", "", "The collection that will user")

	// UCP Grants
	ucpAuth.AddCommand(ucpAuthGrants)
	ucpAuthGrants.AddCommand(ucpAuthGrantsGet)
	ucpAuthGrants.AddCommand(ucpAuthGrantsList)

	if !DiverRO {
		ucpAuthGrants.AddCommand(ucpAuthGrantsSet)
		ucpAuthGrants.AddCommand(ucpAuthGrantsDelete)

	}

}

var ucpAuthGrants = &cobra.Command{
	Use:   "grants",
	Short: "Manage Docker EE Grants",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpAuthGrantsList = &cobra.Command{
	Use:   "list",
	Short: "List Docker EE Grants",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = client.GetGrants(resolve)
		if err != nil {
			parsererr := ucp.ParseUCPError([]byte(err.Error()))
			if parsererr != nil {
				log.Errorf("Error parsing UCP error: %v", parsererr)
				log.Debugf("Response %v", err)
			}
			return
		}
	},
}

var ucpAuthGrantsGet = &cobra.Command{
	Use:   "get",
	Short: "List all rules for a particular grant",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if name == "" {
			cmd.Help()
			log.Fatalln("No role name specified to download")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		rules, err := client.GetRoleRuleset(name, id)
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("%s", rules)
	},
}

var ucpAuthGrantsSet = &cobra.Command{
	Use:   "set",
	Short: "Set a new grant linking a user through a role to a collection",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if name == "" {
			cmd.Help()
			log.Fatalln("No subject specified")
		}

		if ruleset == "" {
			cmd.Help()
			log.Fatalln("No role specified")
		}

		if collection == "" {
			cmd.Help()
			log.Fatalln("No collection specified")
		}

		var grantFlag uint
		switch collectionType {
		case "collection":
			grantFlag = ucptypes.GrantCollection
		case "namespace":
			grantFlag = ucptypes.GrantNamespace
		case "all":
			grantFlag = ucptypes.GrantObject
		default:
			cmd.Help()
			log.Fatalf("Unknown Grant type [%s]", collectionType)
		}

		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}

		err = client.SetGrant(collection, ruleset, name, grantFlag)
		if err != nil {
			ucp.ParseUCPError([]byte(err.Error()))
			return
		}
		log.Infof("Grant for user [%s] created succesfully", name)
	},
}

var ucpAuthGrantsDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a grant in Docker EE",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if name == "" {
			cmd.Help()
			log.Fatalln("No subject specified")
		}

		if ruleset == "" {
			cmd.Help()
			log.Fatalln("No role specified")
		}

		if collection == "" {
			cmd.Help()
			log.Fatalln("No collection specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}

		err = client.DeleteGrant(collection, ruleset, name)
		if err != nil {
			ucp.ParseUCPError([]byte(err.Error()))
			return
		}
		log.Infof("Grant for user [%s] deleted succesfully", name)
	},
}
