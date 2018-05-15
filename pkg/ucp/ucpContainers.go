package ucp

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"

	log "github.com/Sirupsen/logrus"
)

var containerCache []types.ContainerJSONBase

// ListContainerJSON -
func (c *Client) ListContainerJSON() error {
	// Add the /auth/log to the URL

	response, err := c.getAllContainers()
	if err != nil {
		return err
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, response, "", "\t")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(prettyJSON.Bytes()))
	return nil
}

// GetContainerCount - Returns the number of containers running
func (c *Client) GetContainerCount() error {
	response, err := c.getAllContainers()
	if err != nil {
		return err
	}
	var containers []types.Container

	err = json.Unmarshal(response, &containers)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d containers\n", len(containers))
	return nil
}

// GetContainerFromID - this will find a container and return it's struct
func (c *Client) GetContainerFromID(id string) (*types.ContainerJSONBase, error) {

	// Added newline to make debugging clearer (makes a mess of normal output)
	log.Debugf("\nLooking up Container from cache")

	cachedContainer := containerIDCache(id)
	if cachedContainer != nil {
		return cachedContainer, nil
	}

	url := fmt.Sprintf("%s/containers/%s/json", c.UCPURL, id)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	var container types.ContainerJSONBase
	err = json.Unmarshal(response, &container)
	if err != nil {
		return nil, err
	}
	log.Debugf("\nAdding new container to cache for faster lookups")
	containerCache = append(containerCache, container)
	return &container, nil
}

// GetContainerNames - lists the names of all containers
func (c *Client) GetContainerNames() error {
	response, err := c.getAllContainers()
	if err != nil {
		return err
	}
	var containers []types.Container

	err = json.Unmarshal(response, &containers)
	if err != nil {
		return err
	}

	for i := range containers {
		fmt.Println(containers[i].Names)
	}
	return nil
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

func (c *Client) getAllContainers() ([]byte, error) {
	log.Debugf("Retrieving all containers")
	url := fmt.Sprintf("%s/containers/json", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ContainerTop -
func (c *Client) ContainerTop() error {
	response, err := c.getAllContainers()
	if err != nil {
		return err
	}
	var containers []types.Container

	log.Debugf("Parsing all containers")
	err = json.Unmarshal(response, &containers)
	if err != nil {
		return err
	}

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

// networkIDCache will cache networks from an ID lookup, reducing the amount of API calls needed
func containerIDCache(id string) *types.ContainerJSONBase {
	for i := range containerCache {
		if containerCache[i].ID == id {
			return &containerCache[i]
		}
	}
	return nil
}
