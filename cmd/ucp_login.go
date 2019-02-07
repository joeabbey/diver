package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/joeabbey/diver/pkg/ucp"
)

var session int

func init() {

	ucpLogin.Flags().StringVar(&ucpClient.Username, "username", os.Getenv("UCP_USERNAME"), "Username that has permissions to authenticate to Docker EE")
	ucpLogin.Flags().StringVar(&ucpClient.Password, "password", os.Getenv("UCP_PASSWORD"), "Password allowing a user to authenticate to Docker EE")
	ucpLogin.Flags().StringVar(&ucpClient.UCPURL, "url", os.Getenv("UCP_URL"), "URL for Docker EE, e.g. https://10.0.0.1")
	ignoreCert := strings.ToLower(os.Getenv("UCP_INSECURE")) == "true"

	ucpLogin.Flags().BoolVar(&ucpClient.IgnoreCert, "ignorecert", ignoreCert, "Ignore x509 certificate")

	ucpLoginSet.Flags().IntVar(&session, "id", 0, "The session ID to set to active")

	ucpLogin.AddCommand(ucpLoginList)
	ucpLogin.AddCommand(ucpLoginSet)

	UCPRoot.AddCommand(ucpLogin)
}

// UCPLogin - This manages logging in and swapping of contexts
var ucpLogin = &cobra.Command{
	Use:   "login",
	Short: "Authenticate against the Universal Control Pane",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		err := ucpClient.Connect()

		// Check if connection was succesful
		if err != nil {
			cmd.Help()
			log.Fatalf("%v", err)
		} else {
			// If succesfull write the token and annouce as succesful
			err = ucpClient.WriteToken()
			if err != nil {
				log.Errorf("%v", err)
			}
			log.Infof("Succesfully logged into [%s]", ucpClient.UCPURL)
		}
	},
}

// ucpLoginList - This manages logging in and swapping of contexts
var ucpLoginList = &cobra.Command{
	Use:   "list",
	Short: "List all UCP sessions",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		// Retrieve all of the client sessions in the token file
		clientTokens, err := ucp.ReadAllClients()
		if err != nil {
			log.Errorf("%v", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "ID\tAddress\tActive")
		for i := range clientTokens {
			fmt.Fprintf(w, "%d\t%s\t%t\n", i, clientTokens[i].UCPAddress, clientTokens[i].Active)
		}
		w.Flush()
	},
}

// ucpLoginSet - This manages logging in and swapping of contexts
var ucpLoginSet = &cobra.Command{
	Use:   "setActive",
	Short: "Set the active UCP session",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		err := ucp.SetActiveSession(session)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Infof("Set session [%d] to active", session)
	},
}
