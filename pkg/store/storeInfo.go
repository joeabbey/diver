package store

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

//GetUserInfo - Retrieves all subscriptions
func (c *Client) GetUserInfo() ([]byte, error) {
	log.Debugf("Retrieving all subscriptions")
	url := fmt.Sprintf("%s/users/%s", c.STOREURL, c.Username)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}
