package ucp

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/thebsdbox/diver/pkg/ucp/types"
)

//ListConfigs - This will return a list of services
func (c *Client) ListConfigs() ([]ucptypes.ServiceConfig, error) {
	url := fmt.Sprintf("%s/configs", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	// We will get an array of nodes from the API call
	var configs []ucptypes.ServiceConfig

	log.Debugf("Parsing all Configurations")
	err = json.Unmarshal(response, &configs)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

//GetConfig - This will return a list of services
func (c *Client) GetConfig(id string) (*ucptypes.ServiceConfig, error) {
	url := fmt.Sprintf("%s/configs/%s", c.UCPURL, id)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	// We will get an array of nodes from the API call
	var config ucptypes.ServiceConfig

	log.Debugf("Parsing all nodes")
	err = json.Unmarshal(response, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
