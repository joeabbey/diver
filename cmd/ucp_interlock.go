package cmd

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/joeabbey/diver/pkg/ucp"
	"github.com/joeabbey/diver/pkg/ucp/types"
)

var interlockConfig ucptypes.InterlockConfig

func init() {
	ucpInterlock.Flags().IntVar(&interlockConfig.HTTPPort, "http", 80, "HTTP Port")
	ucpInterlock.Flags().IntVar(&interlockConfig.HTTPSPort, "https", 8443, "HTTPS (TLS) Port")
	ucpInterlock.Flags().StringVar(&interlockConfig.Arch, "arch", "x86_64", "Interlock Architecture")
	ucpInterlock.Flags().BoolVar(&interlockConfig.InterlockEnabled, "enabled", true, "--enabled=true/false will enable or disable Interlock")

	UCPRoot.AddCommand(ucpInterlock)

}

var ucpInterlock = &cobra.Command{
	Use:   "interlock",
	Short: "Configure Interlock Layer 7 Routing",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		if strings.HasPrefix(client.UCPVersion, "ucp/2") {
			err = client.ConfigureHRM(interlockConfig)
		} else {
			err = client.ConfigureInterlock(interlockConfig)
		}
		if err != nil {
			// Fatal error if can't configure Interlock
			log.Fatalf("%v", err)

		}
	},
}
