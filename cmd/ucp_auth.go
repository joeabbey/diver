package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
	"github.com/thebsdbox/diver/pkg/ucp/types"
)

var auth ucptypes.Account
var name, org, user, ruleset, collection, collectionType, description string
var admin, inactive, resolve bool

func init() {
	// Auth flags
	ucpAuth.Flags().StringVar(&exportPath, "exportCSV", "", "Export users to a file [csv currently supported]")
	ucpAuth.Flags().BoolVar(&exampleFile, "exampleCSV", false, "Create an example csv file [example_accounts.csv]")

	if !DiverRO {
		ucpAuth.Flags().StringVar(&importPath, "importCSV", "", "Import accounts from a file [csv currently supported]")
	}

	// UCP ROOT
	UCPRoot.AddCommand(ucpAuth)
}

var ucpAuth = &cobra.Command{
	Use:   "auth",
	Short: "Authorisation commands for users/orgs/teams, roles and grants",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if exampleFile == true {
			log.Infof("Creating example CSV file for UCP accounts [example_accounts.csv]")
			err := ucp.CreateExampleAccountCSV()
			if err != nil {
				// Fatal error if can't read the token
				log.Fatalf("%v", err)
			}
			return
		}
		// A file has been passed in, so parse it and return
		if importPath != "" {
			log.Info("Importing Accounts from CSV file")
			client, err := ucp.ReadToken()
			if err != nil {
				// Fatal error if can't read the token
				log.Fatalf("%v", err)
			}
			err = client.ImportAccountsFromCSV(importPath)
			if err != nil {
				log.Fatalf("%v", err)
			}
			log.Info("Import succesfull")
			return
		}

		// Export all users to a csv file at exportPath
		if exportPath != "" {
			log.Infof("Exporting Accounts to file [%s]", exportPath)
			client, err := ucp.ReadToken()
			if err != nil {
				// Fatal error if can't read the token
				log.Fatalf("%v", err)
			}
			err = client.ExportAccountsToCSV(exportPath)
			if err != nil {
				log.Fatalf("%v", err)
			}
			return
		}
		cmd.Help()
	},
}
