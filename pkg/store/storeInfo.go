package store

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

//GetUserInfo - Retrieves all subscriptions
func (c *Client) GetUserInfo(user string) error {
	log.Debugf("Retrieving all subscriptions")
	var url string
	if user == "" {
		url = fmt.Sprintf("%s/users/%s", c.STOREURL, c.Username)
	} else {
		url = fmt.Sprintf("%s/users/%s", c.STOREURL, user)

	}

	response, err := c.getRequest(url, nil)
	if err != nil {
		return err
	}

	var UserInfo struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		FullName string `json:"full_name"`
		Location string `json:"location"`
		Company  string `json:"company"`
		URL      string `json:"profile_url"`
		Type     string `json:"type"`
	}

	err = json.Unmarshal(response, &UserInfo)
	if err != nil {
		return err
	}

	fmt.Printf("ID\t%s\n", UserInfo.ID)
	fmt.Printf("Username\t%s\n", UserInfo.Username)
	fmt.Printf("Full Name\t%s\n", UserInfo.FullName)
	fmt.Printf("Location\t%s\n", UserInfo.Location)
	fmt.Printf("Company\t%s\n", UserInfo.Company)
	fmt.Printf("URL\t%s\n", UserInfo.URL)
	fmt.Printf("Type\t%s\n", UserInfo.Type)

	return nil
}
