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

	ucpAuth.AddCommand(ucpCollections)
	ucpCollections.AddCommand(ucpCollectionsList)
}

var ucpCollections = &cobra.Command{
	Use:   "collections",
	Short: "Manage Docker EE Collections",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		cmd.Help()
	},
}

var ucpCollectionsList = &cobra.Command{
	Use:   "list",
	Short: "List Docker EE Collections",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		collections, err := client.GetCollections()
		const padding = 3
		w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

		fmt.Fprintf(w, "Name\tID\n")

		for i := range collections {
			fmt.Fprintf(w, "%s\t%s\n", collections[i].Name, collections[i].ID)
		}
		w.Flush()
	},
}
