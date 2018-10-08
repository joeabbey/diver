package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/joeabbey/diver/pkg/dtr"
)

var enable bool

func init() {
	dtrSettingsSet.PersistentFlags().BoolVar(&enable, "setting", true, "--setting=true/false will enable or disable this setting")
	dtrSettingsSet.AddCommand(dtrSettingsCreateRepo)
	dtrSettingsSet.AddCommand(dtrSettingsScanning)

}

var dtrSettingsCreateRepo = &cobra.Command{
	Use:   "createrepo",
	Short: "Create a repository on Push",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := dtr.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.DTRCreateRepoOnPush(enable)
		if err != nil {
			// Fatal error if can't return any webhooks
			log.Fatalf("%v", err)
		}
	},
}

var dtrSettingsScanning = &cobra.Command{
	Use:   "scanning",
	Short: "DTR Image scanning",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := dtr.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.DTRScanningEnable(enable)
		if err != nil {
			// Fatal error if can't return any webhooks
			log.Fatalf("%v", err)

		}
	},
}
