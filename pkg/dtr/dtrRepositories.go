package dtr

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/joeabbey/diver/pkg/dtr/types"
)

//ListAllRepositories -
func (c *Client) ListAllRepositories() ([]dtrtypes.DTRRepository, error) {

	//TODO - pageSize needs more help
	url := fmt.Sprintf("%s/api/v0/repositories?pageSize=1000&count=true&refresh_token=%s", c.DTRURL, c.Token)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	// Create temporary struct that contains a simple array
	var repos struct {
		Repositories []dtrtypes.DTRRepository `json:"repositories"`
	}

	err = json.Unmarshal(response, &repos)
	if err != nil {
		return nil, err
	}
	// Return only the array of repos
	return repos.Repositories, nil
}

//ListReposForNamespace -
func (c *Client) ListReposForNamespace(ns string) ([]dtrtypes.DTRRepository, error) {

	//TODO - pageSize needs more help
	url := fmt.Sprintf("%s/api/v0/repositories/%s?pageSize=1000&count=true&refresh_token=%s", c.DTRURL, ns, c.Token)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	// Create temporary struct that contains a simple array
	var repos struct {
		Repositories []dtrtypes.DTRRepository `json:"repositories"`
	}

	err = json.Unmarshal(response, &repos)
	if err != nil {
		return nil, err
	}
	// Return only the array of repos
	return repos.Repositories, nil
}

//CreateRepository -
func (c *Client) CreateRepository(repo dtrtypes.DTRRepository) error {

	url := fmt.Sprintf("%s/api/v0/repositories/%s?refresh_token=%s", c.DTRURL, repo.Namespace, c.Token)

	b, err := json.Marshal(repo)
	if err != nil {
		return err
	}
	log.Debugf("%s", b)
	response, err := c.postRequest(url, b)
	if err != nil {
		log.Debugf("%s", response)
		return err
	}
	return nil
}

//DeleteRepository -
func (c *Client) DeleteRepository(repo dtrtypes.DTRRepository) error {

	url := fmt.Sprintf("%s/api/v0/repositories/%s/%s?refresh_token=%s", c.DTRURL, repo.Namespace, repo.Name, c.Token)

	b, err := json.Marshal(repo)
	if err != nil {
		return err
	}
	log.Debugf("%s", b)
	response, err := c.delRequest(url, b)
	if err != nil {
		log.Debugf("%s", response)
		return err
	}
	return nil
}
