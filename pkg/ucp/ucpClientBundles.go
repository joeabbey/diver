package ucp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/joeabbey/diver/pkg/ucp/types"
)

// GetClientBundle - will download the UCP Client Bundle
func (c *Client) GetClientBundle() error {

	log.Infoln("Downloading the UCP Client Bundle")
	// Create the file
	out, err := os.Create("ucp-bundle.zip")
	if err != nil {
		return err
	}
	defer out.Close()

	log.Debugln("Retrieving Client Bundle")
	url := fmt.Sprintf("%s/api/clientbundle", c.UCPURL)

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

// ListClientBundles - will list UCP Client Bundle
func (c *Client) ListClientBundles(user string) (*ucptypes.ClientBundles, error) {
	if user == "" {
		currentAccount, err := c.AuthStatus()
		if err != nil {
			log.Warn("Session has expired, please login")
			return nil, err
		}
		user = currentAccount.Name
	}
	log.Debugln("Listing Client Bundles")

	url := fmt.Sprintf("%s/accounts/%s/publicKeys", c.UCPURL, user)

	var clientBundles ucptypes.ClientBundles

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	if err != nil {
		log.Debugf("response: %v", err)
		parseerr := ParseUCPError([]byte(response))
		if parseerr != nil {
			log.Debugf("%s", response)
			log.Errorf("Error parsing UCP error: %v", err)
		}
		return nil, err
	}

	err = json.Unmarshal(response, &clientBundles)
	if err != nil {
		return nil, err
	}

	return &clientBundles, nil
}

// RenameClientBundle - will relabel / rename the selected UCP Client Bundle
func (c *Client) RenameClientBundle(user, keyID, name string) error {
	if user == "" {
		currentAccount, err := c.AuthStatus()
		if err != nil {
			log.Warn("Session has expired, please login")
			return err
		}
		user = currentAccount.Name
	}
	log.Infoln("Renaming the UCP Client Bundle")
	url := fmt.Sprintf("%s/accounts/%s/publicKeys/%s", c.UCPURL, user, keyID)

	type bundleLabel struct {
		Label string `json:"label"`
	}

	label := bundleLabel{Label: name}

	b, err := json.Marshal(label)
	if err != nil {
		return nil
	}

	response, err := c.patchRequest(url, b)
	if err != nil {
		return err
	}

	if err != nil {
		log.Debugf("response: %v", err)
		parseerr := ParseUCPError([]byte(response))
		if parseerr != nil {
			log.Debugf("%s", response)
			log.Errorf("Error parsing UCP error: %v", err)
		}
		return err
	}

	return nil
}

// DeleteClientBundle - will delete the selected UCP Client Bundle
func (c *Client) DeleteClientBundle(user, keyID string) error {
	if user == "" {
		currentAccount, err := c.AuthStatus()
		if err != nil {
			log.Warn("Session has expired, please login")
			return err
		}
		user = currentAccount.Name
	}
	log.Infoln("Deleting the UCP Client Bundle")
	url := fmt.Sprintf("%s/accounts/%s/publicKeys/%s", c.UCPURL, user, keyID)

	response, err := c.delRequest(url, nil)
	if err != nil {
		return err
	}

	if err != nil {
		log.Debugf("response: %v", err)
		parseerr := ParseUCPError([]byte(response))
		if parseerr != nil {
			log.Debugf("%s", response)
			log.Errorf("Error parsing UCP error: %v", err)
		}
		return err
	}

	return nil
}
