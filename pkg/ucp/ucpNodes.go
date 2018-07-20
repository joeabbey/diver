package ucp

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/swarm"
)

//ListAllNodes - Retrieves the complete list of all nodes connected to a UCP cluster
func (c *Client) ListAllNodes() ([]swarm.Node, error) {

	url := fmt.Sprintf("%s/nodes", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	// We will get an array of nodes from the API call
	var nodes []swarm.Node

	log.Debugf("Parsing all nodes")
	err = json.Unmarshal(response, &nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
