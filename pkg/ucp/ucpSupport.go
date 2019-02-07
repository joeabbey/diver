package ucp

import (
	"bytes"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// GetSupportDump - will download a UCP support dump
func (c *Client) GetSupportDump() error {

	// TODO add a spinner here
	log.Infoln("Downloading a UCP support dump (may take some time)")
	// Create the file
	out, err := os.Create("support-dump.zip")
	if err != nil {
		return err
	}
	defer out.Close()

	url := fmt.Sprintf("%s/api/support", c.UCPURL)

	response, err := c.postRequest(url, nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, bytes.NewReader(response))
	if err != nil {
		return err
	}

	return nil
}
