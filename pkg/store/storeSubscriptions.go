package store

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

//GetAllSubscriptions - Retrieves all subscriptions
func (c *Client) GetAllSubscriptions(id string) ([]byte, error) {
	log.Debugf("Retrieving all subscriptions")
	url := fmt.Sprintf("%s/?docker_id=%s", c.HUBURL, id)
	log.Debugf("Url = %s", url)
	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}
