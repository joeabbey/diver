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
	diverCmd.AddCommand(UCPRoot)

	ucpLogin.Flags().StringVar(&ucpClient.Username, "username", os.Getenv("UCP_USERNAME"), "Username that has permissions to authenticate to Docker EE")
	ucpLogin.Flags().StringVar(&ucpClient.Password, "password", os.Getenv("UCP_PASSWORD"), "Password allowing a user to authenticate to Docker EE")
	ucpLogin.Flags().StringVar(&ucpClient.UCPURL, "url", os.Getenv("UCP_URL"), "URL for Docker EE, e.g. https://10.0.0.1")
	ignoreCert := strings.ToLower(os.Getenv("UCP_INSECURE")) == "true"

	ucpLogin.Flags().BoolVar(&ucpClient.IgnoreCert, "ignorecert", ignoreCert, "Ignore x509 certificate")

	// Container flags
	ucpContainer.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")
	ucpContainer.Flags().BoolVar(&top, "top", false, "Enable TOP for watching running containers")

	UCPRoot.AddCommand(ucpContainer)
	UCPRoot.AddCommand(ucpLogin)

	// Sub commands
	ucpContainer.AddCommand(ucpContainerTop)
	ucpContainer.AddCommand(ucpContainerList)

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

var ucpContainer = &cobra.Command{
	Use:   "containers",
	Short: "Interact with containers",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
	},
}

var ucpContainerTop = &cobra.Command{
	Use:   "top",
	Short: "A list of containers and their CPU usage like the top command on linux",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.ContainerTop()
		if err != nil {
			log.Fatalf("%v", err)
		}
		return
	},
}

var ucpContainerList = &cobra.Command{
	Use:   "list",
	Short: "List all containers across all nodes in UCP",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.GetContainerNames()
		if err != nil {
			log.Fatalf("%v", err)
		}
		return
	},
}
