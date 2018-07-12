package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
	"github.com/thebsdbox/diver/pkg/ucp/types"
)

func init() {

	// UCP User flags
	ucpAuthUsersDelete.Flags().StringVar(&auth.Name, "name", "", "Existing username")

	ucpAuthUsersList.Flags().BoolVar(&admin, "admin", false, "Retrieve *only* Administrative users")
	ucpAuthUsersList.Flags().BoolVar(&inactive, "inactive", false, "Retrieve *only* inactive users")

	ucpAuthUsersCreate.Flags().StringVar(&auth.FullName, "fullname", "", "The full name of a UCP user or organisation")
	ucpAuthUsersCreate.Flags().StringVar(&auth.Name, "name", "", "The unique username")
	ucpAuthUsersCreate.Flags().StringVar(&auth.Password, "password", "", "A string password for a new user of organisation")
	ucpAuthUsersCreate.Flags().BoolVar(&auth.IsAdmin, "admin", false, "Make this user an administrator")
	ucpAuthUsersCreate.Flags().BoolVar(&auth.IsActive, "active", true, "Enable this user in the Universal Control Plane")

	// UCP USERS
	ucpAuth.AddCommand(ucpAuthUsers)
	ucpAuthUsers.AddCommand(ucpAuthUsersList)
	if !DiverRO {
		ucpAuthUsers.AddCommand(ucpAuthUsersCreate)
		ucpAuthUsers.AddCommand(ucpAuthOrgDelete)
	}
}

var ucpAuthUsers = &cobra.Command{
	Use:   "users",
	Short: "Manage Docker EE Users",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpAuthUsersCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a Docker EE User",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if auth.Name == "" {
			cmd.Help()
			log.Fatalln("No Username specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		auth.IsOrg = false
		err = client.AddAccount(&auth)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var ucpAuthUsersDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Docker EE User",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if auth.Name == "" {
			cmd.Help()
			log.Fatalln("No Username specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = client.DeleteAccount(auth.Name)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var ucpAuthUsersList = &cobra.Command{
	Use:   "list",
	Short: "List all users in Docker EE",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		var users *ucptypes.AccountList
		var accountQuery ucptypes.Account
		accountQuery.IsOrg = false
		if admin {
			accountQuery.IsAdmin = true
		}
		if inactive {
			accountQuery.IsActive = true
		}
		users, err = client.GetAccounts(accountQuery, 1000)
		if err != nil {
			parsererr := ucp.ParseUCPError([]byte(err.Error()))
			if parsererr != nil {
				log.Errorf("Error parsing UCP error: %v", parsererr)
				log.Debugf("Response %v", err)
			}
			return
		}

		if len(users.Accounts) == 0 {
			log.Error("No accounts returned")
			return
		}
		const padding = 3
		w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
		log.Debugf("Found %d Accounts", len(users.Accounts))
		fmt.Fprintln(w, "User Name\tID\tFullname\t")
		for _, acct := range users.Accounts {
			// Not sure why we're still retrieving ORGs even though we said false above - TODO
			if !acct.IsOrg {
				fmt.Fprintf(w, "%s\t%s\t%s\n", acct.Name, acct.ID, acct.FullName)
			}
		}
		w.Flush()
	},
}
