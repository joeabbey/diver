package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"

	log "github.com/Sirupsen/logrus"
)

var ucpClient ucp.Client
var importPath, exportPath, action string
var top, exampleFile bool

func init() {
	ucpLogin.Flags().StringVar(&ucpClient.Username, "username", os.Getenv("UCP_USERNAME"), "Username that has permissions to authenticate to Docker EE")
	ucpLogin.Flags().StringVar(&ucpClient.Password, "password", os.Getenv("UCP_PASSWORD"), "Password allowing a user to authenticate to Docker EE")
	ucpLogin.Flags().StringVar(&ucpClient.UCPURL, "url", os.Getenv("UCP_URL"), "URL for Docker Universal Control Plane, e.g. https://10.0.0.1")
	ignoreCert := strings.ToLower(os.Getenv("UCP_INSECURE")) == "true"

	ucpLogin.Flags().BoolVar(&ucpClient.IgnoreCert, "ignorecert", ignoreCert, "Ignore x509 certificate")

	UCPRoot.AddCommand(ucpLogin)

	// Add UCP and subcommands to the main application
	diverCmd.AddCommand(UCPRoot)

}

// UCPRoot - This is the root of all UCP commands / flags
var UCPRoot = &cobra.Command{
	Use:   "ucp",
	Short: "Universal Control Plane ",
	Run: func(cmd *cobra.Command, args []string) {

		existingClient, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			cmd.Help()
			log.Warn("Unable to find existing session, please login")
			return
		}
		currentAccount, err := existingClient.AuthStatus()
		if err != nil {
			cmd.Help()
			log.Warn("Session has expired, please login")
			return
		}
		cmd.Help()
		fmt.Printf("\n\n")
		log.Infof("Current user [%s]", currentAccount.Name)
		return
	},
}

// UCPRoot - This is the root of all UCP commands / flags
var ucpLogin = &cobra.Command{
	Use:   "login",
	Short: "Authenticate against the Universal Control Pane",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		err := ucpClient.Connect()

		// Check if connection was succesful
		if err != nil {
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
