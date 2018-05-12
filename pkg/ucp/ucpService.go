package ucp

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/swarm"
)

//GetServices - This will return a list of services
func (c *Client) GetServices() error {

	url := fmt.Sprintf("%s/services", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return err
	}

	// We will get an array of services from the API call
	var services []swarm.Service

	log.Debugf("Parsing all services")
	err = json.Unmarshal(response, &services)
	if err != nil {
		return err
	}

	log.Debugf("Found %d services", len(services))

	// Loop through all networks in the array
	for i := range services {
		name := services[i].Spec.Name
		id := services[i].ID
		fmt.Printf("%s \t %s\n", id, name)
	}
	return nil
}

// GetServiceContainers =
func (c *Client) GetServiceContainers(service string) error {

	var filter struct {
		Service []string `json:"service"`
	}

	filter.Service = append(filter.Service, service)

	url := fmt.Sprintf("%s/tasks?filters=%s", c.UCPURL, filter)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return err
	}

	var tasks []swarm.Task

	err = json.Unmarshal(response, &tasks)
	if err != nil {
		return err
	}

	log.Debugf("Found %d tasks", len(tasks))

	// Loop through all networks in the array
	for i := range tasks {
		image := tasks[i].Spec.ContainerSpec.Image
		id := tasks[i].ID
		fmt.Printf("%s \t %s\n", image, id)
	}

	return nil
}
