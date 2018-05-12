package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"

	log "github.com/Sirupsen/logrus"
)

var logLevel = 5
var client ucp.Client
var auth ucp.Account

var importPath, exportPath, action string

var top, exampleFile bool

func init() {
	diverCmd.AddCommand(ucpRoot)

	ucpRoot.Flags().StringVar(&client.Username, "username", os.Getenv("DIVER_USERNAME"), "Username that has permissions to authenticate to Docker EE")
	ucpRoot.Flags().StringVar(&client.Password, "password", os.Getenv("DIVER_PASSWORD"), "Password allowing a user to authenticate to Docker EE")
	ucpRoot.Flags().StringVar(&client.UCPURL, "url", os.Getenv("DIVER_URL"), "URL for Docker EE, e.g. https://10.0.0.1")
	ignoreCert := strings.ToLower(os.Getenv("DIVER_INSECURE")) == "true"

	ucpRoot.Flags().BoolVar(&client.IgnoreCert, "ignorecert", ignoreCert, "Ignore x509 certificate")

	ucpRoot.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")

	// Auth Flags
	ucpAuth.Flags().StringVar(&auth.FullName, "fullname", "", "The full name of a UCP user or organisation")
	ucpAuth.Flags().StringVar(&auth.Name, "username", "", "The unique username organisation")
	ucpAuth.Flags().StringVar(&auth.Password, "password", "", "A string password for a new user of organisation")
	ucpAuth.Flags().BoolVar(&auth.IsAdmin, "admin", false, "Make this user an administrator")
	ucpAuth.Flags().BoolVar(&auth.IsActive, "active", true, "Enable this user in the Universal Control Plane")
	ucpAuth.Flags().BoolVar(&auth.IsOrg, "isorg", false, "Create an Organisation")
	ucpAuth.Flags().StringVar(&importPath, "importCSV", "", "Import accounts from a file [csv currently supported]")
	ucpAuth.Flags().StringVar(&exportPath, "exportCSV", "", "Export users to a file [csv currently supported]")

	ucpAuth.Flags().BoolVar(&exampleFile, "exampleCSV", false, "Create an example csv file [example_accounts.csv]")

	ucpAuth.Flags().StringVar(&action, "action", "create", "Action to be performed [create/delete/update]")
	ucpAuth.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")

	// Container flags
	ucpContainer.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")
	ucpContainer.Flags().BoolVar(&top, "top", false, "Enable TOP for watching running containers")

	ucpRoot.AddCommand(ucpAuth)
	ucpRoot.AddCommand(ucpContainer)
	ucpRoot.AddCommand(ucpCliBundle)
	ucpRoot.AddCommand(ucpNetwork)
	ucpContainer.AddCommand(ucpContainerTop)
	ucpContainer.AddCommand(ucpContainerList)

}

var ucpRoot = &cobra.Command{
	Use:   "ucp",
	Short: "Universal Control Plane ",
	Run: func(cmd *cobra.Command, args []string) {

		existingClient, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Warn("Unable to find existing session, please login")
		} else {

			currentAccount, err := existingClient.AuthStatus()
			if err != nil {
				log.Errorf("%v", err)
			} else {

				log.Infof("Current user [%s]", currentAccount.Name)
				return
			}
		}
		// Error checking flags/variables
		if client.Username == "" {
			cmd.Help()
			log.Fatalln("UCP Username is required")

		}
		if client.Password == "" {
			cmd.Help()
			log.Fatalln("UCP Password is required")
		}
		if client.UCPURL == "" {
			cmd.Help()
			log.Fatalln("UCP URL is required [https://<address/]")
		}

		log.SetLevel(log.Level(logLevel))
		err = client.Connect()

		// Check if connection was succesful
		if err != nil {
			log.Errorf("%v", err)
		} else {
			// If succesfull write the token and annouce as succesful
			err = client.WriteToken()
			if err != nil {
				log.Errorf("%v", err)
			}
			log.Infof("Succesfully logged into [%s]", client.UCPURL)
		}
	},
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
			os.Exit(0)
		} else {
			// Parse flags/variables
			if auth.Name == "" {
				log.Fatalf("No Username has been entered")
			}
			client, err := ucp.ReadToken()
			if err != nil {
				log.Fatalf("%v", err)
			}

			switch action {
			case "create":
				err = client.AddAccount(&auth)
			case "delete":
				err = client.DeleteAccount(auth.Name)
			case "update":
				log.Errorf("Not supported (yet)")
			default:
				log.Errorf("Unknown action [%s]", action)
				cmd.Help()
			}

			if err != nil {
				// Fatal error if can't read the token
				log.Fatalf("%v", err)
			}
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

var ucpCliBundle = &cobra.Command{
	Use:   "client-bundle",
	Short: "download the client bundle for UCP",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.GetClientBundle()
		if err != nil {
			log.Fatalf("%v", err)
		}

	},
}

var ucpNetwork = &cobra.Command{
	Use:   "network",
	Short: "Interact with container networks",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		err = client.GetNetworks()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}
