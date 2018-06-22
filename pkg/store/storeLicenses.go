package store

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

//GetLicense - Retrieves all subscriptions
func (c *Client) GetLicense(subscription string) error {
	log.Debugf("Retrieving all subscriptions")
	url := fmt.Sprintf("%s/%s/license-file", c.HUBURL, subscription)
	log.Debugf("Url = %s", url)
	response, err := c.getRequest(url, nil)
	if err != nil {
		return err
	}
	fmt.Printf("%v", string(response))
	return nil
}
