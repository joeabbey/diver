package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

func init() {

	// UCP ORG Flags
	//ucpAuthOrgCreate.Flags().StringVar(&auth.Name, "name", "", "A unique Organisation name")

	//ucpAuthOrgDelete.Flags().StringVar(&auth.Name, "name", "", "Existing Organisation")
	ucpConfig.AddCommand(ucpConfigList)
	UCPRoot.AddCommand(ucpConfig)
}

var ucpConfig = &cobra.Command{
	Use:   "config",
	Short: "Manage Docker EE Service Configurations",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		cmd.Help()
	},
}

var ucpConfigList = &cobra.Command{
	Use:   "list",
	Short: "List all Docker EE Service Configurations",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		cfgs, err := client.ListConfigs()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		log.Debugf("Found [%d] configurations", len(cfgs))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintf(w, "Name\tID\tVersion\n")
		for i := range cfgs {
			fmt.Fprintf(w, "%s\t%s\t%d\n", cfgs[i].Spec.Name, cfgs[i].ID, cfgs[i].Version.Index)
		}
		w.Flush()

	},
}
