package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/dtr"
)

func init() {

	dtrSettings.AddCommand(dtrSettingsGet)
	dtrSettings.AddCommand(dtrSettingsSet)

	dtrCmd.AddCommand(dtrSettings)

}

var dtrSettings = &cobra.Command{
	Use:   "settings",
	Short: "Interact with Docker Trusted Registry Settings",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var dtrSettingsGet = &cobra.Command{
	Use:   "get",
	Short: "Get settings values from the Docker Trusted Registry",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := dtr.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		settings, err := client.DTRGetSettings()
		if err != nil {
			// Fatal error if can't return any webhooks
			log.Fatalf("%v", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Configuration\tSetting")

		fmt.Fprintf(w, "DTR Host\t%s\n", settings.DtrHost)
		fmt.Fprintf(w, "DTR Replica ID\t%s\n", settings.ReplicaID)
		fmt.Fprintf(w, "SSO\t%t\n", settings.Sso)
		fmt.Fprintf(w, "Create Repo on Push\t%t\n", settings.CreateRepositoryOnPush)
		fmt.Fprintf(w, "Log Protocol\t%s\n", settings.LogProtocol)
		fmt.Fprintf(w, "Log Host\t%s\n", settings.LogHost)
		fmt.Fprintf(w, "Log Level\t%s\n", settings.LogLevel)
		fmt.Fprintf(w, "Storage Volume\t%s\n", settings.StorageVolume)
		fmt.Fprintf(w, "NFS Host\t%s\n", settings.NfsHost)
		fmt.Fprintf(w, "NFS Path\t%s\n", settings.NfsPath)
		fmt.Fprintf(w, "Scanning Enabled\t%t\n", settings.ScanningEnabled)
		fmt.Fprintf(w, "Scanning Online Sync\t%t\n", settings.ScanningSyncOnline)
		fmt.Fprintf(w, "Scanning Auto Recheck\t%t\n", settings.ScanningEnableAutoRecheck)
		w.Flush()
	},
}

var dtrSettingsSet = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values in Docker Trusted Registry",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
