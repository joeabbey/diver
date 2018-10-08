package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/joeabbey/diver/pkg/ucp"
	"github.com/joeabbey/diver/pkg/ucp/types"
)

func init() {

	// UCP ORG Flags
	ucpAuthOrgCreate.Flags().StringVar(&auth.Name, "name", "", "A unique Organisation name")

	ucpAuthOrgDelete.Flags().StringVar(&auth.Name, "name", "", "Existing Organisation")

	// UCP ORG
	ucpAuth.AddCommand(ucpAuthOrg)
	ucpAuthOrg.AddCommand(ucpAuthOrgList)

	if !DiverRO {
		ucpAuthOrg.AddCommand(ucpAuthOrgCreate)
		ucpAuthOrg.AddCommand(ucpAuthOrgDelete)
	}
}

var ucpAuthOrg = &cobra.Command{
	Use:   "org",
	Short: "Manage Docker EE Organisations",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpAuthOrgCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a Docker EE Organisation",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if auth.Name == "" {
			cmd.Help()
			log.Fatalln("No Organisation name specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		auth.IsOrg = true
		err = client.AddAccount(&auth)
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
	},
}

var ucpAuthOrgDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Docker EE Organisation",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if auth.Name == "" {
			cmd.Help()
			log.Fatalln("No Organisation name specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = client.DeleteAccount(auth.Name)
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
	},
}

var ucpAuthOrgList = &cobra.Command{
	Use:   "list",
	Short: "List all organisations in Docker EE",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Create query and set it to organisation type
		var accountQuery ucptypes.Account
		accountQuery.IsOrg = true

		orgs, err := client.GetAccounts(accountQuery, 1000)
		if err != nil {
			parseerr := ucp.ParseUCPError([]byte(err.Error()))
			if parseerr != nil {
				log.Debugf("Error parsing error")
			}
			return
		}

		if len(orgs.Accounts) == 0 {
			log.Error("No accounts returned")
			return
		}
		log.Debugf("Found %d Accounts", len(orgs.Accounts))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)

		fmt.Fprintf(w, "Org Name\tID\tFullname\n")

		for _, acct := range orgs.Accounts {
			fmt.Fprintf(w, "%s\t%s\t%s\n", acct.Name, acct.ID, acct.FullName)
		}
		w.Flush()
	},
}
