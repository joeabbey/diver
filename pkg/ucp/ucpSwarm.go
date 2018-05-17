package ucp

import (
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types/swarm"
)

//GetSwarmInfo - returns all information about the swarm cluster
func (c *Client) GetSwarmInfo() (*swarm.ClusterInfo, error) {
	url := fmt.Sprintf("%s/swarm", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	var swm swarm.ClusterInfo
	err = json.Unmarshal(response, &swm)
	if err != nil {
		return nil, err
	}
	return &swm, nil
}

//SetSwarmCluster - takes an updated cluster configuration and applies it to the existing swarm version
func (c *Client) SetSwarmCluster(version string, swm *swarm.Spec) error {
	url := fmt.Sprintf("%s/swarm/update?version=%s", c.UCPURL, version)

	b, err := json.Marshal(swm)

	if err != nil {
		return err
	}

	_, err = c.postRequest(url, b)
	if err != nil {
		return err
	}

	return nil
}
