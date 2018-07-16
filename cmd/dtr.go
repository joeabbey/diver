package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/dtr"
)

var dtrClient dtr.Client

func init() {
	dtrLogin.Flags().StringVar(&dtrClient.Username, "username", os.Getenv("DTR_USERNAME"), "Username that has permissions to authenticate to Docker EE")
	dtrLogin.Flags().StringVar(&dtrClient.Password, "password", os.Getenv("DTR_PASSWORD"), "Password allowing a user to authenticate to Docker EE")
	dtrLogin.Flags().StringVar(&dtrClient.DTRURL, "url", os.Getenv("DTR_PASSWORD"), "The URL of a Docker Trusted Registry")

	ignoreCert := strings.ToLower(os.Getenv("STORE_INSECURE")) == "true"

	dtrLogin.Flags().BoolVar(&dtrClient.IgnoreCert, "ignorecert", ignoreCert, "Ignore x509 certificate")

	dtrCmd.AddCommand(dtrLogin)
	dtrCmd.AddCommand(dtrInfo)
	dtrInfo.AddCommand(dtrLoginReplicas)
	dtrWebHooks.AddCommand(dtrWebHooksList)
	dtrCmd.AddCommand(dtrWebHooks)

	diverCmd.AddCommand(dtrCmd)

}

var dtrCmd = &cobra.Command{
	Use:   "dtr",
	Short: "Docker Trusted Registry",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var dtrLogin = &cobra.Command{
	Use:   "login",
	Short: "Login to a Docker Trusted Registry",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		err := dtrClient.Connect()

		if err != nil {
			log.Fatalf("%v", err)
		} else {
			// If succesfull write the token and annouce as succesful
			err = dtrClient.WriteToken()
			if err != nil {
				log.Errorf("%v", err)
			}
			log.Infof("Succesfully logged into [%s]", dtrClient.DTRURL)
		}
	},
}

var dtrInfo = &cobra.Command{
	Use:   "info",
	Short: "Information about Docker Trusted Registry",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var dtrLoginReplicas = &cobra.Command{
	Use:   "replicas",
	Short: "Docker Trusted Registry Replicase",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := dtr.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.ListReplicas()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
	},
}

var dtrWebHooks = &cobra.Command{
	Use:   "webhook",
	Short: "Docker Trusted Registry Webhooks",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var dtrWebHooksList = &cobra.Command{
	Use:   "list",
	Short: "List Docker Trusted Registry Webhooks",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := dtr.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		wh, err := client.ListWebhooks()
		if err != nil {
			// Fatal error if can't return any webhooks
			log.Fatalf("%v", err)
		}

		const padding = 3
		w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
		fmt.Fprintln(w, "ID\tKey\tType\tEndpoint\tActive")

		for i := range wh {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%t\n", wh[i].ID, wh[i].Key, wh[i].EndPoint, wh[i].Type, wh[i].InActive)
		}
		w.Flush()
	},
}
