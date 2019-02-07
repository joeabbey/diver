package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/joeabbey/diver/pkg/ucp"
	"github.com/joeabbey/diver/pkg/ucp/types"
)

// This are the configurable options for LDAP
var usernameLDAP, passwordLDAP, serverLDAP, baseDN, usernameAttribute string

func init() {
	ucpLDAPSet.Flags().StringVar(&serverLDAP, "server", "", "The hostname or IP address of the LDAP server")
	ucpLDAPSet.Flags().StringVar(&serverLDAP, "username", "", "A username with read privileges for the LDAP server")
	ucpLDAPSet.Flags().StringVar(&serverLDAP, "password", "", "A password for the LDAP server")
	ucpLDAPSet.Flags().StringVar(&baseDN, "basedn", "", "The base Distinguished Name to search from e.g. dc=ad,dc=company,dc=com")
	ucpLDAPSet.Flags().StringVar(&usernameAttribute, "userAttribute", "", "A password for the LDAP server")

	ucpLDAP.AddCommand(ucpLDAPSet)
	ucpLDAP.AddCommand(ucpLDAPGet)
	UCPRoot.AddCommand(ucpLDAP)
}

var ucpLDAP = &cobra.Command{
	Use:   "ldap",
	Short: "Manage UCP and LDAP",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpLDAPGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve the LDAP configuration",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		cfg, err := client.GetLDAPInfo()
		if err != nil {
			log.Fatalf("%v", err)
		}
		printLDAPCfg(cfg)

	},
}

// ucpLDAPSet will just configure the basics in the initial release
// TODO - Add additional configuration abilities
var ucpLDAPSet = &cobra.Command{
	Use:   "set",
	Short: "Configure the LDAP configuration",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		// Retrieve the existing configuration as we'll need to update it
		existingCfg, err := client.GetLDAPInfo()
		if err != nil {
			log.Fatalf("%v", err)
		}

		if serverLDAP == "" && existingCfg.ServerURL == "" {
			log.Fatalf("No LDAP server has been specified or configured")
		}

		// Update the Server URL
		if serverLDAP != "" {
			// Perform testing on the URL to ensure it matches the input requirements for the UCP API
			if strings.HasPrefix(serverLDAP, "ldap://") {
				existingCfg.ServerURL = serverLDAP
			} else if strings.HasPrefix(serverLDAP, "ldaps://") {
				existingCfg.ServerURL = serverLDAP
			} else {
				log.Fatalln("The LDAP server URL should begin with ldap:// or ldaps://")
			}
		}

		// Update or set the Username for the LDAP server
		if usernameLDAP != "" {
			existingCfg.ReaderDN = usernameLDAP
		}

		// Update or set the Password for the LDAP server
		if passwordLDAP != "" {
			existingCfg.ReaderPassword = passwordLDAP
		}

		// See if a previous userSearch config exists
		log.Debugf("Found [%d] search configurations", len(existingCfg.UserSearchConfigs))
		if len(existingCfg.UserSearchConfigs) == 1 {
			// If creating an entirely new Search Config

			if baseDN == "" || usernameAttribute == "" {
				log.Fatalln("Both --basedn and --userAttribute are needed for a new configuration")
			} else {
				log.Debugln("Updating base Search configuration")
				existingCfg.UserSearchConfigs[0].BaseDN = baseDN
				existingCfg.UserSearchConfigs[0].UserAttr = usernameAttribute
			}
		} else {
			if baseDN != "" || usernameAttribute != "" {
				log.Fatalln("Existing search configurations should be managed in the UI")
			}
		}
		err = client.SetLDAPInfo(existingCfg)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Infoln("Succesfully updated LDAP settings")
	},
}

func printLDAPCfg(cfg *ucptypes.LDAPConfig) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)

	fmt.Fprintf(w, "LDAP Server URL:\t%s\n", cfg.ServerURL)
	fmt.Fprintf(w, "LDAP Username:\t%s\n", cfg.ReaderDN)
	fmt.Fprintf(w, "Skip TLS:\t%t\n", cfg.TLSSkipVerify)
	fmt.Fprintf(w, "JIT User Provisioning:\t%t\n", cfg.JitUserProvisioning)

	w.Flush()

}
