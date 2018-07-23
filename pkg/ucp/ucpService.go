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

	log.Debugf("Found %d tasks for service %s", len(tasks), q.ServiceName)

	// Loop through all networks in the array
	for i := range tasks {

		// Print task ID
		if q.ID {
			task := tasks[i].Status.ContainerStatus.ContainerID
			if q.Resolve {
				resolvedTask, err := c.GetContainerFromID(task)
				if err != nil {
					// Usually we return from all errors, however we may have lost container IDs
					ParseUCPError([]byte(err.Error()))
					// continue goes to the next loop
					continue
				} else {
					fmt.Printf("%s\t", resolvedTask.Name)
				}
			} else {
				fmt.Printf("%s\t", task)
			}
		}

		// Above query will have cached the results if the container was found
		if q.Node {
			containerNode, err := c.GetContainerFromID(tasks[i].Status.ContainerStatus.ContainerID)
			if err != nil {
				// Usually we return from all errors, however we may have lost container IDs
				ParseUCPError([]byte(err.Error()))
				// continue goes to the next loop
				continue
			} else {
				fmt.Printf("%s\t", containerNode.Node.Name)
			}
		}

		// Print all networks attached to task (Only if attachements exist)
		if q.Networks && len(tasks[i].NetworksAttachments) != 0 {
			var networkString string
			for n := range tasks[i].NetworksAttachments {
				for a := range tasks[i].NetworksAttachments[n].Addresses {

					address := tasks[i].NetworksAttachments[n].Addresses[a]
					networkID := tasks[i].NetworksAttachments[n].Network.ID

					// build output from query
					if q.Resolve {
						resolvedNetwork, err := c.GetNetworkFromID(networkID)
						if err != nil {
							return err
						}
						// Build from resolved name
						networkString = networkString + address + "\t" + resolvedNetwork.Name + "\t"
					} else {
						// Build from UUID name
						networkString = networkString + address + "\t" + networkID + "\t"
					}
				}
			}
			fmt.Printf("%s", networkString)
		} else {
			fmt.Printf("Unattached\t")
		}

		if q.State {
			fmt.Printf("%s\t", tasks[i].Status.State)
		}

		// Create a newline for the next task
		fmt.Printf("\n")
	}

	return nil
}
