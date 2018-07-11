package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/thebsdbox/diver/pkg/ucp/types"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

var newCollection ucptypes.Collection
var labelConstraints []string
var parentID string

func init() {
	ucpCollectionsGet.Flags().StringVar(&name, "id", "", "The ID of collection to inspect")

	ucpCollectionsCreate.Flags().StringVar(&name, "name", "", "Name of new collection")
	ucpCollectionsCreate.Flags().StringVar(&parentID, "parent", "", "The ID of the parent collection")

	ucpAuth.AddCommand(ucpCollections)

	ucpCollections.AddCommand(ucpCollectionsCreate)
	ucpCollections.AddCommand(ucpCollectionsList)
	ucpCollections.AddCommand(ucpCollectionsGet)

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

var ucpCollectionsCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a Docker EE Collection",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if name == "" {
			cmd.Help()
			log.Fatalf("No collection ID specified")
		}
		if parentID == "" {
			cmd.Help()
			log.Fatalf("No parent collection ID specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = client.CreateCollection(name, parentID)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var ucpCollectionsGet = &cobra.Command{
	Use:   "get",
	Short: "Get information about a Docker EE Collection",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if name == "" {
			cmd.Help()
			log.Fatalf("No collection ID specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		collection, err := client.GetCollection(name)
		if collection == nil {
			log.Fatalf("Error parsing collection [%s]", name)
		}
		if err != nil {
			log.Fatalf("%v", err)
		}
		const padding = 3
		w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

		fmt.Fprintf(w, "Name\t%s\n", collection.Name)
		fmt.Fprintf(w, "ID\t%s\n", collection.ID)
		fmt.Fprintf(w, "Path\t%s\n", collection.Path)
		fmt.Fprintf(w, "Parents\t\n")

		for i := range collection.ParentIds {
			fmt.Fprintf(w, "\t/%s\n", collection.ParentIds[i])
		}
		fmt.Fprintf(w, "Label Constraints\t\n")
		fmt.Fprintf(w, "\tKey\tValue\tEqual to\tType\n")
		for i := range collection.LabelConstraints {
			fmt.Fprintf(w, "\t%s\t%s\t%t\t%s\n", collection.LabelConstraints[i].LabelKey, collection.LabelConstraints[i].LabelValue, collection.LabelConstraints[i].Equality, collection.LabelConstraints[i].Type)
		}
		w.Flush()
	},
}
