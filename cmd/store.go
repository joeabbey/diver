package cmd

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/thebsdbox/diver/pkg/store"
)

var storeClient store.Client
var id string
var firstActive bool

func init() {
	storeCmd.Flags().StringVar(&storeClient.Username, "username", os.Getenv("STORE_USERNAME"), "Username that has permissions to authenticate to Docker EE")
	storeCmd.Flags().StringVar(&storeClient.Password, "password", os.Getenv("STORE_PASSWORD"), "Password allowing a user to authenticate to Docker EE")
	storeCmd.Flags().StringVar(&storeClient.STOREURL, "storeurl", "https://hub.docker.com/v2", "The Docker Store URL")
	storeCmd.Flags().StringVar(&storeClient.HUBURL, "huburl", "https://store.docker.com/api/billing/v4/subscriptions", "The Docker Hub URL")

	ignoreCert := strings.ToLower(os.Getenv("STORE_INSECURE")) == "true"

	storeCmd.Flags().BoolVar(&storeClient.IgnoreCert, "ignorecert", ignoreCert, "Ignore x509 certificate")

	storeSubscriptionsList.Flags().StringVar(&id, "id", "", "Docker Store ID, by default will take the ID from ~/.storetoken")
	storeSubscriptionsList.Flags().BoolVar(&firstActive, "firstactive", false, "Retrieve first active subscription")

	storeUser.Flags().StringVar(&id, "user", "", "Retrieve information about a specified user")

	storeLicensesGet.Flags().StringVar(&id, "subscription", "", "Set which subscription to retrieve the license")

	storeLicenses.AddCommand(storeLicensesGet)

	storeCmd.AddCommand(storeLicenses)
	storeCmd.AddCommand(storeSubscriptions)
	storeCmd.AddCommand(storeUser)
	storeSubscriptions.AddCommand(storeSubscriptionsList)
	diverCmd.AddCommand(storeCmd)

}

var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Docker Store",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		_, err := store.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Warn("Unable to find existing session, please login")
		} else {
			log.Infof("Existing session found at ~/.storetoken")
			return
		}

		err = storeClient.Connect()

		// Check if connection was succesful
		if err != nil {
			cmd.Help()
			log.Fatalf("%v", err)
		} else {
			// If succesfull write the token and annouce as succesful
			err = storeClient.WriteToken()
			if err != nil {
				log.Errorf("%v", err)
			}
			log.Infof("Succesfully logged into [%s]", storeClient.STOREURL)
		}
	},
}

var storeSubscriptions = &cobra.Command{
	Use:   "subscriptions",
	Short: "Docker Store subscriptions",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var storeSubscriptionsList = &cobra.Command{
	Use:   "ls",
	Short: "List Docker Store subscriptions",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		existingClient, err := store.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			cmd.Help()
			log.Warn("Unable to find existing session, please login")
			return
		}
		if firstActive == true {
			err = existingClient.GetFirstActiveSubscription(id)
			if err != nil {
				log.Fatalf("%v", err)
			}
			return
		}

		err = existingClient.GetAllSubscriptions(id)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}

var storeLicenses = &cobra.Command{
	Use:   "licenses",
	Short: "Docker Store licenses",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var storeLicensesGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve a Subscription license",
	Run: func(cmd *cobra.Command, args []string) {
		existingClient, err := store.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			cmd.Help()
			log.Warn("Unable to find existing session, please login")
			return
		}
		existingClient.GetLicense(id)
	},
}

var storeUser = &cobra.Command{
	Use:   "user",
	Short: "Return Docker Store User Information",
	Run: func(cmd *cobra.Command, args []string) {
		existingClient, err := store.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			cmd.Help()
			log.Warn("Unable to find existing session,")
			return
		}
		err = existingClient.GetUserInfo(id)
		if err != nil {
			log.Fatalf("%v", err)
		}
	},
}
