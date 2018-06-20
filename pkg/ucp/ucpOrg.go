package ucp

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

type roles struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	ServiceRole bool            `json:"servicerole"`
	Operations  json.RawMessage // Captures the raw output of the remaining json object
}

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

	var r []roles

	log.Debugf("Parsing all roles")
	err = json.Unmarshal(response, &r)
	if err != nil {
		return err
	}

	fmt.Printf("ID\t\tService Account\tName\n")

	for i := range r {
		fmt.Printf("%s\t%t\t%s\n", r[i].ID, r[i].ServiceRole, r[i].Name)
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

	var r []roles

	log.Debugf("Parsing all roles")
	err = json.Unmarshal(response, &r)
	if err != nil {
		return "", err
	}

	for i := range r {
		if role == r[i].Name {
			return string(r[i].Operations), nil
		}

	}
	return "", fmt.Errorf("Unable to find role [%s]", role)
}

//CreateRole - This set the role of a user in an organisation
func (c *Client) CreateRole(name, id, ruleset string, serviceAccount bool) error {

	url := fmt.Sprintf("%s/roles", c.UCPURL)

	newrole := roles{
		ID:          id,
		Name:        name,
		ServiceRole: serviceAccount,
		Operations:  json.RawMessage(ruleset),
	}

	b, err := json.Marshal(newrole)

	if err != nil {
		return err
	}

	response, err := c.postRequest(url, b)
	if err != nil {
		return nil
	}

	log.Debugf("%v", string(response))

	return nil
}
