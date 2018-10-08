package cmd

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/joeabbey/diver/pkg/ucp"
)

// This is the ucp version that we will attempt to upgrade to
var ucpversion string

func init() {
	ucpUpgrade.Flags().StringVar(&ucpversion, "version", "", "The version of UCP to upgrade to")

	// UCP ROOT
	UCPRoot.AddCommand(ucpUpgrade)
}

var ucpUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade Docker Universal Control Plane",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}

		// Retrieve the current UCP version
		v, err := client.GetUCPVersion()
		if err != nil {
			log.Fatalf("%v", err)
		}

		// If the user doesn't specify a particular version then print out current and available versions with the help
		if ucpversion == "" {
			cmd.Help()
			versions, err := client.GetAvailavleUCPVersions()
			if err != nil {
				log.Fatalf("%v", err)
			}
			fmt.Printf("\nAvailable Versions\n")
			for i := range versions {

				fmt.Printf("%s\n", versions[i])
			}
			fmt.Printf("\n")
			fmt.Printf("Current Version\n%s\n\n", v)

			log.Fatalln("No --version specified")
		}

		log.Infof("Upgrading from [%s] to [%s]", v, ucpversion)
		log.Warnln("The Universal Control Plane will be un-available for > 5 minutes whilst the upgrade takes place")

		// Give user the opportunity to cancel the procedure
		log.Infoln("The procedure will begin in 10 seconds, press ctrl+c to cancel")
		time.Sleep(time.Second * 10)

		// Begin the upgrade procedure
		err = client.UpgradeUCP(ucpversion)
		if err != nil {
			log.Fatalf("%v", err)
		}

		log.Infoln("Upgrade procedure has begun succesfully")

	},
}
