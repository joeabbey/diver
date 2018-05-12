package ucp

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
)

//GetNetworks -
func (c *Client) GetNetworks() error {

	url := fmt.Sprintf("%s/networks", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return err
	}

	// We will get an array of networks from the API call
	var networks []types.NetworkResource

	log.Debugf("Parsing all networks")
	err = json.Unmarshal(response, &networks)
	if err != nil {
		return err
	}
	log.Debugf("Found %d networks", len(networks))

	// Loop through all networks in the array
	for i := range networks {
		name := networks[i].Name
		id := networks[i].ID
		fmt.Printf("%s \t %s\n", id, name)
	}
	return nil
}
