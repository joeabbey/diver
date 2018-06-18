package ucp

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

//GetOrg -
func (c *Client) GetOrg(orgName string) error {
	log.Debugf("Searching for Org [%s]", orgName)
	return nil
}

//GetRoles - This will return a list of services
func (c *Client) GetRoles() error {

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
