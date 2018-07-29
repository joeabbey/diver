package store

import (
	"bytes"
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
)

// GetLicense - Retrieves all subscriptions
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

// GetCVEDatabase - Retrieves CVE Database
func (c *Client) GetCVEDatabase(subscription, dtrSchema string) error {
	log.Debugf("Retrieving CVE database")

	out, err := os.Create("cve-file.tar")
	if err != nil {
		return err
	}
	defer out.Close()

	url := fmt.Sprintf("%s/%s/cve-file?schema_version=%s", c.HUBURL, subscription, dtrSchema)
	log.Debugf("Url = %s", url)

	if err != nil {
		panic(err)
	}

	response, err := c.getRequest(url, nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, bytes.NewReader(response))
	if err != nil {
		return err
	}

	return nil
}
