package ucp

import (
	"encoding/json"
	"fmt"

	"github.com/thebsdbox/diver/pkg/ucp/types"

	"github.com/docker/docker/api/types"

	log "github.com/Sirupsen/logrus"
)

var containerCache []types.ContainerJSON

// networkIDCache will cache networks from an ID lookup, reducing the amount of API calls needed
func containerIDCache(id string) *types.ContainerJSON {
	for i := range containerCache {
		if containerCache[i].ID == id {
			return &containerCache[i]
		}
	}
	return nil
}

// GetContainerCount - Returns the number of containers running
func (c *Client) GetContainerCount() error {
	containers, err := c.GetAllContainers()
	if err != nil {
		return err
	}

	fmt.Printf("Found %d containers\n", len(containers))
	return nil
}

// GetContainerFromID - this will find a container and return it's struct
func (c *Client) GetContainerFromID(id string) (*types.ContainerJSON, error) {

	// Added newline to make debugging clearer (makes a mess of normal output)
	log.Debugln("Looking up Container from cache")

	cachedContainer := containerIDCache(id)
	if cachedContainer != nil {
		return cachedContainer, nil
	}

	url := fmt.Sprintf("%s/containers/%s/json", c.UCPURL, id)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	var container types.ContainerJSON
	err = json.Unmarshal(response, &container)
	if err != nil {
		return nil, err
	}
	log.Debugln("Adding new container to cache for faster lookups")
	containerCache = append(containerCache, container)
	return &container, nil
}

func (c *Client) getContainerMem(id string) (uint64, uint64, error) {
	url := fmt.Sprintf("%s/containers/%s/stats?stream=false", c.UCPURL, id)

	log.Debugf("Retrieving %s", url)
	response, err := c.getRequest(url, nil)
	if err != nil {
		return 0, 0, err
	}

	var stats types.Stats

	err = json.Unmarshal(response, &stats)
	if err != nil {
		return 0, 0, err
	}

	return stats.MemoryStats.Usage, stats.MemoryStats.Limit, nil

}

//GetAllContainers - Returns all containers in the cluster
func (c *Client) GetAllContainers() ([]types.Container, error) {
	log.Debugf("Retrieving all containers")
	url := fmt.Sprintf("%s/containers/json", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}
	var containers []types.Container

	log.Debugf("Parsing all containers")
	err = json.Unmarshal(response, &containers)
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// ContainerTop -
func (c *Client) ContainerTop() error {
	containers, err := c.GetAllContainers()
	if err != nil {
		return err
	}

	log.Debugf("Parsing all containers")

	for i := range containers {
		name := containers[i].Names
		id := containers[i].ID
		usage, limit, err := c.getContainerMem(id)
		if err != nil {
			return nil
		}
		percentage := ((float64(usage) / float64(limit)) * 100)
		if percentage > 90 {
			fmt.Printf("\033[35m%0.2f%%\033[m  %s %s\n", percentage, id, name)
		} else if percentage > 75 {
			fmt.Printf("\033[35m%0.2f%%\033[m  %s %s\n", percentage, id, name)
		} else {
			fmt.Printf("\033[32m%0.2f%%\033[m  %s %s\n", percentage, id, name)
		}
	}
	return nil
}

//GetContainerProcesses - Returns all containers in the cluster
func (c *Client) GetContainerProcesses(id string) (*ucptypes.ContainerProcesses, error) {
	log.Debugf("Retrieving processers from container")
	url := fmt.Sprintf("%s/containers/%s/top?ps_args=-ef", c.UCPURL, id)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}
	var processList ucptypes.ContainerProcesses

	err = json.Unmarshal(response, &processList)
	if err != nil {
		return nil, err
	}

	log.Debugf("%s", response)
	return &processList, nil
}
