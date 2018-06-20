package ucp

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

//GetOrg - TODO
func (c *Client) GetOrg(orgName string) error {
	log.Debugf("Searching for Org [%s]", orgName)
	return nil
}

//GetRoles - This will return a list of services
func (c *Client) GetRoles() error {

	url := fmt.Sprintf("%s/roles", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return err
	}

	var roles []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ServiceRole bool   `json:"servicerole"`
		//Operations  interface{} `json:"operations,omitempty"`
		Operations json.RawMessage // Captures the raw output of the remaining json object
	}

	log.Debugf("Parsing all roles")
	err = json.Unmarshal(response, &roles)
	if err != nil {
		return err
	}

	fmt.Printf("ID\t\tService Account\tName\n")

	for i := range roles {
		fmt.Printf("%s\t%t\t%s\n", roles[i].ID, roles[i].ServiceRole, roles[i].Name)
	}
	return nil
}

//GetRole - This will return a list of services
func (c *Client) GetRole(role string) (string, error) {

	url := fmt.Sprintf("%s/roles", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return "", err
	}

	var roles []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ServiceRole bool   `json:"servicerole"`
		//Operations  interface{} `json:"operations,omitempty"`
		Operations json.RawMessage // Captures the raw output of the remaining json object
	}

	log.Debugf("Parsing all roles")
	err = json.Unmarshal(response, &roles)
	if err != nil {
		return "", err
	}

	for i := range roles {
		if role == roles[i].Name {
			return string(roles[i].Operations), nil
		}

	}
	return "", fmt.Errorf("Unable to find role [%s]", role)
}

//SetRole - This set the role of a user in an organisation
func (c *Client) SetRole(user, org, role string) error {

	url := fmt.Sprintf("%s/roles", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil
	}

	var roles []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ServiceRole bool   `json:"servicerole"`
	}

	log.Debugf("Parsing all roles")
	err = json.Unmarshal(response, &roles)
	if err != nil {
		return err
	}

	fmt.Printf("Name\t\tID\t\tService Account\n")

	for i := range roles {
		fmt.Printf("%s\t%s\t%t\n", roles[i].Name, roles[i].ID, roles[i].ServiceRole)
	}

	return nil
}
