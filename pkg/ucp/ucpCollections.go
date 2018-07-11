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

// GetCollections - This will get all collections
func (c *Client) GetCollections() ([]ucptypes.Collection, error) {
	data, err := c.returnAllCollections()

	collections := []ucptypes.Collection{}

	err = json.Unmarshal(data, &collections)
	if err != nil {
		return nil, err
	}
	return collections, nil
}

// CreateCollection - This will get all accounts
func (c *Client) CreateCollection(name, parentID string) error {
	data := map[string]string{
		"name":      name,
		"parent_id": parentID,
	}

	b, err := json.Marshal(data)

	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/collections", c.UCPURL)

	log.Debugf("Built URL [%s]", url)
	_, err = c.postRequest(url, b)
	if err != nil {
		return err
	}
	return nil
}

// GetCollection - This will get all accounts
func (c *Client) GetCollection(collectionID string) (*ucptypes.Collection, error) {
	// Build the URL (TODO set limit)

	url := fmt.Sprintf("%s/collections/%s", c.UCPURL, collectionID)

	log.Debugf("Built URL [%s]", url)
	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	foundCollection := ucptypes.Collection{}

	err = json.Unmarshal(response, &foundCollection)
	if err != nil {
		return nil, err
	}
	return &foundCollection, nil
}
