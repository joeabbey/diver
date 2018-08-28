package ucp

import (
	"encoding/json"
	"fmt"

	diff "github.com/thebsdbox/jd/lib"

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
	Node     bool
	// Resolve UUIDs to Name
	Resolve bool
}

//GetService - This will return a list of services
func (c *Client) GetService(service string) (*swarm.Service, error) {

	url := fmt.Sprintf("%s/services/%s", c.UCPURL, service)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	log.Debugf("Parsing all services")
	var svc swarm.Service
	err = json.Unmarshal(response, &svc)
	if err != nil {
		return nil, err
	}
	return &svc, nil
}

//GetAllServices - This will return a list of services
func (c *Client) GetAllServices() ([]swarm.Service, error) {

	url := fmt.Sprintf("%s/services", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	// We will get an array of services from the API call
	var services []swarm.Service

	log.Debugf("Parsing all services")
	err = json.Unmarshal(response, &services)
	if err != nil {
		return nil, err
	}

	log.Debugf("Found %d services", len(services))
	return services, nil

}

// GetServiceTasks - This returns all tasks associated with a service
func (c *Client) GetServiceTasks(serviceName string) ([]swarm.Task, error) {

	// Build JSON object => e.g. {"service": ["task_test"]}

	// TODO - this is a hack as html.escapestring() wont escape "{}:"
	beginEncode := "%7B%22service%22%3A%5B%22"
	endEncode := "%22%5D%7D"
	encodeString := beginEncode + serviceName + endEncode

	url := fmt.Sprintf("%s/tasks?filters=%s", c.UCPURL, encodeString)

	log.Debugf("Built url %s", url)

	response, err := c.getRequest(url, nil)
	if err != nil {
		//TODO - Must be a nicer method for this
		ParseUCPError(response)
		return nil, err
	}

	var tasks []swarm.Task

	err = json.Unmarshal(response, &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// ReapFailedTasks - This returns all tasks associated with a service
func (c *Client) ReapFailedTasks(serviceName string, rmvol, kill bool) error {

	tasks, err := c.GetServiceTasks(serviceName)
	if err != nil {
		return err
	}

	log.Debugf("Found %d tasks, finding tasks in \"failed\" status", len(tasks))

	if len(tasks) == 0 {
		return fmt.Errorf("No tasks found for service [%s]", serviceName)
	}

	var failedTasks []string

	for i := range tasks {
		if tasks[i].Status.State == "failed" {
			// Look for an existing container that matches the ID
			_, err := c.GetContainerFromID(tasks[i].Status.ContainerStatus.ContainerID)
			// If we find one, we know that it exists and can be added to be reaped
			if err == nil {
				failedTasks = append(failedTasks, tasks[i].Status.ContainerStatus.ContainerID)
			}
		}
	}

	if len(failedTasks) == 0 {
		log.Info("No tasks in a \"failed\" state")
	}

	for i := range failedTasks {
		fmt.Printf("Removing %s\n", failedTasks[i])

		// DELETE CODE GOES HERE

		// url := fmt.Sprintf("%s/services/%s", c.UCPURL, failedTasks[i])
		// _, err = c.delRequest(url, nil)
		// if err != nil {
		// 	return err
		// }

	}

	return nil
}

//GetServiceDifference - This will compare a service spec and previous spec and output the differences
func (c *Client) GetServiceDifference(serviceName string, pretty bool) (string, error) {

	svc, err := c.GetService(serviceName)
	if err != nil {
		return "", err
	}

	if svc.PreviousSpec == nil {
		return "", fmt.Errorf("No previous service spec")
	}

	// Retrieve the two spec JSON objects
	spec, err := json.Marshal(svc.Spec)
	if err != nil {
		return "", err
	}

	prevspec, err := json.Marshal(svc.PreviousSpec)
	if err != nil {
		return "", err
	}

	log.Debugln("Both Specifications have been unmarshalled, performing a diff on them")
	a, _ := diff.ReadJsonString(string(spec))
	b, _ := diff.ReadJsonString(string(prevspec))

	return b.Diff(a).Render(pretty), nil
}
