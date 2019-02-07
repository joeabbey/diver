package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/joeabbey/diver/pkg/ucp"
)

func init() {

	ucpSupport.AddCommand(ucpSupportDownload)

	UCPRoot.AddCommand(ucpSupport)
}

var ucpSupport = &cobra.Command{
	Use:   "support",
	Short: "Download the client bundle for UCP",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpSupportDownload = &cobra.Command{
	Use:   "get",
	Short: "Download the support dump for UCP",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.GetSupportDump()
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}
