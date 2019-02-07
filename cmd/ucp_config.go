package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"text/tabwriter"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/joeabbey/diver/pkg/ucp"
)

var path string

func init() {

	ucpConfigGet.Flags().StringVar(&id, "id", "", "The ID of Configuration to get")
	ucpConfigGet.Flags().StringVar(&path, "path", "", "The path to write the configuration to")

	ucpConfigCreate.Flags().StringVar(&id, "name", "", "The name of Configuration to create")
	ucpConfigCreate.Flags().StringVar(&path, "path", "", "The path of the configuration to upload")

	ucpConfig.AddCommand(ucpConfigGet)
	ucpConfig.AddCommand(ucpConfigCreate)
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

var ucpConfigGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve a Docker EE Service Configuration",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}

		if id == "" {
			cmd.Help()
			log.Fatalf("No Configuration ID Specified")
		}

		if path == "" {
			cmd.Help()
			log.Fatalf("No Path to save config data Specified")
		}

		cfg, err := client.GetConfig(id)
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		data, err := base64.StdEncoding.DecodeString(cfg.Spec.Data)
		if err != nil {
			log.Fatalf("%v", err)
		}

		err = ioutil.WriteFile(path, data, 0644)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Infof("Succesfully written configuration [%s] to [%s]", cfg.Spec.Name, path)
	},
}

var ucpConfigCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a Docker EE Service Configuration",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			log.Fatalf("%v", err)
		}
		if id == "" {
			cmd.Help()
			log.Fatalf("No Configuration Name Specified")
		}

		if path == "" {
			cmd.Help()
			log.Fatalf("No Path to config data Specified")
		}
		err = client.CreateConfig(id, path)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Infof("Succesfully created config [%s] from file [%s]", id, path)
	},
}
