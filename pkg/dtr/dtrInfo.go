package dtr

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/thebsdbox/diver/pkg/dtr/types"
)

func (c *Client) dtrClusterStatus() (*dtrtypes.DTRCluster, error) {

	url := fmt.Sprintf("%s/api/v0/meta/cluster_status?refresh_token=%s", c.DTRURL, c.Token)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}
	//log.Debugf("%v", string(response))
	var info dtrtypes.DTRCluster

	err = json.Unmarshal(response, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

//DTRClusterStatus -
func (c *Client) DTRClusterStatus() (*dtrtypes.DTRSettings, error) {

	url := fmt.Sprintf("%s/api/v0/meta/settings?refresh_token=%s", c.DTRURL, c.Token)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}
	//log.Debugf("%v", string(response))
	var info dtrtypes.DTRSettings

	err = json.Unmarshal(response, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

//ListReplicas -
func (c *Client) ListReplicas() error {
	cluster, err := c.dtrClusterStatus()
	if err != nil {
		return err
	}

	replicas := cluster.ReplicaHealth

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "Replica\tStatus")

	for replica, status := range replicas {
		fmt.Fprintf(w, "%s\t%s\n", replica, status)
	}
	w.Flush()

	return nil
}
