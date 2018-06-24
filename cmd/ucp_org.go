package cmd

import (
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

var auth ucp.Account
var name, ruleset, collection string
var admin, inactive, resolve bool

func init() {
	// Auth flags
	ucpAuth.Flags().StringVar(&importPath, "importCSV", "", "Import accounts from a file [csv currently supported]")
	ucpAuth.Flags().StringVar(&exportPath, "exportCSV", "", "Export users to a file [csv currently supported]")
	ucpAuth.Flags().BoolVar(&exampleFile, "exampleCSV", false, "Create an example csv file [example_accounts.csv]")
	ucpAuth.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")

	// User/Org Create flags
	ucpAuthOrgCreate.Flags().StringVar(&auth.Name, "name", "", "A unique Organisation name")
	ucpAuthOrgCreate.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")

	ucpAuthUsersCreate.Flags().StringVar(&auth.FullName, "fullname", "", "The full name of a UCP user or organisation")
	ucpAuthUsersCreate.Flags().StringVar(&auth.Name, "name", "", "The unique username")
	ucpAuthUsersCreate.Flags().StringVar(&auth.Password, "password", "", "A string password for a new user of organisation")
	ucpAuthUsersCreate.Flags().BoolVar(&auth.IsAdmin, "admin", false, "Make this user an administrator")
	ucpAuthUsersCreate.Flags().BoolVar(&auth.IsActive, "active", true, "Enable this user in the Universal Control Plane")
	ucpAuthUsersCreate.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")

	// User/Org Delete flags
	ucpAuthOrgDelete.Flags().StringVar(&auth.Name, "name", "", "Existing Organisation")
	ucpAuthOrgDelete.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")
	ucpAuthOrgList.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")

	ucpAuthUsersDelete.Flags().StringVar(&auth.Name, "name", "", "Existing username")
	ucpAuthUsersDelete.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")
	ucpAuthUsersList.Flags().BoolVar(&admin, "admin", false, "Retrieve *only* Administrative users")
	ucpAuthUsersList.Flags().BoolVar(&inactive, "inactive", false, "Retrieve *only* inactive users")
	ucpAuthUsersList.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")

	ucpAuthRolesGet.Flags().StringVar(&name, "rolename", "", "Name of the role retrieve")
	ucpAuthRolesCreate.Flags().StringVar(&name, "rolename", "", "Name of the role to create")
	ucpAuthRolesCreate.Flags().StringVar(&ruleset, "ruleset", "", "Path to a ruleset (JSON) to be used")
	ucpAuthRolesCreate.Flags().BoolVar(&admin, "service", false, "New role is a service account")

	// UCP ORG
	ucpAuth.AddCommand(ucpAuthOrg)
	ucpAuthOrg.AddCommand(ucpAuthOrgCreate)
	ucpAuthOrg.AddCommand(ucpAuthOrgDelete)
	ucpAuthOrg.AddCommand(ucpAuthOrgList)

	// TODO - UCP TEAMS
	ucpAuth.AddCommand(ucpAuthTeams)

	// UCP ROLES
	ucpAuth.AddCommand(ucpAuthRoles)
	ucpAuthRoles.AddCommand(ucpAuthRolesList)
	ucpAuthRoles.AddCommand(ucpAuthRolesGet)
	ucpAuthRoles.AddCommand(ucpAuthRolesCreate)

	// UCP USERS
	ucpAuth.AddCommand(ucpAuthUsers)
	ucpAuthUsers.AddCommand(ucpAuthUsersCreate)
	ucpAuthUsers.AddCommand(ucpAuthOrgDelete)
	ucpAuthUsers.AddCommand(ucpAuthUsersList)

	// UCP Grants
	ucpAuth.AddCommand(ucpAuthGrants)
	ucpAuthGrants.AddCommand(ucpAuthGrantsSet)
	ucpAuthGrants.AddCommand(ucpAuthGrantsGet)
	ucpAuthGrants.AddCommand(ucpAuthGrantsList)
	ucpAuthGrantsList.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")
	ucpAuthGrantsList.Flags().BoolVar(&resolve, "resolve", false, "Resolve the UUIDs to subject,role and grant names")
	ucpAuthGrantsSet.Flags().StringVar(&name, "subject", "", "The subject (user/org) that will be used")
	ucpAuthGrantsSet.Flags().StringVar(&ruleset, "role", "", "The role providing the capabilites")
	ucpAuthGrantsSet.Flags().StringVar(&collection, "collection", "", "The collection that will user")

	// UCP ROOT
	UCPRoot.AddCommand(ucpAuth)

}

var ucpAuth = &cobra.Command{
	Use:   "auth",
	Short: "Authorisation commands for users, groups and teams",
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
			log.Info("Importing Accounts from file")
			client, err := ucp.ReadToken()
			if err != nil {
				// Fatal error if can't read the token
				log.Fatalf("%v", err)
			}
			log.Debugf("Started parsing [%s]", importPath)
			err = client.ImportAccountsFromCSV(importPath)
			if err != nil {
				log.Fatalf("%v", err)
			}
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

		//orgs, err := client.GetAllOrgs()
		var accountQuery ucp.Account
		accountQuery.IsOrg = true

		orgs, err := client.GetAccounts(accountQuery, 1000)

		if err != nil {
			err = ucp.ParseUCPError([]byte(err.Error()))
			if err != nil {
				log.Errorf("Error parsing UCP error: %v", err)
			}
			log.Fatalf("%v", err)
		}

		if len(orgs.Accounts) == 0 {
			log.Error("No accounts returned")
			return
		}
		log.Debugf("Found %d Accounts", len(orgs.Accounts))
		fmt.Printf("Org Name\tFullname\n")
		for _, acct := range orgs.Accounts {
			fmt.Printf("%s\t%s\n", acct.Name, acct.FullName)
		}
	},
}

var ucpAuthUsers = &cobra.Command{
	Use:   "users",
	Short: "Manage Docker EE Users",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
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
			// Fatal error if can't read the token
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
			// Fatal error if can't read the token
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
		var users *ucp.AccountList
		var accountQuery ucp.Account
		accountQuery.IsOrg = false
		if admin {
			accountQuery.IsAdmin = true
		}
		if inactive {
			accountQuery.IsActive = false
		}
		users, err = client.GetAccounts(accountQuery, 1000)
		if err != nil {
			err = ucp.ParseUCPError([]byte(err.Error()))
			if err != nil {
				log.Errorf("Error parsing UCP error: %v", err)
			}
			log.Fatalf("%v", err)
		}

		if len(users.Accounts) == 0 {
			log.Error("No accounts returned")
			return
		}
		log.Debugf("Found %d Accounts", len(users.Accounts))
		fmt.Printf("User Name\tFullname\n")
		for _, acct := range users.Accounts {
			// Not sure why we're still retrieving ORGs even though we said false above - TODO
			if !acct.IsOrg {
				fmt.Printf("%s\t%s\n", acct.Name, acct.FullName)
			}
		}

	},
}

var ucpAuthTeams = &cobra.Command{
	Use:   "teams",
	Short: "Manage Docker EE Teams",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
	},
}

var ucpAuthRoles = &cobra.Command{
	Use:   "roles",
	Short: "Manage Docker EE Roles",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpAuthRolesList = &cobra.Command{
	Use:   "list",
	Short: "List Docker EE Roles",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = client.GetRoles()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var ucpAuthRolesGet = &cobra.Command{
	Use:   "get",
	Short: "List all rules for a particular role",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if name == "" {
			cmd.Help()
			log.Fatalln("No role specified to download")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		rules, err := client.GetRoleRuleset(name)
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("%s", rules)
	},
}

var ucpAuthRolesCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a new role based upon a ruleset",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if name == "" {
			cmd.Help()
			log.Fatalln("No role specified to download")
		}

		rulefile, err := ioutil.ReadFile(ruleset)
		if err != nil {
			log.Fatalf("%v", err)
		}

		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}

		err = client.CreateRole(name, name, string(rulefile), admin)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Infof("Role [%s] created succesfully", name)
	},
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
			log.Fatalf("%v", err)
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
			log.Fatalln("No role specified to download")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		rules, err := client.GetRoleRuleset(name)
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

		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}

		err = client.SetGrant(collection, ruleset, name, ucp.GrantCollection)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Infof("Grant created succesfully", name)
	},
}
