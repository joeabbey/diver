package ucp

import (
	log "github.com/Sirupsen/logrus"
)

//GetOrg -
func (c *Client) GetOrg(orgName string) error {
	log.Debugf("Searching for Org [%s]", orgName)
	return nil
}
