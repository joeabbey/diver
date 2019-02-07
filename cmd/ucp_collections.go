package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/joeabbey/diver/pkg/ucp/types"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/joeabbey/diver/pkg/ucp"
)

var newCollection ucptypes.Collection
var newConstraint ucptypes.CollectionLabelConstraints
var labelConstraints []string
var parentID string

func init() {
	ucpCollectionsGet.Flags().StringVar(&name, "id", "", "The ID of collection to inspect")

	ucpCollectionsCreate.Flags().StringVar(&name, "name", "", "Name of new collection")
	ucpCollectionsCreate.Flags().StringVar(&parentID, "parent", "", "The ID of the parent collection")

	ucpCollectionsDelete.Flags().StringVar(&name, "id", "", "ID of the collection to delete")

	ucpCollectionsSet.Flags().StringVar(&name, "id", "", "The ID of the collection to update")
	ucpCollectionsSet.Flags().StringVar(&newConstraint.LabelKey, "key", "", "The label Key")
	ucpCollectionsSet.Flags().StringVar(&newConstraint.LabelValue, "value", "", "The label value")
	ucpCollectionsSet.Flags().StringVar(&newConstraint.Type, "type", "", "Type is either a \"node\" or \"engine\" constraint")
	ucpCollectionsSet.Flags().BoolVar(&newConstraint.Equality, "equals", true, "The constraint is that the key \"equals\" the value")

	ucpCollectionsSetDefault.Flags().StringVar(&collection, "collection", "", "ID of the collection to add a user to")
	ucpCollectionsSetDefault.Flags().StringVar(&id, "user", "", "ID of the User to add to a collection")

	ucpAuth.AddCommand(ucpCollections)

	ucpCollections.AddCommand(ucpCollectionsList)
	ucpCollections.AddCommand(ucpCollectionsGet)

	if !DiverRO {
		ucpCollections.AddCommand(ucpCollectionsSet)
		ucpCollections.AddCommand(ucpCollectionsCreate)
		ucpCollections.AddCommand(ucpCollectionsDelete)
		ucpCollections.AddCommand(ucpCollectionsSetDefault)

	}
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
			log.Errorf("Unable to create collection [%s]", name)
			log.Fatalf("%v", err)
		}
		log.Infof("Succesfully created collection [%s] in parent collection [%s]", name, parentID)
	},
}

var ucpCollectionsDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Docker EE Collection",
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
		err = client.DeleteCollection(name)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Infof("Succesfully deleted collection [%s]", name)

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

var ucpCollectionsSet = &cobra.Command{
	Use:   "set",
	Short: "Set configuration details about a Docker EE Collection",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if name == "" {
			cmd.Help()
			log.Fatalf("No collection ID specified")
		}
		if newConstraint.LabelKey == "" {
			cmd.Help()
			log.Fatalf("No constraint key specified")
		}
		if newConstraint.LabelValue == "" {
			cmd.Help()
			log.Fatalf("No constraint value specified")
		}
		// Only two types currently exist as part of the collection types
		if newConstraint.Type != "node" && newConstraint.Type != "engine" {
			cmd.Help()
			log.Fatalf("Unknown Constraint type [%s]", newConstraint.Type)
		}
		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		err = client.SetCollection(name, &newConstraint)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Infof("Succesfully updated collection [%s]", name)
	},
}

var ucpCollectionsSetDefault = &cobra.Command{
	Use:   "default",
	Short: "Set the default collection for a user",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if collection == "" {
			cmd.Help()
			log.Fatalf("No collection specified")
		}
		if id == "" {
			cmd.Help()
			log.Fatalf("No username/id specified")
		}
		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.SetDefaultCollection(collection, id)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Infof("Succesfully set user [%s] to use collection [%s]", id, collection)
	},
}
