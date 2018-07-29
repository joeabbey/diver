package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

var username string

func init() {

	ucpCliBundleList.Flags().StringVar(&username, "username", "", "Username to list bundle for (default: logged-in user)")

	ucpCliBundleDelete.Flags().StringVar(&username, "username", "", "Username to delete bundle for (default: logged-in user)")

	ucpCliBundleRename.Flags().StringVar(&username, "username", "", "Username to rename bundle for (default: logged-in user)")

	ucpCliBundle.AddCommand(ucpCliBundleDownload)
	ucpCliBundle.AddCommand(ucpCliBundleList)

	if !DiverRO {
		ucpCliBundle.AddCommand(ucpCliBundleRename)
		ucpCliBundle.AddCommand(ucpCliBundleDelete)
	}

	UCPRoot.AddCommand(ucpCliBundle)
}

// ucpCliBundle kept as an alias for `diver ucp client-bundle download` to preserve compatibility, may be deprecated later?
var ucpCliBundle = &cobra.Command{
	Use:   "client-bundle",
	Short: "Download the client bundle for UCP",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.GetClientBundle()
		if err != nil {
			log.Fatalf("%v", err)
		}

	},
}

var ucpCliBundleDownload = &cobra.Command{
	Use:   "get",
	Short: "Download the client bundle for UCP",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.GetClientBundle()
		if err != nil {
			log.Fatalf("%v", err)
		}

	},
}

var ucpCliBundleList = &cobra.Command{
	Use:   "ls",
	Short: "List client bundles for UCP",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		clientBundles, err := client.ListClientBundles(username)
		if err != nil {
			log.Fatalf("%v", err)
		}

		const padding = 3
		w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
		fmt.Fprintln(w, "ID\tLabel")

		for _, publicKey := range clientBundles.AccountPublicKeys {
			fmt.Fprintf(w, "%s\t%s\n", publicKey.ID, publicKey.Label)
		}
		w.Flush()

	},
}

var ucpCliBundleRename = &cobra.Command{
	Use:   "rename BUNDLE_ID NEW_LABEL",
	Short: "Rename client bundle for UCP",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.RenameClientBundle(username, args[0], args[1])
		if err != nil {
			log.Fatalf("%v", err)
		}

	},
}

var ucpCliBundleDelete = &cobra.Command{
	Use:   "rm BUNDLE_ID",
	Short: "Delete client bundle for UCP",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.DeleteClientBundle(username, args[0])
		if err != nil {
			log.Fatalf("%v", err)
		}

	},
}
