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
	ucpRoot.AddCommand(ucpAuth)

}

var ucpRoot = &cobra.Command{
	Use:   "ucp",
	Short: "Universal Control Plane ",
	Run: func(cmd *cobra.Command, args []string) {
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

			user := ucp.NewUser("dan finneran", "dan", "password", true, true, false)
			err = client.AddAccount(user)
			if err != nil {
				log.Errorf("%v\n", err)
			}
			err = client.DeleteAccount("dan1")
			if err != nil {
				log.Errorf("%v\n", err)
			}
		}
	},
}

var ucpAuth = &cobra.Command{
	Use:   "Auth",
	Short: "Authorisation commands for users, groups and teams",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
