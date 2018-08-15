package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
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

		if ucpversion == "" {
			cmd.Help()
			versions, err := client.GetAvailavleUCPVersions()
			if err != nil {
				log.Fatalf("%v", err)
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
			fmt.Fprintln(w, "\nAvailable Versions")
			for i := range versions {

				fmt.Fprintf(w, "%s\n", versions[i])
			}
			fmt.Printf("\n")
			w.Flush()

			log.Fatalln("No --version specified")
		}

		// Retrieve the current UCP version
		v, err := client.GetUCPVersion()
		if err != nil {
			log.Fatalf("%v", err)
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
