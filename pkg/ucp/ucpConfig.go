package ucp

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/joeabbey/diver/pkg/ucp/types"
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

//CreateConfig - This will return a list of services
func (c *Client) CreateConfig(name, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/configs/create", c.UCPURL)

	var config struct {
		Data string `json:"Data"`
		Name string `json:"Name"`
	}
	config.Name = name
	config.Data = base64.StdEncoding.EncodeToString(data)

	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	response, err := c.postRequest(url, b)
	if err != nil {
		ParseUCPError(response)
		return err
	}
	return nil
}
