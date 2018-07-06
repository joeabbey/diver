package cmd

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/thebsdbox/diver/pkg/ucp"
)

// swarminfo
var swarmVersionInfo bool
var swarmVersion string

//raftinfo
var raftHB, raftET int

func init() {

	ucpSwarm.Flags().BoolVar(&swarmVersionInfo, "version", false, "Get the current Swarm cluster version Number")

	ucpSwarmRaftSetInfo.Flags().IntVar(&raftHB, "heartbeat", -1, "Set the Raft Heartbeat")
	ucpSwarmRaftSetInfo.Flags().IntVar(&raftET, "electiontick", -1, "Set the Raft Election Tick")
	ucpSwarmRaftSetInfo.Flags().StringVar(&swarmVersion, "version", "", "Set the version of the Raft Cluster to update")

	ucpSwarm.AddCommand(ucpSwarmRaft)
	ucpSwarmRaft.AddCommand(ucpSwarmRaftGetInfo)
	ucpSwarmRaft.AddCommand(ucpSwarmRaftSetInfo)

	UCPRoot.AddCommand(ucpSwarm)
}

var ucpSwarm = &cobra.Command{
	Use:   "swarm",
	Short: "Manage Docker Swarm settings",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		swarm, err := client.GetSwarmInfo()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		if swarmVersionInfo {
			fmt.Printf("version: %d\n", swarm.Version.Index)
			return
		}
		cmd.Help()
	},
}

var ucpSwarmRaft = &cobra.Command{
	Use:   "raft",
	Short: "Manage Docker Swarm Raft settings",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var ucpSwarmRaftGetInfo = &cobra.Command{
	Use:   "get",
	Short: "Get Docker Swarm Raft settings",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		swarm, err := client.GetSwarmInfo()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		fmt.Printf("ElectionTick: %d\n", swarm.Spec.Raft.ElectionTick)
		fmt.Printf("HeartBeatTick: %d\n", swarm.Spec.Raft.HeartbeatTick)
		fmt.Printf("KeepOldSnapshots: %d\n", swarm.Spec.Raft.KeepOldSnapshots)
		fmt.Printf("LogEntriesForSlowFollowers: %d\n", swarm.Spec.Raft.LogEntriesForSlowFollowers)
		fmt.Printf("SnapShotInterval: %d\n", swarm.Spec.Raft.SnapshotInterval)

	},
}
var ucpSwarmRaftSetInfo = &cobra.Command{
	Use:   "set",
	Short: "Set Docker Swarm Raft settings",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.Level(logLevel))

		client, err := ucp.ReadToken()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
		swarm, err := client.GetSwarmInfo()
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}

		if raftET != -1 {
			swarm.Spec.Raft.ElectionTick = raftET
		}

		if raftHB != -1 {
			swarm.Spec.Raft.HeartbeatTick = raftHB
		}
		err = client.SetSwarmCluster(swarmVersion, &swarm.Spec)
		if err != nil {
			// Fatal error if can't read the token
			log.Fatalf("%v", err)
		}
	},
}
