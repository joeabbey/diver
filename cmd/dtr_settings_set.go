package cmd

import (
	"github.com/joeabbey/diver/pkg/dtr"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var enable bool

func init() {
	dtrSettingsSet.PersistentFlags().BoolVar(&enable, "setting", true, "--setting=true/false will enable or disable this setting")
	dtrSettingsSet.AddCommand(dtrSettingsCreateRepo)
	dtrSettingsSet.AddCommand(dtrSettingsScanning)
	dtrSettingsSet.AddCommand(dtrSettingsOnline)

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

var dtrSettingsOnline = &cobra.Command{
	Use:   "online",
	Short: "Scanning online sync",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := dtr.ReadToken()
		if err != nil {
			//can't read token
			log.Fatalf("%v", err)
		}
		err = client.DTROnlineScan(enable)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}
