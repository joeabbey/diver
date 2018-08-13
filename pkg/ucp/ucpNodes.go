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

//GetNode - Retrieves the complete list of all nodes connected to a UCP cluster
func (c *Client) GetNode(id string) (swarm.Node, error) {

	url := fmt.Sprintf("%s/nodes/%s", c.UCPURL, id)
	// We will get the struct of a node from the API call
	var node swarm.Node

	response, err := c.getRequest(url, nil)
	if err != nil {
		return node, err
	}

	log.Debugf("Parsing Node details")
	err = json.Unmarshal(response, &node)
	if err != nil {
		return node, err
	}

	return node, nil
}

//SetNodeLabel - Retrieves the complete list of all nodes connected to a UCP cluster
func (c *Client) SetNodeLabel(id, k, v string) error {

	log.Debugln("Retrieving information about existing configuration")
	node, err := c.GetNode(id)
	if err != nil {
		return err
	}

	// Modify the node spec labels
	node.Spec.Labels[k] = v

	b, err := json.Marshal(node.Spec)
	if err != nil {
		return err
	}
	log.Debugf("%s", b)
	url := fmt.Sprintf("%s/nodes/%s/update?version=%d", c.UCPURL, id, node.Version.Index)

	response, err := c.postRequest(url, b)
	if err != nil {
		ParseUCPError(response)
		return err
	}
	return nil
}
