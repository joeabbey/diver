package ucp

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

// GetUCPVersion - this returns the running version of the current session
func (c *Client) GetUCPVersion() (string, error) {

	url := fmt.Sprintf("%s/version", c.UCPURL)
	response, err := c.getRequest(url, nil)
	if err != nil {
		return "", err
	}
	log.Debugf("%s", response)

	var responseData map[string]interface{}

	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return "", err
	}
	if responseData["Version"] != nil {
		return responseData["Version"].(string), nil
	}

	return "", fmt.Errorf("Couldn't determine UCP version")
}

//UpgradeUCP - this will start the upgrade procedure to a newer version
func (c *Client) UpgradeUCP(version string) error {

	url := fmt.Sprintf("%s/api/ucpversionupdate?imageVersion=%s", c.UCPURL, version)

	response, err := c.postRequest(url, nil)
	if err != nil {
		parseerr := ParseUCPError(response)
		if err != nil {
			log.Errorf("Error parsing UCP error: %v", parseerr)
		}
		return err
	}

	log.Debugf("%v", string(response))

	return nil
}

//GetAvailavleUCPVersions - this will start the upgrade procedure to a newer version
func (c *Client) GetAvailavleUCPVersions() ([]string, error) {

	url := fmt.Sprintf("%s/api/ucpversions", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		parseerr := ParseUCPError(response)
		if err != nil {
			log.Errorf("Error parsing UCP error: %v", parseerr)
		}
		return nil, err
	}

	var availableVersions []string

	err = json.Unmarshal(response, &availableVersions)
	if err != nil {
		return nil, err
	}
	return availableVersions, nil
}
