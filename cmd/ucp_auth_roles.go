package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

func init() {
	// UCP ROLES flags
	ucpAuthRolesGet.Flags().StringVar(&name, "rolename", "", "Name of the role to retrieve")
	ucpAuthRolesGet.Flags().StringVar(&id, "id", "", "ID of the role to retrieve")

	ucpAuthRolesCreate.Flags().StringVar(&name, "rolename", "", "Name of the role to create")
	ucpAuthRolesCreate.Flags().StringVar(&ruleset, "ruleset", "", "Path to a ruleset (JSON) to be used")
	ucpAuthRolesCreate.Flags().BoolVar(&admin, "service", false, "New role is a system role")

	// UCP ROLES
	ucpAuth.AddCommand(ucpAuthRoles)
	ucpAuthRoles.AddCommand(ucpAuthRolesList)
	ucpAuthRoles.AddCommand(ucpAuthRolesTotal)
	if !DiverRO {
		ucpAuthRoles.AddCommand(ucpAuthRolesGet)
		ucpAuthRoles.AddCommand(ucpAuthRolesCreate)
	}
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
		if name == "" && id == "" {
			cmd.Help()
			log.Fatalln("No role specified to download")
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

var ucpAuthRolesTotal = &cobra.Command{
	Use:   "totalrole",
	Short: "returns the TOTAL ruleset",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		rules, err := client.TotalRole()
		if err != nil {
			log.Fatalf("%v", err)
		}
		var prettyJSON bytes.Buffer
		err = json.Indent(&prettyJSON, rules, "", "\t")
		if err != nil {
			log.Fatalf("%v", err)
		}
		fmt.Printf("%s", prettyJSON.Bytes())
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
