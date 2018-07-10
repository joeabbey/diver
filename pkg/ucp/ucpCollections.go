package ucp

import (
	"encoding/json"
	"fmt"

	"github.com/thebsdbox/diver/pkg/ucp/types"

	log "github.com/Sirupsen/logrus"
)

func (c *Client) returnAllCollections() ([]byte, error) {

	// Build the URL (TODO set limit)

	url := fmt.Sprintf("%s/collections", c.UCPURL)

	log.Debugf("Built URL [%s]", url)
	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// GetCollections - This will get all accounts
func (c *Client) GetCollections() ([]ucptypes.Collection, error) {
	data, err := c.returnAllCollections()

	collections := []ucptypes.Collection{}

	err = json.Unmarshal(data, &collections)
	if err != nil {
		return nil, err
	}
	return collections, nil
}
