package store

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

//UserAccount - The struct that is returned by interacting with the Docker Store
type UserAccount struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Location string `json:"location"`
	Company  string `json:"company"`
	URL      string `json:"profile_url"`
	Type     string `json:"type"`
}

func (c *Client) userInfo(user string) (*UserAccount, error) {
	log.Debugf("Retrieving all subscriptions")
	var url string
	if user == "" {
		url = fmt.Sprintf("%s/users/%s", c.STOREURL, c.Username)
	} else {
		url = fmt.Sprintf("%s/users/%s", c.STOREURL, user)

	}

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}
	var userInfo UserAccount

	err = json.Unmarshal(response, &userInfo)
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

//GetUserInfo - Retrieves all subscriptions
func (c *Client) GetUserInfo(user string) error {
	userInfo, err := c.userInfo(user)

	if err != nil {
		return err
	}

	fmt.Printf("ID\t%s\n", userInfo.ID)
	fmt.Printf("Username\t%s\n", userInfo.Username)
	fmt.Printf("Full Name\t%s\n", userInfo.FullName)
	fmt.Printf("Location\t%s\n", userInfo.Location)
	fmt.Printf("Company\t%s\n", userInfo.Company)
	fmt.Printf("URL\t%s\n", userInfo.URL)
	fmt.Printf("Type\t%s\n", userInfo.Type)

	return nil
}
