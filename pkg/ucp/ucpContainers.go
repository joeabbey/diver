package ucp

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
)

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
	var container []types.Container

	err = json.Unmarshal(response, &container)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d containers\n", len(container))
	return nil
}

func (c *Client) getAllContainers() ([]byte, error) {
	url := fmt.Sprintf("%s/containers/json", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}
