package ucp

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/swarm"
)

// ServiceQuery - This collection of bools builds the output from a service query
type ServiceQuery struct {
	// The name of the service to query
	ServiceName string

	// The query output
	ID       bool
	Networks bool
	State    bool

	// Resolve UUIDs to Name
	Resolve bool
}

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

// QueryServiceContainers - This takes a query struct and builds output
func (c *Client) QueryServiceContainers(q *ServiceQuery) error {

	// Build JSON object => e.g. {"service": ["task_test"]}

	// TODO - this is a hack as html.escapestring() wont escape "{}:"
	beginEncode := "%7B%22service%22%3A%5B%22"
	endEncode := "%22%5D%7D"
	encodeString := beginEncode + q.ServiceName + endEncode

	url := fmt.Sprintf("%s/tasks?filters=%s", c.UCPURL, encodeString)

	log.Debugf("Built url %s", url)

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
		image := tasks[i].Status.ContainerStatus.ContainerID

		var networkString string

		for n := range tasks[i].NetworksAttachments {
			for a := range tasks[i].NetworksAttachments[n].Addresses {
				networkString = networkString + "\t" + tasks[i].NetworksAttachments[n].Addresses[a]
			}

		}

		fmt.Printf("%s \t %s\n", image, networkString)
	}

	return nil
}
