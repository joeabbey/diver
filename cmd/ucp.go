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

var filepath, action string

var top bool

func init() {
	diverCmd.AddCommand(ucpRoot)
	//client := ucp.Client{}

	ucpRoot.Flags().StringVar(&client.Username, "username", os.Getenv("DIVER_USERNAME"), "Username that has permissions to authenticate to Docker EE")
	ucpRoot.Flags().StringVar(&client.Password, "password", os.Getenv("DIVER_PASSWORD"), "Password allowing a user to authenticate to Docker EE")
	ucpRoot.Flags().StringVar(&client.UCPURL, "url", os.Getenv("DIVER_URL"), "URL for Docker EE, e.g. https://10.0.0.1")
	ignoreCert := strings.ToLower(os.Getenv("DIVER_INSECURE")) == "true"

	ucpRoot.Flags().BoolVar(&client.IgnoreCert, "ignorecert", ignoreCert, "Ignore x509 certificate")

	ucpRoot.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")

	// Add subcommands
	ucpAuth.Flags().StringVar(&auth.FullName, "fullname", "", "The full name of a UCP user or organisation")
	ucpAuth.Flags().StringVar(&auth.Name, "username", "", "The unique username organisation")
	ucpAuth.Flags().StringVar(&auth.Password, "password", "", "A string password for a new user of organisation")
	ucpAuth.Flags().BoolVar(&auth.IsAdmin, "admin", false, "Make this user an administrator")
	ucpAuth.Flags().BoolVar(&auth.IsActive, "active", true, "Enable this user in the Universal Control Plane")
	ucpAuth.Flags().BoolVar(&auth.IsOrg, "isorg", false, "Create an Organisation")
	ucpAuth.Flags().StringVar(&filepath, "file", "", "Read users from a file [csv currently supported]")
	ucpAuth.Flags().StringVar(&action, "action", "create", "Action to be performed [create/delete/update]")
	ucpAuth.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")
	ucpRoot.AddCommand(ucpAuth)

	ucpContainer.Flags().IntVar(&logLevel, "logLevel", 4, "Set the logging level [0=panic, 3=warning, 5=debug]")

	ucpContainer.Flags().BoolVar(&top, "top", false, "Enable TOP for watching running containers")
	ucpRoot.AddCommand(ucpContainer)
	ucpRoot.AddCommand(ucpCliBundle)
}

var ucpRoot = &cobra.Command{
	Use:   "ucp",
	Short: "Universal Control Plane ",
	Run: func(cmd *cobra.Command, args []string) {

		// Error checking flags/variables
		if client.Username == "" {
			log.Errorln("UCP Username is required")
			cmd.Help()
		}
		if client.Password == "" {
			log.Errorln("UCP Password is required")
			cmd.Help()
		}
		if client.UCPURL == "" {
			log.Errorln("UCP URL is required [https://<address/]")
			cmd.Help()
		}

		log.SetLevel(log.Level(logLevel))
		err := client.Connect()
		if err != nil {
			log.Errorf("%v", err)
		} else {
			// err = client.ListNetworks()
			// if err != nil {
			// 	log.Errorf("%v\n", err)
			// }
			// err = client.ListContainerJSON()
			// if err != nil {
			// 	log.Errorf("%v\n", err)
			// }
			// err = client.GetClientBundle()
			// if err != nil {
			// 	log.Errorf("%v\n", err)
			// }

			// user := ucp.NewUser("dan finneran", "dan", "password", true, true, false)
			// err = client.AddAccount(user)
			// if err != nil {
			// 	log.Errorf("%v\n", err)
			// }
			// err = client.DeleteAccount("dan1")
			// if err != nil {
			// 	log.Errorf("%v\n", err)
			// }

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

		// if action == "" {
		// 	cmd.Help()
		// 	log.Fatalf("--action is a required flag")
		// }
		// A file has been passed in, so parse it and return
		if filepath != "" {
			_, err := ucp.ReadToken()
			if err != nil {
				// Fatal error if can't read the torken
				log.Fatalf("%v", err)
			}
			log.Debugf("Started parsing [%s]", filepath)

		} else {
			// Parse flags/variables

			client, err := ucp.ReadToken()
			if err != nil {
				// Fatal error if can't read the torken
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
				// Fatal error if can't read the torken
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

		if top == true {
			client, err := ucp.ReadToken()
			if err != nil {
				// Fatal error if can't read the torken
				log.Fatalf("%v", err)
			}
			err = client.ContainerTop()
			if err != nil {
				// Fatal error if can't read the torken
				log.Fatalf("%v", err)
			}
		}
	},
}

var ucpCliBundle = &cobra.Command{
	Use:   "client-bundle",
	Short: "download the client bundle for UCP",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the torken
			log.Fatalf("%v", err)
		}
		err = client.GetClientBundle()
		if err != nil {
			// Fatal error if can't read the torken
			log.Fatalf("%v", err)
		}

	},
}
