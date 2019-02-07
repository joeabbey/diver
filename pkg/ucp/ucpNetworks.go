package ucp

import (
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	log "github.com/sirupsen/logrus"
)

var networkCache []types.NetworkResource

// networkIDCache will cache networks from an ID lookup, reducing the amount of API calls needed
func networkIDCache(id string) *types.NetworkResource {
	for i := range networkCache {
		if networkCache[i].ID == id {
			return &networkCache[i]
		}
	}
	return nil
}

//GetNetworks -
func (c *Client) GetNetworks() ([]types.NetworkResource, error) {

	url := fmt.Sprintf("%s/networks", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	// We will get an array of networks from the API call
	var networks []types.NetworkResource

	log.Debugf("Parsing all networks")
	err = json.Unmarshal(response, &networks)
	if err != nil {
		return nil, err
	}
	log.Debugf("Found %d networks", len(networks))

	return networks, nil
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

// NetworkDisconnectContainer - This will remove a container from the network
func (c *Client) NetworkDisconnectContainer(containerID, networkID string, force bool) error {

	var spec struct {
		Container string `json:"Container"`
		Force     bool   `json:"Force"`
	}

	spec.Container = containerID
	spec.Force = force

	b, err := json.Marshal(spec)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/networks/%s/disconnect", c.UCPURL, networkID)

	response, err := c.postRequest(url, b)
	if err != nil {
		ParseUCPError(response)
		return err
	}
	return nil
}

//NetworkConnectContainer - Connect a container to a network
func (c *Client) NetworkConnectContainer(containerID, networkID, ipv4, ipv6 string) error {

	var spec struct {
		Container      string `json:"Container"`
		EndpointConfig struct {
			IPAMConfig struct {
				IPv4Address string `json:"IPv4Address,omitempty"`
				IPv6Address string `json:"IPv6Address,omitempty"`
			} `json:"IPAMConfig"`
		} `json:"EndpointConfig"`
	}

	spec.Container = containerID
	spec.EndpointConfig.IPAMConfig.IPv4Address = ipv4
	spec.EndpointConfig.IPAMConfig.IPv6Address = ipv6

	b, err := json.Marshal(spec)

	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/networks/%s/connect", c.UCPURL, networkID)

	response, err := c.postRequest(url, b)
	if err != nil {
		ParseUCPError(response)
		return err
	}
	return nil
}
