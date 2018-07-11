package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

var bld ucp.BuildPlan

func init() {
	// Service flags
	ucpBuild.Flags().StringVar(&bld.GHURL, "github", "", "URL to a dockerfile on github")
	ucpBuild.Flags().StringVar(&bld.Tag, "tag", "", "Docker tags to apply to the image")
	//ucpBuild.Flags().StringVar(&bld.BuildHost, "buildhost", "", "Docker engine where Image will be built")

	// Add Service to UCP root commands
	if !DiverRO {
		UCPRoot.AddCommand(ucpBuild)
	}
}

var ucpBuild = &cobra.Command{
	Use:   "build",
	Short: "Build a new Image on Docker EE",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		err = client.BuildImage(&bld)
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
	},
}
