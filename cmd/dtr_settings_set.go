package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/dtr"
)

var enable bool

func init() {
	dtrSettingsCreateRepo.Flags().BoolVar(&enable, "enable", true, "Enable the creation of a repository on push")
	dtrSettingsSet.AddCommand(dtrSettingsCreateRepo)

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
