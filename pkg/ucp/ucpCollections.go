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

// DeleteCollection - This will get all accounts
func (c *Client) DeleteCollection(collectionID string) error {
	// Build the URL (TODO set limit)

	url := fmt.Sprintf("%s/collections/%s", c.UCPURL, collectionID)

	log.Debugf("Built URL [%s]", url)
	_, err := c.delRequest(url, nil)
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

// SetCollection - This will get all accounts
func (c *Client) SetCollection(collectionID string, constraint *ucptypes.CollectionLabelConstraints) error {

	// Find the existing Collection
	collection, err := c.GetCollection(collectionID)
	if err != nil {
		return err
	}

	// Create a temporary struct with the correct JSON schema
	var patchedCollection struct {
		LabelConstraints []ucptypes.CollectionLabelConstraints `json:"label_constraints"`
	}

	patchedCollection.LabelConstraints = collection.LabelConstraints
	// add the new contraint to the existing ones
	patchedCollection.LabelConstraints = append(patchedCollection.LabelConstraints, *constraint)

	// Marshall the JSON to bytes
	b, err := json.Marshal(patchedCollection)
	if err != nil {
		return err
	}

	log.Debugf("%s", b)

	url := fmt.Sprintf("%s/collections/%s", c.UCPURL, collectionID)
	log.Debugf("Built URL [%s]", url)

	_, err = c.patchRequest(url, b)
	if err != nil {
		return err
	}
	return nil
}

// SetDefaultCollection - This will get all accounts
func (c *Client) SetDefaultCollection(collectionID, user string) error {

	// Create a temporary struct with the correct JSON schema
	var collection struct {
		ID string `json:"id"`
	}
	collection.ID = collectionID

	b, err := json.Marshal(collection)
	url := fmt.Sprintf("%s/defaultCollection/%s", c.UCPURL, user)

	response, err := c.putRequest(url, b)
	if err != nil {
		ParseUCPError(response)
		return err
	}
	return nil
}
