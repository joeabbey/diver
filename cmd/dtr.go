package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/dtr"

	"github.com/thebsdbox/diver/pkg/dtr/types"
)

var dtrClient dtr.Client
var webhook dtrtypes.DTRWebHook
var repository dtrtypes.DTRRepository

func init() {
	dtrLogin.Flags().StringVar(&dtrClient.Username, "username", os.Getenv("DTR_USERNAME"), "Username that has permissions to authenticate to Docker EE")
	dtrLogin.Flags().StringVar(&dtrClient.Password, "password", os.Getenv("DTR_PASSWORD"), "Password allowing a user to authenticate to Docker EE")
	dtrLogin.Flags().StringVar(&dtrClient.DTRURL, "url", os.Getenv("DTR_PASSWORD"), "The URL of a Docker Trusted Registry")

	ignoreCert := strings.ToLower(os.Getenv("STORE_INSECURE")) == "true"

	dtrLogin.Flags().BoolVar(&dtrClient.IgnoreCert, "ignorecert", ignoreCert, "Ignore x509 certificate")

	dtrWebHooksCreate.Flags().StringVar(&webhook.EndPoint, "endpoint", "", "The Endpoint that will be used as part of the webhook")
	dtrWebHooksCreate.Flags().StringVar(&webhook.Key, "repo", "", "The Repository that the webhook belongs too")
	dtrWebHooksCreate.Flags().StringVar(&webhook.Type, "type", "", "The type of webhook")

	dtrWebHooksDelete.Flags().StringVar(&id, "id", "", "ID of the webhook to delete")

	dtrRepoList.Flags().StringVar(&org, "namespace", "", "The Namespace/Organisation that holds the repositories")

	dtrRepoCreate.Flags().StringVar(&repository.Namespace, "namespace", "", "The Namespace/Organisation that will hold the repositories")
	dtrRepoCreate.Flags().StringVar(&repository.Name, "name", "", "The Name of the new repository")
	dtrRepoCreate.Flags().StringVar(&repository.ShortDescription, "description", "", "A Description about the repository")
	dtrRepoCreate.Flags().StringVar(&repository.Visibility, "visibility", "public", "If the repository should be \"public\" or \"private\"")
	dtrRepoCreate.Flags().BoolVar(&repository.ImmutableTags, "immutable", false, "Repository tags are immutable")
	dtrRepoCreate.Flags().BoolVar(&repository.ScanOnPush, "scan", false, "Vulnerability scans enabled on push")
	dtrRepoCreate.Flags().BoolVar(&repository.EnableManifestLists, "manifest", true, "Enable Repository manifest lists")

	dtrRepoDelete.Flags().StringVar(&repository.Namespace, "namespace", "", "The Namespace/Organisation that holds the repository")
	dtrRepoDelete.Flags().StringVar(&repository.Name, "name", "", "The Name of the repository to delete")

	dtrCmd.AddCommand(dtrLogin)
	dtrCmd.AddCommand(dtrInfo)

	dtrInfo.AddCommand(dtrLoginReplicas)

	dtrWebHooks.AddCommand(dtrWebHooksList)
	dtrWebHooks.AddCommand(dtrWebHooksCreate)
	dtrWebHooks.AddCommand(dtrWebHooksDelete)

	dtrRepos.AddCommand(dtrRepoList)
	dtrRepos.AddCommand(dtrRepoCreate)
	dtrRepos.AddCommand(dtrRepoDelete)

	dtrCmd.AddCommand(dtrWebHooks)
	dtrCmd.AddCommand(dtrRepos)

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
	Short: "Docker Trusted Registry Replicas",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := dtr.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		dc, err := client.DTRClusterStatus()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Replica\tNode")

		for replica, settings := range dc.ReplicaSettings {
			fmt.Fprintf(w, "%s\t%s\n", replica, settings.Node)
		}
		w.Flush()

	},
}

var dtrRepos = &cobra.Command{
	Use:   "repos",
	Short: "Docker Trusted Registry Repositories",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var dtrRepoList = &cobra.Command{
	Use:   "list",
	Short: "List Docker Trusted Registry Repositories",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		client, err := dtr.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		var r []dtrtypes.DTRRepository

		if org == "" {
			log.Debugf("No namespace, returning all repositories")
			r, err = client.ListAllRepositories()
			if err != nil {
				// Fatal error if can't return any webhooks
				log.Fatalf("%v", err)
			}
		} else {
			r, err = client.ListReposForNamespace(org)
			if err != nil {
				// Fatal error if can't return any webhooks
				log.Fatalf("%v", err)
			}
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "Repo\tID\tName\tNamespace\tDescription")

		for i := range r {
			repoName := r[i].Namespace + "/" + r[i].Name
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", repoName, r[i].ID, r[i].Name, r[i].Namespace, r[i].ShortDescription)
		}
		w.Flush()
	},
}

var dtrRepoCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a Docker Trusted Registry Repository",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if repository.Name == "" {
			cmd.Help()
			log.Fatalf("No repository name specified")
		}

		if repository.Namespace == "" {
			cmd.Help()
			log.Fatalf("No repository namespace specified")
		}

		client, err := dtr.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.CreateRepository(repository)
		if err != nil {
			// Fatal error if can't return any webhooks
			log.Fatalf("%v", err)
		}
		log.Infof("New Repository [%s] created for namespace [%s]", repository.Name, repository.Namespace)
	},
}

var dtrRepoDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Docker Trusted Registry Repository",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if repository.Name == "" {
			cmd.Help()
			log.Fatalf("No repository name specified")
		}

		if repository.Namespace == "" {
			cmd.Help()
			log.Fatalf("No repository namespace specified")
		}

		client, err := dtr.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.DeleteRepository(repository)
		if err != nil {
			// Fatal error if can't return any webhooks
			log.Fatalf("%v", err)
		}
		log.Infof("Repository [%s] deleted from namespace [%s]", repository.Name, repository.Namespace)
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

		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabPadding, ' ', 0)
		fmt.Fprintln(w, "ID\tKey\tEndpoint\tType\tInActive")

		for i := range wh {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%t\n", wh[i].ID, wh[i].Key, wh[i].EndPoint, wh[i].Type, wh[i].InActive)
		}
		w.Flush()
	},
}

var dtrWebHooksCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a Docker Trusted Registry Webhooks",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if webhook.Type == "" {
			cmd.Help()
			log.Fatalf("No Webhook type specified")
		}

		if webhook.EndPoint == "" {
			cmd.Help()
			log.Fatalf("No Webhook endpoint specified")
		}

		if webhook.Key == "" {
			cmd.Help()
			log.Fatalf("No repository for the Webhook specified")
		}

		client, err := dtr.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.CreateWebhook(webhook)
		if err != nil {
			// Fatal error if can't return any webhooks
			log.Fatalf("%v", err)
		}
		log.Infof("New Webook type:[%s] created for endpoint [%s]", webhook.Type, webhook.EndPoint)
	},
}

var dtrWebHooksDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Docker Trusted Registry Webhooks",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))
		if id == "" {
			log.Fatalf("No DTR webhook specified")
		}

		client, err := dtr.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.DeleteWebhook(id)
		if err != nil {
			// Fatal error if can't return any webhooks
			log.Fatalf("%v", err)
		}
		log.Infof("Webhook [%s] succesfully deleted", id)
	},
}
