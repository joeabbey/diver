package ucp

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
)

var networkCache []types.NetworkResource

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

// GetNetworkFromID - this will find a container and return it's struct
func (c *Client) GetNetworkFromID(id string) (*types.NetworkResource, error) {
	// Added newline to make debugging clearer (makes a mess of normal output)
	log.Debugf("\nLooking up Network from cache")

	cachedNetwork := networkIDCache(id)
	if cachedNetwork != nil {
		return cachedNetwork, nil
	}

	log.Debugf("Network not found in cache, using API lookup")
	url := fmt.Sprintf("%s/networks/%s", c.UCPURL, id)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	var network types.NetworkResource
	err = json.Unmarshal(response, &network)
	if err != nil {
		return nil, err
	}
	// Add network to cache to speed further lookups
	networkCache = append(networkCache, network)
	return &network, nil
}

// networkIDCache will cache networks from an ID lookup, reducing the amount of API calls needed
func networkIDCache(id string) *types.NetworkResource {
	for i := range networkCache {
		if networkCache[i].ID == id {
			return &networkCache[i]
		}
	}
	return nil
}
