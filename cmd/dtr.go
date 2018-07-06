package cmd

import (
	"os"
	"strings"

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
